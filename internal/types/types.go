// Code generated by goctl. DO NOT EDIT.
package types

type Column struct {
	ColumnName    string `json:"columnName"`
	IsRequired    string `json:"isRequired"`
	IsPk          string `json:"isPk"`
	Sort          int32  `json:"sort"`
	IsIncrement   string `json:"isIncrement"`
	ColumnType    string `json:"columnType"`
	ColumnComment string `json:"columnComment"`
	ColumnDefault string `json:"columnDefault"`
}

type Database struct {
	DbName    string  `json:"dbName"`
	DbComment string  `json:"dbComment"`
	DbType    string  `json:"dbType"`
	Tables    []Table `json:"tables"`
}

type GetColumnsRequest struct {
	DbName    string `form:"dbName"`
	TableName string `form:"tableName"`
}

type GetColumnsResponse struct {
	Columns []Column `json:"columns"`
}

type GetDatabasesRequest struct {
}

type GetDatabasesResponse struct {
	Database []Database `json:"databases"`
}

type GetTablesRequest struct {
	DbName string `form:"dbName"` // 必填参数
}

type GetTablesResponse struct {
	Tables []Table `json:"tables"`
}

type ShowCreateTableRequest struct {
	DbName    string `form:"dbName"`
	TableName string `form:"tableName"`
}

type ShowCreateTableResponse struct {
	Ddl string `json:"ddl"`
}

type Table struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`
	CreateTime   string `json:"createTime"`
}
