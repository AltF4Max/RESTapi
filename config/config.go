package config

import (
  "time"
)

var Dsn string = "root:@tcp(127.0.0.1:3306)/something"//Data Source Name
//"root:@tcp(127.0.0.1:3306)/something"
//"user:password@tcp(db:3306)/something"
var JwtSecret = []byte("f4keraven")
type ServerConfig struct {
  Addr         string
  ReadTimeout  time.Duration
  WriteTimeout time.Duration
}
var Config = ServerConfig{
  Addr:         ":8080",
  ReadTimeout:  5 * time.Second,
  WriteTimeout: 10 * time.Second,
}
