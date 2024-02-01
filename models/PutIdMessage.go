package models

import(
  "time"
  "gorm.io/gorm"
)

func (m *Messages)PutIdMessage(id int, db *gorm.DB)error{
  header:=m.Header
  message:=m.Message
  if err := db.First(&m, id).Error; err!=nil{
    return err
  }
  now:=time.Now()
  m.Updated=now.Format("2006-01-02 15:04:05")
  m.Message=message
  m.Header=header
  if err := db.Save(&m).Error; err!=nil{
    return err
  }
  return nil
}
