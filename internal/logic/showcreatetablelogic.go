package logic

import (
	"context"

	"github.com/quxionglie/showdb/internal/svc"
	"github.com/quxionglie/showdb/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowCreateTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowCreateTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowCreateTableLogic {
	return &ShowCreateTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowCreateTableLogic) ShowCreateTable(req *types.ShowCreateTableRequest) (resp *types.ShowCreateTableResponse, err error) {
	dbStore, err := l.svcCtx.GetDbStore(req.DbName)
	if err != nil {
		l.Logger.Error("dbStore not found")
		return nil, err
	}

	ddl, err := dbStore.ShowCreateTable(l.ctx, req.TableName)
	return &types.ShowCreateTableResponse{
		Ddl: ddl,
	}, err
}
