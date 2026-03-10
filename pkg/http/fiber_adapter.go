/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"bufio"
	"context"
	"io"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/logger"
)

type httpRequest struct {
	Ctx *fiber.Ctx
}

func NewHttpRequest(ctx *fiber.Ctx) *httpRequest {
	return &httpRequest{
		Ctx: ctx,
	}
}

func (request *httpRequest) GetParam(key string, defaultValue ...string) string {
	return request.Ctx.Params(key, defaultValue...)
}

func (request *httpRequest) GetParamInt(key string, defaultValue ...int) (int, error) {
	return request.Ctx.ParamsInt(key, defaultValue...)
}

func (request *httpRequest) GetBody() []byte {
	return request.Ctx.Body()
}

func (request *httpRequest) GetHeader(key string) string {
	return request.Ctx.Get(key)
}

func (request *httpRequest) GetHeaders() map[string]string {
	headers := make(map[string]string)
	for key, values := range request.Ctx.GetReqHeaders() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return headers
}

func (request *httpRequest) GetMethod() string {
	return request.Ctx.Method()
}

func (request *httpRequest) GetPath() string {
	return request.Ctx.Path()
}

func (request *httpRequest) GetQuery(key string, defaultValue ...string) string {
	return request.Ctx.Query(key, defaultValue...)
}

func (request *httpRequest) GetQueryInt(key string, defaultValue ...int) int {
	return request.Ctx.QueryInt(key, defaultValue...)

}

func (request *httpRequest) Queries() map[string]string {
	return request.Ctx.Queries()
}

func (request *httpRequest) GetIP() string {
	return request.Ctx.IP()
}

func (request *httpRequest) GetUserAgent() string {
	return string(request.Ctx.Context().UserAgent())
}

func (request *httpRequest) GetHost() string {
	return request.Ctx.Hostname()
}

func (request *httpRequest) ParseBody(any any) error {
	return request.Ctx.BodyParser(any)
}

func (request *httpRequest) File(key string) (*multipart.FileHeader, error) {
	return request.Ctx.FormFile(key)
}

func (request *httpRequest) FormData(key string) []string {
	values := request.Ctx.FormValue(key)
	return strings.Split(values, ",")
}

func (request *httpRequest) SaveFile(file *multipart.FileHeader, destination string) error {
	return request.Ctx.SaveFile(file, destination)
}

func (request *httpRequest) GetCookie(name string) string {
	cookie := request.Ctx.Cookies(name)
	if cookie == "" {
		return ""
	}
	return cookie
}

func (request *httpRequest) Context() context.Context {
	return request.Ctx.Context()
}

func (request *httpRequest) Validate(any any) error {
	if err := validate.Struct(any); err != nil {
		return ErrValidation(err)
	}
	return nil
}

func (request *httpRequest) ValidateBody(any any) error {
	if err := request.ParseBody(any); err != nil {
		logger.E(err, logger.Fields{"context": "ValidateBody", "action": "parse body"})
		return ErrBadRequest(err)
	}

	if err := validate.Struct(any); err != nil {
		logger.E(err, logger.Fields{"context": "ValidateBody", "action": "validation"})
		if errs, ok := err.(validator.ValidationErrors); ok {
			fields := make(map[string][]string)
			for _, e := range errs {
				fields[e.Field()] = append(fields[e.Field()], e.Tag())
			}

			return ErrValidation(err, fields)
		}

		return ErrValidation(err)
	}

	return nil
}

func (request *httpRequest) Auth() *vayload.HttpAuth {
	authToken := request.Ctx.Locals(HTTP_AUTH_KEY)
	auth, ok := authToken.(*vayload.HttpAuth)
	if !ok {
		return &vayload.HttpAuth{
			UserId:      snowflake.ID(0),
			Role:        "",
			AccessToken: "",
		}
	}

	return auth
}

func (request *httpRequest) GetLocal(key string) any {
	return request.Ctx.Locals(key)
}

func (request *httpRequest) Locals(key string, value any) any {
	if value != nil {
		request.Ctx.Locals(key, value)
		return value
	}

	return request.Ctx.Locals(key)
}

func (request *httpRequest) Next() error {
	return request.Ctx.Next()
}

type httpResponse struct {
	ctx *fiber.Ctx
}

func NewHttpResponse(ctx *fiber.Ctx) *httpResponse {
	return &httpResponse{
		ctx: ctx,
	}
}

func (response *httpResponse) SetStatus(status int) {
	response.ctx.Status(status)
}

func (response *httpResponse) SetHeader(key string, value string) {
	response.ctx.Set(key, value)
}

func (response *httpResponse) Send(body []byte) error {
	return response.ctx.Send(body)
}

func (response *httpResponse) JSON(data any) error {
	return response.ctx.JSON(data)
}

func (response *httpResponse) Json(body any) error {
	return response.ctx.JSON(body)
}

func (response *httpResponse) File(path string) error {
	return response.ctx.SendFile(path)
}

func (response *httpResponse) Stream(stream io.Reader) error {
	return response.ctx.SendStream(stream)
}

func (response *httpResponse) Status(status int) vayload.HttpResponse {
	response.ctx.Status(status)
	return response
}

func (response *httpResponse) Redirect(path string, status int) error {
	return response.ctx.Redirect(path, status)
}

func (response *httpResponse) SetBodyStreamWriter(writer vayload.StreamWriter) error {
	response.ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		writer(w)
	}))

	return nil
}

func (response *httpResponse) Cookie(cookie *vayload.Cookie) vayload.HttpResponse {
	fiberCookie := fiber.Cookie{
		Name:        cookie.Name,
		Value:       cookie.Value,
		Path:        cookie.Path,
		Domain:      cookie.Domain,
		MaxAge:      cookie.MaxAge,
		Expires:     cookie.Expires,
		Secure:      cookie.Secure,
		HTTPOnly:    cookie.HttpOnly,
		SameSite:    cookie.SameSite,
		SessionOnly: cookie.SessionOnly,
	}
	response.ctx.Cookie(&fiberCookie)
	return response
}

func (response *httpResponse) Cookies(cookies ...*vayload.Cookie) vayload.HttpResponse {
	for _, c := range cookies {
		fiberCookie := fiber.Cookie{
			Name:        c.Name,
			Value:       c.Value,
			Path:        c.Path,
			Domain:      c.Domain,
			MaxAge:      c.MaxAge,
			Expires:     c.Expires,
			Secure:      c.Secure,
			HTTPOnly:    c.HttpOnly,
			SameSite:    c.SameSite,
			SessionOnly: c.SessionOnly,
		}
		response.ctx.Cookie(&fiberCookie)
	}
	return response
}

func FiberWrap(handler vayload.HttpHandler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := NewHttpRequest(ctx)
		res := NewHttpResponse(ctx)

		if err := handler(req, res); err != nil {
			return err
		}

		return nil
	}
}
