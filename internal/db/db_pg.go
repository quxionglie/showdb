package db

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	defaultPgModel struct {
		conn sqlx.SqlConn
	}
)

var _ DbStore = (*defaultPgModel)(nil)

func NewPgStore(conn sqlx.SqlConn) *defaultPgModel {
	return &defaultPgModel{
		conn: conn,
	}
}

func (m *defaultPgModel) GetTables(ctx context.Context) ([]*DbTable, error) {
	query := `  select table_name, table_comment, create_time, update_time
            from (
                SELECT c.relname AS table_name,
                        obj_description(c.oid) AS table_comment,
                        CURRENT_TIMESTAMP AS create_time,
                        CURRENT_TIMESTAMP AS update_time
                FROM pg_class c
                    LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
                WHERE (c.relkind = ANY (ARRAY ['r'::"char", 'p'::"char"]))
                    AND c.relname != 'spatial_%'::text
                    AND n.nspname = 'public'::name
                    AND n.nspname <> ''::name
            ) list_table
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

func (m *defaultPgModel) GetColumns(ctx context.Context, tableName string) ([]*DbColumn, error) {
	query := `SELECT column_name, is_required, is_pk, sort, column_comment, is_increment, column_type,
                      column_default
            FROM (
                SELECT c.relname AS table_name,
                       a.attname AS column_name,
                       d.description AS column_comment,
                       CASE WHEN a.attnotnull AND con.conname IS NULL THEN 1 ELSE 0
                       END AS is_required,
                       CASE WHEN con.conname IS NOT NULL THEN 1 ELSE 0
                       END AS is_pk,
                       a.attnum AS sort,
                       CASE WHEN "position"(pg_get_expr(ad.adbin, ad.adrelid),
                           ((c.relname::text || '_'::text) || a.attname::text) || '_seq'::text) > 0 THEN 1 ELSE 0
                       END AS is_increment,
                       btrim(
                           CASE WHEN t.typelem <>  0::oid AND t.typlen = '-1'::integer THEN 'ARRAY'::text ELSE
                                CASE WHEN t.typtype = 'd'::"char" THEN format_type(t.typbasetype, NULL::integer)
                                ELSE format_type(a.atttypid, NULL::integer) END
                           END, '"'::text
                       ) AS column_type,
                       col.column_default as column_default
                FROM pg_attribute a
                    JOIN (pg_class c JOIN pg_namespace n ON c.relnamespace = n.oid) ON a.attrelid = c.oid
                    LEFT JOIN pg_description d ON d.objoid = c.oid AND a.attnum = d.objsubid
                    LEFT JOIN pg_constraint con ON con.conrelid = c.oid AND (a.attnum = ANY (con.conkey))
                    LEFT JOIN pg_attrdef ad ON a.attrelid = ad.adrelid AND a.attnum = ad.adnum
                    LEFT JOIN pg_type t ON a.atttypid = t.oid
                    LEFT JOIN information_schema.columns as col ON col.table_name = c.relname and col.column_name = a.attname
                WHERE (c.relkind = ANY (ARRAY ['r'::"char", 'p'::"char"]))
                    AND a.attnum > 0
                    AND n.nspname = 'public'::name
                ORDER BY c.relname, a.attnum
            ) temp
            WHERE table_name = ($1)
                AND column_type <>  '-'
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

func (m *defaultPgModel) ShowCreateTable(ctx context.Context, tableName string) (string, error) {
	return "-- 尚未支持PG查询 TODO", nil
}
