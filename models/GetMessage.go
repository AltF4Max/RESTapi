package models

import(
  "gorm.io/gorm"
)

func GetMessage(loginStr string,db *gorm.DB)([]Messages, error){
  var messages []Messages
  if err := db.Where("login = ?", loginStr).Find(&messages).Error; err!=nil{
    return nil, err
  }
  return messages, nil
}
