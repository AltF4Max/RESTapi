package main

import (
  "fmt"
  "strings"
  "gorm.io/gorm"
  "net/http"
  "encoding/json"
  "context"
  "github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
  "RESTapi/handlers"
  "RESTapi/models"
  "RESTapi/dbgorm"
  "RESTapi/config"
)

func main(){
  fmt.Println("Start")
  db, sqlDB, err:= dbgorm.ConnectionDB()
  if err!=nil{
    fmt.Println("dbgorm.ConnectionDB()", err)
    return
  }
  defer sqlDB.Close()

  r := mux.NewRouter()
  server := &http.Server{
    Addr:         config.Config.Addr,
    Handler:      r,
    ReadTimeout:  config.Config.ReadTimeout,
    WriteTimeout: config.Config.WriteTimeout,
  }

  r.HandleFunc("/login", handlers.Login(db)).Methods("POST")
  r.HandleFunc("/register", handlers.Register(db)).Methods("POST")
  r.Use(verifyTokenFromDB(db, config.JwtSecret, "/login", "/register"))
  r.HandleFunc("/api/messages", handlers.PostMsg(db)).Methods("POST")
  r.HandleFunc("/api/messages", handlers.GetMsg(db)).Methods("GET")
  r.HandleFunc("/api/messages/{id}", handlers.GetMessageId(db)).Methods("GET")
  r.HandleFunc("/api/messages/{id}", handlers.PutMessageId(db)).Methods("PUT")
  r.HandleFunc("/api/messages/{id}", handlers.DeleteMessageId(db)).Methods("DELETE")

  err = server.ListenAndServe()
  if err!=nil{
    fmt.Println("Error starting server:", err)
    return
  }
}
func verifyTokenFromDB(db *gorm.DB, jwtSecret []byte, excludePaths ...string) mux.MiddlewareFunc {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      for _, path := range excludePaths {
        if r.URL.Path == path {
          next.ServeHTTP(w, r)
          return
        }
      }

      tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
      _, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {//token
        return jwtSecret, nil
      })
      if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        answer:=models.Answer{Msg: "Unauthorized"}
        json.NewEncoder(w).Encode(answer)
        return
      }

      //if _, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {//claims
        u:=models.Users{}
        u.TokenJWT=tokenString
        if err := db.Where("token_jwt = ?", u.TokenJWT).First(&u).Error; err != nil {
          w.WriteHeader(http.StatusUnauthorized)
          answer:=models.Answer{Msg: "Unauthorized"}
          json.NewEncoder(w).Encode(answer)
          return
        }
        /*if claims.UserID != u.Login {//Кажеттся это не нужно
          w.WriteHeader(http.StatusUnauthorized)
          answer:=models.Answer{Msg: "Unauthorized"}
          json.NewEncoder(w).Encode(answer)
          return
        }*/
        ctx := context.WithValue(r.Context(), "User", u)
        next.ServeHTTP(w, r.WithContext(ctx))
        return
    /*}
    w.WriteHeader(http.StatusUnauthorized)
    answer:=models.Answer{Msg: "Unauthorized"}
    json.NewEncoder(w).Encode(answer)
    return*/
  })
}
}
