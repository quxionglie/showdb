package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Databases []*Database
}

type Database struct {
	Name       string `json:",required"`
	DbType     string `json:",required"`
	Host       string `json:",required"`
	Port       int    `json:",required"`
	Username   string `json:",required"`
	Password   string `json:",required"`
	Database   string `json:",required"`
	Charset    string `json:",optional"`
	ServerName string `json:",optional"`
}

func (c *Config) GetDatabase(name string) *Database {
	for _, db := range c.Databases {
		if db.Name == name {
			return db
		}
	}
	return nil
}
