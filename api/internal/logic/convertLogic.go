// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/api/internal/svc"
	"shortener/api/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 处理长连接数据；数据不为空，长连接能请求通，是否已转链过，避免循环转链
	if ok := connect.Get(req.LongURL); !ok {
		return nil, errors.New("无效的链接")
	}

	md5Value := md5.Sum([]byte(req.LongURL))
	u, err := l.svcCtx.ShortURLModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})

	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该连接已被转为%s", u.Surl.String)
		}
		logx.Errorw("ShortURLModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	basePath, err := urltool.GetBasePath(req.LongURL)
	if err != nil {
		logx.Errorw("url.Parse failed", logx.LogField{Key: "lurl", Value: req.LongURL}, logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	_, err = l.svcCtx.ShortURLModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("该链接已是短链了")
		}
		logx.Errorw("ShortURLModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 从数据库中取号码
	shortURL := ""
	for {
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		// 号码转短链
		shortURL = base62.Int2String(seq)
		if _, ok := l.svcCtx.ShortBlackListMap[shortURL]; !ok {
			break
		}
	}

	err = l.svcCtx.Filter.AddCtx(l.ctx, []byte(shortURL))
	if err != nil {
		logx.Errorw("Filter.AddCtx failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, errors.New("过滤器映射失败")
	}

	// 存储长，短链接的映射
	if _, err := l.svcCtx.ShortURLModel.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongURL, Valid: true},
		Md5:  sql.NullString{String: md5Value, Valid: true},
		Surl: sql.NullString{String: shortURL, Valid: true},
	}); err != nil {
		logx.Errorw("ShortURLModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 返回短链接
	resURL := l.svcCtx.Config.ShortDomain + "/" + shortURL
	return &types.ConvertResponse{
		ShortURL: resURL,
	}, nil
}
