// Package dto
package dto

import "github.com/labstack/echo/v4"

type HttpCode int

const (
	HttpCodeOk               HttpCode = 200
	HttpCodeBadRequest       HttpCode = 400
	HttpCodeUnauthorized     HttpCode = 401
	HttpCodePermissionDenied HttpCode = 403
	HttpCodeNotFound         HttpCode = 404
	HttpCodeConflict         HttpCode = 409
	HttpCodeTooManyRequests  HttpCode = 429
	HttpCodeInternalError    HttpCode = 500
)

func (hc HttpCode) Code() int {
	return int(hc)
}

type ApiStatus struct {
	StatusName  string
	Description string
	HttpCode    HttpCode
}

func NewApiStatus(statusName, description string, httpCode HttpCode) *ApiStatus {
	return &ApiStatus{
		StatusName:  statusName,
		Description: description,
		HttpCode:    httpCode,
	}
}

type ApiResponse[T any] struct {
	HttpCode int    `json:"-"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Data     T      `json:"data"`
}

func (res *ApiResponse[T]) Response(ctx echo.Context) error {
	return ctx.JSON(res.HttpCode, res)
}

func NewApiResponse[T any](codeStatus *ApiStatus, data T) *ApiResponse[T] {
	return &ApiResponse[T]{
		HttpCode: codeStatus.HttpCode.Code(),
		Code:     codeStatus.StatusName,
		Message:  codeStatus.Description,
		Data:     data,
	}
}
