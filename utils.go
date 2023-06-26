package fiber

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unsafe"

	"github.com/gofiber/fiber/v2"

	"github.com/goravel/framework/contracts/config"
	httpcontract "github.com/goravel/framework/contracts/http"
)

func pathToFiberPath(relativePath string) string {
	return bracketToColon(mergeSlashForPath(relativePath))
}

func middlewaresToFiberHandlers(middlewares []httpcontract.Middleware) []fiber.Handler {
	var fiberHandlers []fiber.Handler
	for _, item := range middlewares {
		fiberHandlers = append(fiberHandlers, middlewareToFiberHandler(item))
	}

	return fiberHandlers
}

func middlewaresToFiberAny(middlewares []httpcontract.Middleware) []any {
	var fiberHandlers []any
	for _, item := range middlewares {
		fiberHandlers = append(fiberHandlers, middlewareToFiberHandler(item))
	}

	return fiberHandlers
}

func (r *FiberGroup) handlerToFiberHandler(handler httpcontract.HandlerFunc) []fiber.Handler {
	handlers := r.middlewaresList()
	handlers = append(handlers, func(fiberCtx *fiber.Ctx) error {
		handler(NewFiberContext(fiberCtx))

		return nil
	})

	return handlers
}

func (r *FiberGroup) handlerToFiberHandlers(handler httpcontract.HandlerFunc) []fiber.Handler {
	handlers := r.middlewaresList()
	handlers = append(handlers, func(fiberCtx *fiber.Ctx) error {
		handler(NewFiberContext(fiberCtx))

		return nil
	})

	return handlers
}

func handlerToFiberHandler(handler httpcontract.HandlerFunc) fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		handler(NewFiberContext(fiberCtx))
		return nil
	}
}

func middlewareToFiberHandler(handler httpcontract.Middleware) fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		return handler(NewFiberContext(fiberCtx))
	}
}

func getDebugLog(config config.Config) fiber.Handler {
	// TODO: add debug log

	return nil
}

func colonToBracket(relativePath string) string {
	arr := strings.Split(relativePath, "/")
	var newArr []string
	for _, item := range arr {
		if strings.HasPrefix(item, ":") {
			item = "{" + strings.ReplaceAll(item, ":", "") + "}"
		}
		newArr = append(newArr, item)
	}

	return strings.Join(newArr, "/")
}

func bracketToColon(relativePath string) string {
	compileRegex := regexp.MustCompile(`{(.*?)}`)
	matchArr := compileRegex.FindAllStringSubmatch(relativePath, -1)

	for _, item := range matchArr {
		relativePath = strings.ReplaceAll(relativePath, item[0], ":"+item[1])
	}

	return relativePath
}

func mergeSlashForPath(path string) string {
	path = strings.ReplaceAll(path, "//", "/")

	return strings.ReplaceAll(path, "//", "/")
}

func runningInConsole() bool {
	args := os.Args

	return len(args) >= 2 && args[1] == "artisan"
}

func handleException(ctx *FiberContext, err error) {
	ExceptionFacade.Report(err)
	ExceptionFacade.Render(ctx, err)
}

// ConvertRequest converts a fasthttp.Request to an http.Request.
// forServer should be set to true when the http.Request is going to be passed to a http.Handler.
//
// The http.Request must not be used after the fasthttp handler has returned!
// Memory in use by the http.Request will be reused after your handler has returned!
func ConvertRequest(ctx *fasthttp.RequestCtx, r *http.Request, forServer bool) error {
	body := ctx.PostBody()
	strRequestURI := b2s(ctx.RequestURI())

	rURL, err := url.ParseRequestURI(strRequestURI)
	if err != nil {
		return err
	}

	r.Method = b2s(ctx.Method())
	r.Proto = b2s(ctx.Request.Header.Protocol())
	if r.Proto == "HTTP/2" {
		r.ProtoMajor = 2
	} else {
		r.ProtoMajor = 1
	}
	r.ProtoMinor = 1
	r.ContentLength = int64(len(body))
	r.RemoteAddr = ctx.RemoteAddr().String()
	r.Host = b2s(ctx.Host())
	r.TLS = ctx.TLSConnectionState()
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.URL = rURL

	if forServer {
		r.RequestURI = strRequestURI
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	} else if len(r.Header) > 0 {
		for k := range r.Header {
			delete(r.Header, k)
		}
	}

	ctx.Request.Header.VisitAll(func(k, v []byte) {
		sk := b2s(k)
		sv := b2s(v)

		switch sk {
		case "Transfer-Encoding":
			r.TransferEncoding = append(r.TransferEncoding, sv)
		default:
			r.Header.Set(sk, sv)
		}
	})

	return nil
}

func b2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}
