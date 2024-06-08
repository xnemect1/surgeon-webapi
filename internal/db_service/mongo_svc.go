package db_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbService[DocType interface{}] interface {
    CreateDocument(ctx context.Context, id string, document *DocType) error
    FindDocument(ctx context.Context, id string) (*DocType, error)
    UpdateDocument(ctx context.Context, id string, document *DocType) error
	GetAllDocuments(ctx context.Context) ([]*DocType, error)
    DeleteDocument(ctx context.Context, id string) error
    Disconnect(ctx context.Context) error
}

var ErrNotFound = fmt.Errorf("document not found")
var ErrConflict = fmt.Errorf("conflict: document already exists")

type MongoServiceConfig struct {
    ServerHost string
    ServerPort int
    UserName   string
    Password   string
    DbName     string
    Collection string
    Timeout    time.Duration
}

type mongoSvc[DocType interface{}] struct {
    MongoServiceConfig
    client     atomic.Pointer[mongo.Client]
    clientLock sync.Mutex
}

func NewMongoService[DocType interface{}](config MongoServiceConfig) DbService[DocType] {
	enviro := func(name string, defaultValue string) string {
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return defaultValue
	}

	svc := &mongoSvc[DocType]{}
	svc.MongoServiceConfig = config

	if svc.ServerHost == "" {
		svc.ServerHost = enviro("SURGEON_API_MONGODB_HOST", "localhost")
	}

	if svc.ServerPort == 0 {
		port := enviro("SURGEON_API_MONGODB_PORT", "27017")
		if port, err := strconv.Atoi(port); err == nil {
			svc.ServerPort = port
		} else {
			log.Printf("Invalid port value: %v", port)
			svc.ServerPort = 27017
		}
	}

	if svc.UserName == "" {
		svc.UserName = enviro("SURGEON_API_MONGODB_USERNAME", "")
	}

	if svc.Password == "" {
		svc.Password = enviro("SURGEON_API_MONGODB_PASSWORD", "")
	}

	if svc.DbName == "" {
		svc.DbName = enviro("SURGEON_API_MONGODB_DATABASE", "xnemect-surgeon")
	}

	if svc.Collection == "" {
		svc.Collection = enviro("SURGEON_API_MONGODB_COLLECTION", "surgeon")
	}

	if svc.Timeout == 0 {
		seconds := enviro("SURGEON_API_MONGODB_TIMEOUT_SECONDS", "10")
		if seconds, err := strconv.Atoi(seconds); err == nil {
			svc.Timeout = time.Duration(seconds) * time.Second
		} else {
			log.Printf("Invalid timeout value: %v", seconds)
			svc.Timeout = 10 * time.Second
		}
	}

	log.Printf(
		"MongoDB config: //%v@%v:%v/%v/%v",
		svc.UserName,
		svc.ServerHost,
		svc.ServerPort,
		svc.DbName,
		svc.Collection,
	)
	return svc
}

func (this *mongoSvc[DocType]) connect(ctx context.Context) (*mongo.Client, error) {
    // optimistic check
    client := this.client.Load()
    if client != nil {
        return client, nil
    }

    this.clientLock.Lock()
    defer this.clientLock.Unlock()
    // pesimistic check
    client = this.client.Load()
    if client != nil {
        return client, nil
    }

    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()

    var uri = fmt.Sprintf("mongodb://%v:%v", this.ServerHost, this.ServerPort)
    log.Printf("Using URI: " + uri)

    if len(this.UserName) != 0 {
        uri = fmt.Sprintf("mongodb://%v:%v@%v:%v", this.UserName, this.Password, this.ServerHost, this.ServerPort)
    }

    if client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetConnectTimeout(10*time.Second)); err != nil {
        return nil, err
    } else {
        this.client.Store(client)
        return client, nil
    }
}

func (this *mongoSvc[DocType]) Disconnect(ctx context.Context) error {
    client := this.client.Load()

    if client != nil {
        this.clientLock.Lock()
        defer this.clientLock.Unlock()

        client = this.client.Load()
        defer this.client.Store(nil)
        if client != nil {
            if err := client.Disconnect(ctx); err != nil {
                return err
            }
        }
    }
    return nil
}

func (this *mongoSvc[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(this.Collection)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil: // no error means there is conflicting document
        return ErrConflict
    case mongo.ErrNoDocuments:
        // do nothing, this is expected
    default: // other errors - return them
        return result.Err()
    }

    _, err = collection.InsertOne(ctx, document)
    return err
}

