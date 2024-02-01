package handlers

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "strconv"
  "gorm.io/gorm"
  "RESTapi/models"
)

func GetMessageId(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    user, ok := r.Context().Value("User").(models.Users)
    if !ok{
      w.WriteHeader(http.StatusInternalServerError)
      answer := models.Answer{Msg: "Error 500 please try again later"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    id := mux.Vars(r)
    s := id["id"]
    i, err := strconv.Atoi(s)
    if err!=nil{
      w.WriteHeader(http.StatusBadRequest)
      answer:=models.Answer{Msg: "Invalid ID"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    err=user.CheckAccessWrite(i, db)
    if err!=nil{
      w.WriteHeader(http.StatusForbidden)
      answer:=models.Answer{Msg: "No access"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    message:=models.Messages{}
    err=message.GetIdMessage(i, db)
    if err!=nil{
      w.WriteHeader(http.StatusNotFound)
      answer:=models.Answer{Msg: "The requested resource was not found on the server"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    json.NewEncoder(w).Encode(message)
  }
}
