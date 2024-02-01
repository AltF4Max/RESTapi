package main

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "regexp"
  "github.com/DATA-DOG/go-sqlmock"
  "github.com/stretchr/testify/assert"
)

type TestStruct struct{
  Token string
  Url   string
  statusCode int
  fakeDB     int
}
func TestVerifyTokenFromDB(t *testing.T){
  testStruct:= []TestStruct{
    {Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbiIsImV4cCI6MTcxMTQ3MjQ4OH0.UI79VMNx3uKlqhB3U_VKihm3cmOLDEy0Nrvh8wELmh8", Url: "/test", statusCode: 200, fakeDB: 1},
    {Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZGNjbWRkZGluIiwiZXhwIjoxNzExOTg3ODA0fQ.TMFIhkcpM3KiJirRWJyi89f5VZB_9v8qX4zcMt68ReQ", Url: "/test", statusCode: 401},
    {Token: "", Url: "/test", statusCode: 401},
    {Url: "/login", statusCode: 200},
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

  testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
  })

  excludePaths := []string{"/login"} // Путь, который мы хотим исключить из проверки токена
  handler := verifyTokenFromDB(gormDB, []byte("f4keraven"), excludePaths...)(testHandler)

  for _, oneTest := range testStruct{
    if oneTest.fakeDB == 1 {
      rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "token_jwt", "email", "last_name", "first_name"}).
      AddRow(1,"admin", "passwordhash", "tokenjwt", "email", "Vika", "Kek")
      mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE token_jwt = ? ORDER BY `users`.`id` LIMIT 1")).
      WithArgs(oneTest.Token).
      WillReturnRows(rows)
    }
    req, _ := http.NewRequest("GET", oneTest.Url, nil)
    req.Header.Set("Authorization", "Bearer "+oneTest.Token)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)
    assert.Equal(t, oneTest.statusCode, rr.Code, "Error")
  }
}
