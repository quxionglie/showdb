package svc

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/quxionglie/showdb/internal/config"
	"github.com/quxionglie/showdb/internal/db"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
	"strings"
)

type ServiceContext struct {
	Config   config.Config
	SqlConns map[string]sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		SqlConns: initSqlConn(c),
	}
}

func initSqlConn(c config.Config) map[string]sqlx.SqlConn {
	sqlConns := map[string]sqlx.SqlConn{}
	for _, database := range c.Databases {
		sqlConn, err := buildSqlConn(database.Name, c)
		if err != nil {
			logx.Errorf("buildSqlStore err=%v", err)
			os.Exit(-1)
		}
		if sqlConn != nil {
			sqlConns[database.Name] = sqlConn
		}
	}
	return sqlConns
}

func buildSqlConn(dbName string, c config.Config) (sqlx.SqlConn, error) {
	dbCfg := c.GetDatabase(dbName)
	if strings.ToLower(dbCfg.DbType) == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=%s&loc=Local",
			dbCfg.Username,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Database,
			dbCfg.Charset,
		)
		conn := sqlx.NewMysql(dsn)
		return conn, nil
	} else if strings.ToLower(dbCfg.DbType) == "oracle" {
		connStr := go_ora.BuildUrl(dbCfg.Host, dbCfg.Port, dbCfg.ServerName, dbCfg.Username, dbCfg.Password, nil)
		//connStr := go_ora.BuildUrl("192.168.0.96", 1528, "uat", "ehc", "ehc", nil)
		//oracle://ehc:ehc@192.168.0.96:1528/uat
		conn := sqlx.NewSqlConn("oracle", connStr)
		return conn, nil
	} else if strings.ToLower(dbCfg.DbType) == "pg" {
		//dsn := "host=X.X.X.X port=54321 user=postgres password=admin123 dbname=postgres sslmode=disable"
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Username,
			dbCfg.Password,
			dbCfg.Database)
		conn := sqlx.NewSqlConn("postgres", dsn)
		return conn, nil
	}
	return nil, errors.New("db type not support")
}

func (s *ServiceContext) GetDbStore(dbName string) (db.DbStore, error) {
	dbCfg := s.Config.GetDatabase(dbName)
	if dbCfg == nil {
		return nil, fmt.Errorf("db not found")
	}
	conn := s.SqlConns[dbName]
	if strings.ToLower(dbCfg.DbType) == "mysql" {
		sqlStore := db.NewMysqlStore(conn)
		return sqlStore, nil
	} else if strings.ToLower(dbCfg.DbType) == "oracle" {
		sqlStore := db.NewOracleStore(conn)
		return sqlStore, nil
	} else if strings.ToLower(dbCfg.DbType) == "pg" {
		sqlStore := db.NewPgStore(conn)
		return sqlStore, nil
	}
	return nil, fmt.Errorf("db not found")
}
