package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	defaultMysqlModel struct {
		conn sqlx.SqlConn
	}
)

var _ DbStore = (*defaultMysqlModel)(nil)

type CreateTableDdl struct {
	TableName string         `db:"Table" json:"table_name"`
	Ddl       sql.NullString `db:"Create Table" json:"ddl"`
}

func NewMysqlStore(conn sqlx.SqlConn) *defaultMysqlModel {
	return &defaultMysqlModel{
		conn: conn,
	}
}

func (m *defaultMysqlModel) GetTables(ctx context.Context) ([]*DbTable, error) {
	query := `select table_name, table_comment, create_time, update_time
					from information_schema.tables
					where table_schema = (select database())
					order by create_time desc
				`
	var resp []*DbTable
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMysqlModel) GetColumns(ctx context.Context, tableName string) ([]*DbColumn, error) {
	query := `SELECT column_name,
                   (case when (is_nullable = 'no' and column_key != 'PRI') then '1' else '0' end) as is_required,
                   (case when column_key = 'PRI' then '1' else '0' end) as is_pk,
                   ordinal_position as sort,
                   column_comment,
                   column_default,
                   (case when extra = 'auto_increment' then '1' else '0' end) as is_increment,
                   column_type
			FROM information_schema.COLUMNS
			WHERE table_schema = DATABASE()
			  AND TABLE_NAME = ?
			order by ordinal_position
				`
	var resp []*DbColumn
	err := m.conn.QueryRowsCtx(ctx, &resp, query, tableName)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMysqlModel) ShowCreateTable(ctx context.Context, tableName string) (string, error) {
	query := fmt.Sprintf(`show create table %s`, tableName)
	var resp CreateTableDdl
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp.Ddl.String, nil
	case sqlx.ErrNotFound:
		return "", ErrNotFound
	default:
		return "", err
	}
}
