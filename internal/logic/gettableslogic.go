package logic

import (
	"context"
	"github.com/quxionglie/showdb/internal/svc"
	"github.com/quxionglie/showdb/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type GetTablesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTablesLogic {
	return &GetTablesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTablesLogic) GetTables(req *types.GetTablesRequest) (resp *types.GetTablesResponse, err error) {
	dbStore, err := l.svcCtx.GetDbStore(req.DbName)
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

	return &types.GetTablesResponse{
		Tables: tables,
	}, nil
}
