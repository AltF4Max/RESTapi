package models

import (
"errors"
"time"
"gorm.io/gorm"
"golang.org/x/crypto/bcrypt"
"github.com/dgrijalva/jwt-go"
"RESTapi/config"
)

func (u *Users)RegisterUser(db *gorm.DB)error{
  ok:=busyUser(u.Login, db)
  if ok==true{
    return errors.New("Login busy")
  }
  passwordByte, err:=CreationPws(u.PasswordHash)
  if err!=nil{
    return err
  }
  u.PasswordHash=string(passwordByte)
  u.TokenJWT, err=CreationToken(u.Login)
  if err!=nil{
    return err
  }
  if err := db.Create(&u).Error; err!=nil{
    return errors.New("Error create new user")
  }
  return nil
}

func busyUser(loginStr string, db *gorm.DB)bool{
  var user Users
  result := db.Where("login = ?", loginStr).First(&user)
  if result.Error==nil{
    return true
  }
  return false
}

func CreationPws(p string)([]byte, error){
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
  if err != nil {
    err = errors.New("Error CreationPws")
    return nil, err
  }
  return hashedPassword, nil
}

func CreationToken(u string)(string, error){
  claims:=&CustomClaims{
    UserID:         u,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: time.Now().Add(time.Hour * 6000).Unix(),
    },
  }
  tokenString, err:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.JwtSecret)
  if err!=nil{
    err=errors.New("Error creating token")
    return "", err
  }
  return tokenString, nil
}
