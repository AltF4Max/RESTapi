package handlers

import(
  "testing"
  "strings"
  "net/http"
  "net/http/httptest"
  "gorm.io/gorm"
  "regexp"
  "gorm.io/driver/mysql"
  "github.com/DATA-DOG/go-sqlmock"
)

func TestLogin(t *testing.T) {
  testStruct:= []TestStruct{
    {strJson: `{"login": "userwiki", "password": "admin"}`, statusCode: 200, fakeDB: 1},
    {strJson: `{"login": "userwiki", "password": "admin"}`, statusCode: 400},
    {strJson: `{"login": "", "password": "admin"}`, statusCode: 400},
    {strJson: ``, statusCode: 400},
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
    if oneTest.fakeDB==1{
      rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "token_jwt", "email", "last_name", "first_name"}).
      AddRow(1,"userwiki", "$2a$10$s0BO.GcCwDawOOAQwNpoc.I/uKUCiVPFZ0hljtiBvcYtl1lY0km/6", "tokenOK","email@email", "Wiki", "Kek")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("userwiki").
      WillReturnRows(rows)
    }
    req, err := http.NewRequest("POST", "/login", strings.NewReader(oneTest.strJson))
    if err!=nil{
      t.Fatalf("Failed to create request: %v", err)
    }
    rr := httptest.NewRecorder()
    handler := Login(gormDB)
    handler.ServeHTTP(rr, req)

    if rr.Code != oneTest.statusCode{
      t.Errorf("Expected status code %v, got %v", oneTest.statusCode, rr.Code)
    }
  }
}
