package handlers

import(
  "testing"
  "strings"
  "net/http"
  "net/http/httptest"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "github.com/DATA-DOG/go-sqlmock"
  "regexp"
)

type TestStruct struct{
  strJson    string
  statusCode int
  fakeDB     int
}
func TestRegister(t *testing.T){
  testStruct:=[]TestStruct{
    {strJson: `{"login": "admin", "password": "testpass", "email": "test@example.com", "lastname": "Test", "firstname": "User"}`, statusCode: 200, fakeDB: 1},
    {strJson: `{"login": "admin", "password": "testpass", "email": "test@example.com", "lastname": "Test", "firstname": "User"}`, statusCode: 500},
    {strJson: `{"login": "admin", "password": "testpass", "email": "", "lastname": "Test", "firstname": "User"}`, statusCode: 400},
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
      rows := sqlmock.NewRows([]string{})
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("admin").
      WillReturnRows(rows)

      mock.ExpectBegin()
      mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
      WillReturnResult(sqlmock.NewResult(1, 1))
      mock.ExpectCommit()
    }
    req, err := http.NewRequest("POST", "/register", strings.NewReader(oneTest.strJson))
    if err!=nil{
      t.Fatalf("Failed to create request: %v", err)
    }
    rr := httptest.NewRecorder()
    handler := Register(gormDB)
    handler.ServeHTTP(rr, req)
    if rr.Code!=oneTest.statusCode{
      t.Errorf("Expected status code %v, got %v", oneTest.statusCode, rr.Code)
    }
  }
}
