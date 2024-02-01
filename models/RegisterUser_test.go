package models

import (
  "fmt"
  "testing"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "regexp"
  "github.com/DATA-DOG/go-sqlmock"
  "errors"
)

type TestStructReg struct{
  errFinal    error
  fakeDB      int
}
func TestRegisterUser(t *testing.T){
  testStructReg:=[]TestStructReg{
    {fakeDB: 1},
    {errFinal: errors.New("Error create new user"),fakeDB: 2},
    {errFinal: errors.New("Login busy"),fakeDB: 3},
  }
  user:= Users{
    Login:        "uservika",
    PasswordHash: "passwordhash",
    Email:        "email",
    LastName:     "Vika",
    FirstName:    "Kek",
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

  for _, oneTest:=range testStructReg{
    if oneTest.fakeDB==1{
      rows := sqlmock.NewRows([]string{})
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("uservika").
      WillReturnRows(rows)

      mock.ExpectBegin()
      mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
      WillReturnResult(sqlmock.NewResult(1, 1))
      mock.ExpectCommit()
    } else if oneTest.fakeDB==2{
      rows := sqlmock.NewRows([]string{})
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("uservika").
      WillReturnRows(rows)
    } else if oneTest.fakeDB==3{
      rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "token_jwt", "email", "last_name", "first_name"}).
      AddRow(1,"uservika", "password", "token", "email@email", "Vika", "Kek")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE login = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs("uservika").
      WillReturnRows(rows)
    }
    err=user.RegisterUser(gormDB)
    if fmt.Sprintf("%v", err)!=fmt.Sprintf("%v", oneTest.errFinal){
      t.Errorf("Expected error %v, got %v", oneTest.errFinal, err)
    }
  }
}
