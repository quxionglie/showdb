package db

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	defaultOracleModel struct {
		conn sqlx.SqlConn
	}
)

type DbTableOracle struct {
	TableName    string         `db:"TABLE_NAME" json:"table_name"`
	TableComment sql.NullString `db:"TABLE_COMMENT" json:"table_comment"`
	CreateTime   sql.NullTime   `db:"CREATE_TIME" json:"create_time"`
	UpdateTime   sql.NullTime   `db:"UPDATE_TIME" json:"update_time"`
}

type DbColumnOracle struct {
	ColumnName    string         `db:"COLUMN_NAME" json:"column_name"`
	IsRequired    string         `db:"IS_REQUIRED" json:"is_required"`
	IsPk          string         `db:"IS_PK" json:"is_pk"`
	Sort          int32          `db:"SORT" json:"sort"`
	IsIncrement   string         `db:"IS_INCREMENT" json:"is_increment"`
	ColumnType    string         `db:"COLUMN_TYPE" json:"column_type"`
	ColumnComment sql.NullString `db:"COLUMN_COMMENT" json:"column_comment"`
	ColumnDefault sql.NullString `db:"COLUMN_DEFAULT" json:"column_default"`
}

var _ DbStore = (*defaultOracleModel)(nil)

func NewOracleStore(conn sqlx.SqlConn) *defaultOracleModel {
	return &defaultOracleModel{
		conn: conn,
	}
}

func (m *defaultOracleModel) GetTables(ctx context.Context) ([]*DbTable, error) {
	query := ` select lower(dt.table_name) as table_name, dtc.comments as table_comment, 
        uo.created as create_time, uo.last_ddl_time as update_time
            from user_tables dt, user_tab_comments dtc, user_objects uo
            where dt.table_name = dtc.table_name
            and dt.table_name = uo.object_name
            and uo.object_type = 'TABLE'
 			order by create_time desc
				`
	var resp []*DbTableOracle
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return convertTable(resp), nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func convertTable(in []*DbTableOracle) []*DbTable {
	if in == nil {
		return nil
	}
	out := []*DbTable{}
	for _, table := range in {
		out = append(out, &DbTable{
			TableName:    table.TableName,
			TableComment: table.TableComment,
			CreateTime:   table.CreateTime,
			UpdateTime:   table.UpdateTime,
		})
	}
	return out
}

func convertColumn(in []*DbColumnOracle) []*DbColumn {
	if in == nil {
		return nil
	}
	out := []*DbColumn{}
	for _, col := range in {
		out = append(out, &DbColumn{
			ColumnName:    col.ColumnName,
			IsRequired:    col.IsRequired,
			IsPk:          col.IsPk,
			Sort:          col.Sort,
			IsIncrement:   col.IsIncrement,
			ColumnType:    col.ColumnType,
			ColumnComment: col.ColumnComment,
			ColumnDefault: col.ColumnDefault,
		})
	}
	return out
}

func (m *defaultOracleModel) GetColumns(ctx context.Context, tableName string) ([]*DbColumn, error) {
	query := `select lower(temp.column_name) as column_name,
                    (case when (temp.nullable = 'N'  and  temp.constraint_type != 'P') then '1' else '0' end) as is_required,
                    (case when temp.constraint_type = 'P' then '1' else '0' end) as is_pk,
                    temp.column_id as sort,
                    temp.comments as column_comment,
                    (case when temp.constraint_type = 'P' then '1' else '0' end) as is_increment,
                    lower(temp.data_type) as column_type,
					temp.column_default as column_default
            from (
                select col.column_id, col.column_name,col.nullable, col.data_type, colc.comments, 
                       uc.constraint_type, row_number()
                    over (partition by col.column_name order by uc.constraint_type desc) as row_flg,
                col.data_default as column_default
                from user_tab_columns col
                left join user_col_comments colc on colc.table_name = col.table_name and colc.column_name = col.column_name
                left join user_cons_columns ucc on ucc.table_name = col.table_name and ucc.column_name = col.column_name
                left join user_constraints uc on uc.constraint_name = ucc.constraint_name
                where col.table_name = upper(:1)
            ) temp
            WHERE temp.row_flg = 1
            ORDER BY temp.column_id
				`
	var resp []*DbColumnOracle
	err := m.conn.QueryRowsCtx(ctx, &resp, query, tableName)
	switch err {
	case nil:
		return convertColumn(resp), nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOracleModel) ShowCreateTable(ctx context.Context, tableName string) (string, error) {
	query := `select DBMS_METADATA.GET_DDL('TABLE',:1,'') as ddl  from DUAL`
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, tableName)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return "", ErrNotFound
	default:
		return "", err
	}
}
