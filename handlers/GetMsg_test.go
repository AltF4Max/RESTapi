package handlers
import(
  "testing"
  "regexp"
  "github.com/DATA-DOG/go-sqlmock"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "net/http/httptest"
  "net/http"
  "RESTapi/models"
  "context"

)
func TestGetMsg(t *testing.T) {
  testStruct:= []TestStructMessage{
    {statusCode: 200, contextKey: "User", fakeDB: 1},
    {statusCode: 404, contextKey: "User"},
    {statusCode: 500},
  }
  u:=models.Users{
    Id:           22,
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
      AddRow(1, "admin", "text", "testtext", "1.1", "1.1").
      AddRow(2, "userwiki", "text", "test", "11.11", "12.12").
      AddRow(3, "userwiki", "text", "testtext", "11.12", "12.12")

      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `messages` WHERE login = ?")).
      WithArgs("userwiki").
      WillReturnRows(rows)
    }
    req, err := http.NewRequest("GET", "/api/messages/", nil)
    if err!=nil{
      t.Fatalf("Failed to create request: %v", err)
    }
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    handler := GetMsg(gormDB)
    handler.ServeHTTP(rr, req)

    if rr.Code != oneTest.statusCode{
      t.Errorf("Expected status code %v, got %v", oneTest.statusCode, rr.Code)
    }
  }
}
