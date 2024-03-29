package handlers

import (
  "net/http"
  "gorm.io/gorm"
  "encoding/json"
  "golang.org/x/crypto/bcrypt"
  "RESTapi/models"
)

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var user models.Users
    err := json.NewDecoder(r.Body).Decode(&user)
    if err!=nil{
      w.WriteHeader(http.StatusBadRequest)
      answer:=models.Answer{Msg: "Error reading data from request body"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    if len(user.Login)==0 || len(user.PasswordHash)==0{
      w.WriteHeader(http.StatusBadRequest)
      answer:=models.Answer{Msg: "Login or password field is empty"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    token, err:=authorization(db, user.Login, user.PasswordHash)
    if err!=nil{
      w.WriteHeader(http.StatusBadRequest)
      answer:=models.Answer{Msg: "Wrong password or login"}
      json.NewEncoder(w).Encode(answer)
      return
    }
    answer:=models.Answer{Msg: token}
    json.NewEncoder(w).Encode(answer)
  }
}
func authorization(db *gorm.DB, u string, p string)(string, error){
  var user models.Users
	if err := db.Where("login = ?", u).First(&user).Error; err!=nil{
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p))
	if err != nil {
		return "", err//Пароль не совпадает
	}
	return user.TokenJWT, nil
}
