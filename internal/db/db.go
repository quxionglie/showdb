package db

import (
	"context"
	"database/sql"
)

type DbTable struct {
	TableName    string         `db:"table_name" json:"table_name"`
	TableComment sql.NullString `db:"table_comment" json:"table_comment"`
	CreateTime   sql.NullTime   `db:"create_time" json:"create_time"`
	UpdateTime   sql.NullTime   `db:"update_time" json:"update_time"`
}

type DbColumn struct {
	ColumnName    string         `db:"column_name" json:"column_name"`
	IsRequired    string         `db:"is_required" json:"is_required"`
	IsPk          string         `db:"is_pk" json:"is_pk"`
	Sort          int32          `db:"sort" json:"sort"`
	IsIncrement   string         `db:"is_increment" json:"is_increment"`
	ColumnType    string         `db:"column_type" json:"column_type"`
	ColumnComment sql.NullString `db:"column_comment" json:"column_comment"`
	ColumnDefault sql.NullString `db:"column_default" json:"column_default"`
}

// 数据库接口
type DbStore interface {
	// 查询表列表
	GetTables(ctx context.Context) ([]*DbTable, error)

	//根据表名查询字段列表
	GetColumns(ctx context.Context, tableName string) ([]*DbColumn, error)

	//查看建表ddl
	ShowCreateTable(ctx context.Context, tableName string) (string, error)
}
