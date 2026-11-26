// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package handler

import (
	"net/http"

	"shortener/api/internal/logic"
	"shortener/api/internal/svc"
	"shortener/api/internal/types"

	"github.com/go-playground/validator"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			logx.Errorw("validator check failed", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
