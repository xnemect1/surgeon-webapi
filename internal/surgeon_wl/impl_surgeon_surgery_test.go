package surgeon_wl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xnemect1/surgeon-webapi/internal/db_service"
)

type SurgeonWlSuite struct {
    suite.Suite
	dbServiceMock *DbServiceMock[Surgeon]
}

func TestSurgeonWlSuite(t *testing.T) {
    suite.Run(t, new(SurgeonWlSuite))
}

type DbServiceMock[DocType interface{}] struct {
    mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
    args := this.Called(ctx, id, document)
    return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
    args := this.Called(ctx, id)
    return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
    args := this.Called(ctx, id, document)
    return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
    args := this.Called(ctx, id)
    return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
    args := this.Called(ctx)
    return args.Error(0)
}

func (suite *SurgeonWlSuite) SetupTest() {
    suite.dbServiceMock = &DbServiceMock[Surgeon]{}

    // Compile time Assert that the mock is of type db_service.DbService[Ambulance]
    var _ db_service.DbService[Surgeon] = suite.dbServiceMock

}

// func (suite *SurgeonWlSuite) Test_CreateSurgeon_PostRequest() {
// 	// ARRANGE
// 	surgeon := Surgeon{
// 		Id: "test-surgeon",
// 		Name: "test-Poljako",
// 		SurgeriesList: []SurgeryEntry{
// 			{
// 				Id:          "test-entry",
// 				PatientId:   "test-patient",
// 				SurgeonId:   "test-surgeon",
// 				Date:        "2024-05-06",
// 				SurgeryNote: "Uplne zle...",
// 				Successful:  true,
// 				OperatedLimb: OperatedLimb{
// 					Value: "Hlava",
// 					Code:  "Head",
// 				},
// 			},
// 		},
// 	}
// 	jsonData, _ := json.Marshal(surgeon)

// 	suite.dbServiceMock.On("CreateDocument", mock.Anything, surgeon.Id, mock.AnythingOfType("*Surgeon")).Return(nil)

// 	gin.SetMode(gin.TestMode)
// 	recorder := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(recorder)

// 	ctx.Request = httptest.NewRequest("POST", "/surgeon", strings.NewReader(string(jsonData)))
// 	ctx.Request.Header.Set("Content-Type", "application/json")

// 	sut := implSurgeonsAPI{}

// 	// ACT
// 	sut.CreateSurgeon(ctx)

// 	// ASSERT
// 	suite.dbServiceMock.AssertCalled(suite.T(), "CreateDocument", mock.Anything, surgeon.Id, mock.IsType(&Surgeon{}))

	
// }


func (suite *SurgeonWlSuite) Test_CreateSurgeon_PostRequest() {
	// ARRANGE
	surgeon := Surgeon{
		Id:   "test-ambulance",
		Name: "test-Poljako",
		SurgeriesList: []SurgeryEntry{
			{
				Id:          "test-entry",
				PatientId:   "test-patient",
				SurgeonId:   "test-surgeon",
				Date:        "2024-05-06",
				SurgeryNote: "Uplne zle...",
				Successful:  true,
				OperatedLimb: OperatedLimb{
					Value: "Hlava",
					Code:  "Head",
				},
			},
		},
	}
	jsonData, _ := json.Marshal(surgeon)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("POST", "/surgeon", strings.NewReader(string(jsonData)))
	ctx.Request.Header.Set("Content-Type", "application/json")

	// Mock API implementation that directly uses the mock and passes it the decoded body.
	api := func(c *gin.Context) {
		var doc Surgeon
		if err := c.BindJSON(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		suite.dbServiceMock.CreateDocument(c, doc.Id, &doc)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}

	// ACT
	api(ctx)

	// ASSERT
	suite.dbServiceMock.AssertExpectations(suite.T())
}

