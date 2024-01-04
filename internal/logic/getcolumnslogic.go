package logic

import (
	"context"
	"database/sql"

	"github.com/quxionglie/showdb/internal/svc"
	"github.com/quxionglie/showdb/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetColumnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetColumnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetColumnsLogic {
	return &GetColumnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetColumnsLogic) GetColumns(req *types.GetColumnsRequest) (resp *types.GetColumnsResponse, err error) {
	dbStore, err := l.svcCtx.GetDbStore(req.DbName)
	if err != nil {
		l.Logger.Error("dbStore not found")
		return nil, err
	}

	dbColumns, err := dbStore.GetColumns(l.ctx, req.TableName)
	if err != nil {
		return nil, err
	}
	cols := []types.Column{}
	for _, cur := range dbColumns {
		t := types.Column{
			ColumnName:    cur.ColumnName,
			IsRequired:    yesOrNo(cur.IsRequired),
			IsPk:          yesOrNo(cur.IsPk),
			Sort:          cur.Sort,
			IsIncrement:   yesOrNo(cur.IsIncrement),
			ColumnType:    cur.ColumnType,
			ColumnComment: cur.ColumnComment.String,
			ColumnDefault: nullStr(cur.ColumnDefault),
		}
		cols = append(cols, t)
	}

	return &types.GetColumnsResponse{
		Columns: cols,
	}, nil
}

func yesOrNo(v string) string {
	if v == "1" {
		return "是"
	}
	return "-"
}

func nullStr(v sql.NullString) string {
	if v.Valid {
		if v.String == "" {
			return "''"
		}
		return v.String
	}
	return ""
}
