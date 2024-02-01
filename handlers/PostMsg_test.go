package handlers

import(
  "testing"
  "strings"
  "net/http"
  "net/http/httptest"
  "context"
  "regexp"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "github.com/DATA-DOG/go-sqlmock"
  "RESTapi/models"
)

func TestPostMsg(t *testing.T) {
  testStruct:= []TestStructMessage{
    {strJson: `{"Header": "Wiki", "Message": "testtext"}`, statusCode: 200, contextKey: "User" , fakeDB: 1},
    {strJson: `{"Header": "Wiki", "Message": "testtext"}`, statusCode: 500, contextKey: "User"},
    {strJson: `{"Header": "Wiki", "Message": ""}`, statusCode: 400, contextKey: "User"},
    {strJson: ``, statusCode: 400, contextKey: "User"},
    {statusCode: 500},
  }
  u:=models.Users{
    Id:           22,
    Login:        "uservika",
    PasswordHash: "nope",
    TokenJWT:     "nope",
    Email:        "nope",
    LastName:     "nope",
    FirstName:    "nope",
  }
  mockDB, mock, err := sqlmock.New()
  if err != nil {
      t.Fatalf("Failed to create mock database: %v", err)
  }
  defer mockDB.Close()

  gormDB, err := gorm.Open(mysql.New(mysql.Config{
    Conn:                      mockDB,
    SkipInitializeWithVersion: true,
  }), &gorm.Config{})
  if err != nil {
    t.Fatalf("Failed to open GORM database: %v", err)
  }

  for _, oneTest := range testStruct {
    ctx := context.WithValue(context.Background(), oneTest.contextKey, u)
    if oneTest.fakeDB==1{
      mock.ExpectBegin()
      mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `messages`")).
      WillReturnResult(sqlmock.NewResult(1, 1))
      mock.ExpectCommit()
    }
    req, err := http.NewRequest("POST", "/api/messages/", strings.NewReader(oneTest.strJson))
    if err!=nil{
      t.Fatalf("Failed to create request: %v", err)
    }
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    handler := PostMsg(gormDB)
    handler.ServeHTTP(rr, req)

    if rr.Code != oneTest.statusCode{
      t.Errorf("Expected status code %v, got %v", oneTest.statusCode, rr.Code)
    }
  }
}
