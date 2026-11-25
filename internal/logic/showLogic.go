// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package logic

import (
	"context"
	"database/sql"
	"errors"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	ok, err := l.svcCtx.Filter.ExistsCtx(l.ctx, []byte(req.ShortURL))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("短链不存在")
	}
	// 从redis中取lurl
	lurl, err := l.svcCtx.ShortURLRedis.GetOrLoad(l.ctx, req.ShortURL, func(ctx context.Context) (string, error) {
		row, err := l.svcCtx.ShortURLModel.FindOneBySurl(ctx, sql.NullString{String: req.ShortURL, Valid: true})
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		if err != nil {
			logx.WithContext(ctx).Errorf("find short URL: %v", err)
			return "", err
		}
		return row.Lurl.String, nil
	})
	if err != nil {
		return nil, err
	}
	// 返回响应，在handler中进行重定向
	return &types.ShowResponse{
		LongURL: lurl,
	}, nil
}
