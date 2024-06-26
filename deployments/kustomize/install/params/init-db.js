const mongoHost = process.env.SURGEON_API_MONGODB_HOST;
const mongoPort = process.env.SURGEON_API_MONGODB_PORT;

const mongoUser = process.env.SURGEON_API_MONGODB_USERNAME;
const mongoPassword = process.env.SURGEON_API_MONGODB_PASSWORD;

const database = process.env.SURGEON_API_MONGODB_DATABASE;
const collection = process.env.SURGEON_API_MONGODB_COLLECTION;

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while (true) {
  try {
    connection = Mongo(
      `mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`
    );
    break;
  } catch (exception) {
    print(`Cannot connect to mongoDB: ${exception}`);
    print(`Will retry after ${retrySeconds} seconds`);
    sleep(retrySeconds * 1000);
  }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames();
if (databases.includes(database)) {
  const dbInstance = connection.getDB(database);
  collections = dbInstance.getCollectionNames();
  if (collections.includes(collection)) {
    print(
      `Collection '${collection}' already exists in database '${database}'`
    );
    process.exit(0);
  }
}

// initialize
// create database and collection
const db = connection.getDB(database);
db.createCollection(collection);

// create indexes
db[collection].createIndex({ id: 1 });

//insert sample data
let result = db[collection].insertMany([
  {
    Id: "1",
    Name: "MuDr. Andrej Poljako",
    Surgeries: [
      {
        Id: "1",
        PatientId: "123-pat1",
        Date: "2024-05-06",
        SurgeryNote: "Uplne vsetko v poriadku.",
        Successful: true,
        OperatedLimb: {
          Value: "Prava ruka",
          Code: "Right hand",
        },
      },
      {
        Id: "2",
        PatientId: "123-pat2",
        Date: "2024-01-01",
        SurgeryNote: "Zle dopadla operacia, nema oko.",
        Successful: false,
        OperatedLimb: {
          Value: "Hlava",
          Code: "Head",
        },
      },
    ],
  },
  {
    Id: "2",
    Name: "MuDr. Hipko",
    Surgeries: [
      {
        Id: "3",
        PatientId: "876-pat3",
        Date: "2024-05-06",
        SurgeryNote: "Prebehlo bez problemov.",
        Successful: true,
        OperatedLimb: {
          Value: "Lava noha",
          Code: "Left leg",
        },
      },
      {
        Id: "4",
        PatientId: "123-pat1",
        Date: "2024-01-01",
        SurgeryNote: "Odlupila sa cast mozgu.",
        Successful: false,
        OperatedLimb: {
          Value: "Hlava",
          Code: "Head",
        },
      },
    ],
  },
]);

if (result.writeError) {
  console.error(result);
  print(`Error when writing the data: ${result.errmsg}`);
}

// exit with success
process.exit(0);
