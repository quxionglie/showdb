package logic

import (
	"context"
	"time"

	"github.com/quxionglie/showdb/internal/svc"
	"github.com/quxionglie/showdb/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDatabasesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDatabasesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDatabasesLogic {
	return &GetDatabasesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDatabasesLogic) GetDatabases(req *types.GetDatabasesRequest) (resp *types.GetDatabasesResponse, err error) {
	dbs := []types.Database{}

	for _, database := range l.svcCtx.Config.Databases {
		curDb := types.Database{
			DbName:    database.Name,
			DbComment: database.Database,
			DbType:    database.DbType,
		}

		dbStore, err := l.svcCtx.GetDbStore(database.Name)
		if err != nil {
			l.Logger.Error("dbStore not found")
			return nil, err
		}
		dbTables, err := dbStore.GetTables(l.ctx)
		if err != nil {
			return nil, err
		}

		tables := []types.Table{}
		for _, cur := range dbTables {
			t := types.Table{
				TableName:    cur.TableName,
				TableComment: cur.TableComment.String,
				CreateTime:   cur.CreateTime.Time.Format(time.RFC3339),
			}
			tables = append(tables, t)
		}
		curDb.Tables = tables

		dbs = append(dbs, curDb)
	}

	return &types.GetDatabasesResponse{
		Database: dbs,
	}, nil
}
