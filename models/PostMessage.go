package models

import(
  "time"
  "gorm.io/gorm"
)

func (m *Messages)PostMessage(loginStr string, db *gorm.DB)error{
  now:=time.Now()
  timeNow := now.Format("2006-01-02 15:04:05")
  m.Login=loginStr
  m.Created=timeNow
  m.Updated=timeNow

  tx := db.Begin()//Позволяет выполнять несколько операций как единое целое
  if err := tx.Omit("Id").Create(m).Error; err!=nil{
    tx.Rollback()//Если возникла ошибка при создании записи, откатываем транзакцию
    return err
  }
  tx.Commit()//Save
  return nil
}
