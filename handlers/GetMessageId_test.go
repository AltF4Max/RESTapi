package handlers

import(
  "testing"
  "net/http"
  "net/http/httptest"
  "context"
  "github.com/gorilla/mux"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "github.com/DATA-DOG/go-sqlmock"
  "regexp"
  "RESTapi/models"
)

func TestGetMessageId(t *testing.T) {
  testStruct:= []TestStructMessage{
    {statusCode: 200, urlId: "11", contextKey: "User", fakeDB: 1},
    {statusCode: 404, urlId: "11", contextKey: "User", fakeDB: 2},
    {statusCode: 403, urlId: "11", contextKey: "User"},
    {statusCode: 400, urlId: "11a", contextKey: "User"},
    {statusCode: 500, urlId: "11"},
  }
  u:=models.Users{
    Id:           23,
    Login:        "userwiki",
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

  for _,oneTest:=range testStruct{
    ctx := context.WithValue(context.Background(), oneTest.contextKey, u)
    if oneTest.fakeDB==1{
      rows := sqlmock.NewRows([]string{"id", "login", "header", "message", "created", "updated"}).
      AddRow(11, "userwiki", "text", "testtext", "11.11", "12.11")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `messages` WHERE id = ? ORDER BY `messages`.`id` LIMIT 1")).
      WithArgs(11).
      WillReturnRows(rows)

      rows2 := sqlmock.NewRows([]string{"id", "login", "header", "message", "created", "updated"}).
      AddRow(11, "userwiki", "text", "testtext", "11.11", "12.11")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `messages` WHERE id = ? ORDER BY `messages`.`id` LIMIT 1")).
      WithArgs(11).
      WillReturnRows(rows2)
    } else if oneTest.fakeDB==2{
      rows := sqlmock.NewRows([]string{"id", "login", "header", "message", "created", "updated"}).
      AddRow(11, "userwiki", "text", "testtext", "11.11", "12.11")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `messages` WHERE id = ? ORDER BY `messages`.`id` LIMIT 1")).
      WithArgs(11).
      WillReturnRows(rows)
    }
    req, err := http.NewRequest("GET", "/api/messages/"+oneTest.urlId, nil)
    if err!=nil{
      t.Fatalf("Failed to create request: %v", err)
    }
    vars := map[string]string{
      "id": oneTest.urlId,
    }

    req = req.WithContext(ctx)
    req = mux.SetURLVars(req, vars)
    rr := httptest.NewRecorder()
    handler := GetMessageId(gormDB)
    handler.ServeHTTP(rr, req)

    if rr.Code != oneTest.statusCode{
      t.Errorf("Expected status code %v, got %v", oneTest.statusCode, rr.Code)
    }
  }
}