func (this *mongoSvc[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.connect(ctx)
    if err != nil {
        return nil, err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(this.Collection)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return nil, ErrNotFound
    default: // other errors - return them
        return nil, result.Err()
    }
    var document *DocType
    if err := result.Decode(&document); err != nil {
        return nil, err
    }
    return document, nil
}

func (this *mongoSvc[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(this.Collection)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return ErrNotFound
    default: // other errors - return them
        return result.Err()
    }
    _, err = collection.ReplaceOne(ctx, bson.D{{Key: "id", Value: id}}, document)
    return err
}

func (this *mongoSvc[DocType]) DeleteDocument(ctx context.Context, id string) error {
    ctx, contextCancel := context.WithTimeout(ctx, this.Timeout)
    defer contextCancel()
    client, err := this.connect(ctx)
    if err != nil {
        return err
    }
    db := client.Database(this.DbName)
    collection := db.Collection(this.Collection)
    result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
    switch result.Err() {
    case nil:
    case mongo.ErrNoDocuments:
        return ErrNotFound
    default: // other errors - return them
        return result.Err()
    }
    _, err = collection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
    return err
}

func (this *mongoSvc[DocType]) GetAllDocuments(ctx context.Context) ([]*DocType, error) {
    ctx, cancel := context.WithTimeout(ctx, this.Timeout)
    defer cancel()

    client, err := this.connect(ctx)
    if err != nil {
        return nil, err
    }

    db := client.Database(this.DbName)
    collection := db.Collection(this.Collection)

    cursor, err := collection.Find(ctx, bson.D{}) // Find all documents
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var documents []*DocType
    if err = cursor.All(ctx, &documents); err != nil {
        return nil, err
    }

    return documents, nil
}

// SEEDER
type Surgeon struct {
	Id   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type SurgeryEntry struct {
	Id           string       `bson:"id" json:"id"`
	SurgeonId    string       `bson:"surgeonId" json:"surgeonId"`
	PatientId    string       `bson:"patientId" json:"patientId"`
	Date         string       `bson:"date" json:"date"`
	Successful   bool         `bson:"successful" json:"successful"`
	SurgeryNote  string       `bson:"surgeryNote,omitempty" json:"surgeryNote,omitempty"`
	OperatedLimb OperatedLimb `bson:"operatedLimb" json:"operatedLimb"`
}

type OperatedLimb struct {
	Value string `bson:"value" json:"value"`
	Code  string `bson:"code" json:"code"`
}

func (svc *mongoSvc[DocType]) SeedDatabase(ctx context.Context) error {
    db := svc.client.Load().Database(svc.DbName)

    surgeonsCollection := db.Collection("surgeons")
    surgeriesCollection := db.Collection("surgeries")

    surgeons := []interface{}{
        Surgeon{Id: "1", Name: "MuDr. Andrej Poljak"},
        Surgeon{Id: "2", Name: "MuDr. Jan Vrba"},
    }

    surgeries := []interface{}{
        SurgeryEntry{Id: "s1", SurgeonId: "1", PatientId: "p1", Date: "2024-01-01", Successful: true, SurgeryNote: "Vyber znamienka", OperatedLimb: OperatedLimb{Value: "Lava ruka", Code: "Left hand"}},
        SurgeryEntry{Id: "s2", SurgeonId: "1", PatientId: "p2", Date: "2024-01-02", Successful: false, SurgeryNote: "Vymena kolenneho klbu", OperatedLimb: OperatedLimb{Value: "Prava noha", Code: "Right leg"}},
        SurgeryEntry{Id: "s3", SurgeonId: "2", PatientId: "p3", Date: "2024-01-03", Successful: true, SurgeryNote: "Amputacia ruky", OperatedLimb: OperatedLimb{Value: "Prava ruka", Code: "Right hand"}},
        SurgeryEntry{Id: "s4", SurgeonId: "2", PatientId: "p4", Date: "2024-01-04", Successful: false, SurgeryNote: "Vymena ACL", OperatedLimb: OperatedLimb{Value: "Lava noha", Code: "Left leg"}},
    }

    if _, err := surgeonsCollection.InsertMany(ctx, surgeons); err != nil {
        log.Printf("Failed to insert surgeons: %v", err)
        return err
    }

    if _, err := surgeriesCollection.InsertMany(ctx, surgeries); err != nil {
        log.Printf("Failed to insert surgeries: %v", err)
        return err
    }

    return nil
}


