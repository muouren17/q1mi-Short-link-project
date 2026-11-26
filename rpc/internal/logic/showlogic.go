package logic

import (
	"context"
	"database/sql"
	"errors"

	"shortener/rpc/internal/svc"
	"shortener/rpc/types/shortener/rpc/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShowLogic) Show(in *types.ShowRequest) (*types.ShowResponse, error) {
	ok, err := l.svcCtx.Filter.ExistsCtx(l.ctx, []byte(in.ShortURL))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("短链不存在")
	}
	// 从redis中取lurl
	lurl, err := l.svcCtx.ShortURLRedis.GetOrLoad(l.ctx, in.ShortURL, func(ctx context.Context) (string, error) {
		row, err := l.svcCtx.ShortURLModel.FindOneBySurl(ctx, sql.NullString{String: in.ShortURL, Valid: true})
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
	return &types.ShowResponse{LongURL: lurl}, nil
}
