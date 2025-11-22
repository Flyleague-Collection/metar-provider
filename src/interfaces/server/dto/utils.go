// Package dto
package dto

import (
	"github.com/labstack/echo/v4"
)

var (
	ErrLackParam         = NewApiStatus("PARAM_MISS", "缺少参数", HttpCodeBadRequest)
	ErrInvalidParam      = NewApiStatus("PARAM_INVALID", "非法参数", HttpCodeBadRequest)
	ErrErrorParam        = NewApiStatus("PARAM_ERROR", "参数错误", HttpCodeBadRequest)
	ErrServerError       = NewApiStatus("SERVER_ERROR", "服务器错误", HttpCodeInternalError)
	SuccessHandleRequest = NewApiStatus("SUCCESS", "Success", HttpCodeOk)
)

func ErrorResponse(ctx echo.Context, codeStatus *ApiStatus) error {
	return NewApiResponse[any](codeStatus, nil).Response(ctx)
}

func TextResponse(ctx echo.Context, httpCode int, content string) error {
	return ctx.String(httpCode, content)
}
