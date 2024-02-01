package models

import(
  "gorm.io/gorm"
  "errors"
)

func (u *Users)CheckAccessWrite(id int,db *gorm.DB)error{
  var message Messages
  if err := db.First(&message, "id = ?", id).Error; err!=nil{
    return err
  }
  if u.Login!=message.Login{
    err := errors.New("Request someone else's date")
    return err
  }
  return nil
}
