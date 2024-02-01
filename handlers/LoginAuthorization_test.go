package handlers

import(
  "testing"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "github.com/DATA-DOG/go-sqlmock"
  "regexp"
)

type AuthorizationStruct struct{
  Login       string
  Password    string
  TokenAnswer string
  fakeDB      int
}
func TestLoginAuthorization(t *testing.T) {
  authorizationStruct:= []AuthorizationStruct{
    {Login: "admin", Password: "admin", TokenAnswer: "tokenOK", fakeDB: 1},
    {Login: "admin", Password: "badpassword", TokenAnswer: "", fakeDB: 1},
    {Login: "admin", Password: "admin", TokenAnswer: ""},
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

  for _, oneTest := range authorizationStruct {
    if oneTest.fakeDB==1{
      rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "token_jwt", "email", "last_name", "first_name"}).
      AddRow(1,"admin", "$2a$10$s0BO.GcCwDawOOAQwNpoc.I/uKUCiVPFZ0hljtiBvcYtl1lY0km/6", "tokenOK", "email@email", "Wiki", "Kek")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("admin").
      WillReturnRows(rows)
    }
    token, _:= authorization(gormDB, oneTest.Login, oneTest.Password)//Можно сверять по ошибкам
    if token!=oneTest.TokenAnswer{
      t.Errorf("Expected token %v, got %v", oneTest.TokenAnswer, token)
    }
  }
}
