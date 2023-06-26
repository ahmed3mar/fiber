package fiber

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"

	httpcontract "github.com/goravel/framework/contracts/http"
)

type FiberResponse struct {
	instance *fiber.Ctx
	origin   httpcontract.ResponseOrigin
}

func (r *FiberResponse) Status(code int) httpcontract.ResponseStatus {
	return NewFibertatus(r.instance, code)
}

type FiberStatus struct {
	instance *fiber.Ctx
	status   int
}

func NewFibertatus(instance *fiber.Ctx, code int) httpcontract.ResponseSuccess {
	return &FiberStatus{instance, code}
}

func (r *FiberStatus) Data(contentType string, data []byte) {
	_ = r.instance.Status(r.status).Type(contentType).Send(data)
}

func (r *FiberStatus) Json(obj any) {
	r.instance.Status(r.status).JSON(obj)
}

func (r *FiberStatus) String(format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(r.status).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(r.status).Type(format).SendString(values[0].(string))
}

func NewFiberResponse(instance *fiber.Ctx, origin httpcontract.ResponseOrigin) *FiberResponse {
	return &FiberResponse{instance, origin}
}

func (r *FiberResponse) Data(code int, contentType string, data []byte) {
	_ = r.instance.Status(code).Type(contentType).Send(data)
}

func (r *FiberResponse) Download(filepath, filename string) {
	_ = r.instance.Download(filepath, filename)
}

func (r *FiberResponse) File(filepath string) {
	_ = r.instance.SendFile(filepath)
}

func (r *FiberResponse) Header(key, value string) httpcontract.Response {
	r.instance.Set(key, value)

	return r
}

func (r *FiberResponse) Json(code int, obj any) {
	_ = r.instance.Status(code).JSON(obj)
}

func (r *FiberResponse) Origin() httpcontract.ResponseOrigin {
	return r.origin
}

func (r *FiberResponse) Redirect(code int, location string) {
	_ = r.instance.Redirect(location, code)
}

func (r *FiberResponse) String(code int, format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(code).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(code).Type(format).SendString(values[0].(string))
}

func (r *FiberResponse) Success() httpcontract.ResponseSuccess {
	return NewFiberSuccess(r.instance)
}

func (r *FiberResponse) Writer() http.ResponseWriter {
	// Fiber doesn't support this
	return nil
}

type FiberSuccess struct {
	instance *fiber.Ctx
}

func NewFiberSuccess(instance *fiber.Ctx) httpcontract.ResponseSuccess {
	return &FiberSuccess{instance}
}

func (r *FiberSuccess) Data(contentType string, data []byte) {
	_ = r.instance.Type(contentType).Send(data)
}

func (r *FiberSuccess) Json(obj any) {
	_ = r.instance.Status(http.StatusOK).JSON(obj)
}

func (r *FiberSuccess) String(format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(http.StatusOK).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(http.StatusOK).Type(format).SendString(values[0].(string))
}

func FiberResponseMiddleware() httpcontract.Middleware {
	return func(ctx httpcontract.Context) error {
		blw := &BodyWriter{body: bytes.NewBufferString("")}
		switch ctx := ctx.(type) {
		case *FiberContext:
			// TODO: implement
			blw.Writer = ctx.Instance().Response().BodyWriter()
		}

		ctx.WithValue("responseOrigin", blw)
		return ctx.Request().Next()
	}
}

type BodyWriter struct {
	io.Writer
	body *bytes.Buffer
}

func (w *BodyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (w *BodyWriter) Size() int {
	//TODO implement me
	panic("implement me")
}

func (w *BodyWriter) Status() int {
	//TODO implement me
	panic("implement me")
}

func (w *BodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.Writer.Write(b)
}

func (w *BodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)

	return w.Writer.Write([]byte(s))
}

func (w *BodyWriter) Body() *bytes.Buffer {
	return w.Writer.(*bytes.Buffer)
}
