package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	httpcontract "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
)

type FiberGroup struct {
	instance          *fiber.App
	originPrefix      string
	prefix            string
	originMiddlewares []httpcontract.Middleware
	middlewares       []httpcontract.Middleware
	lastMiddlewares   []httpcontract.Middleware
}

func NewFiberGroup(instance *fiber.App, prefix string, originMiddlewares []httpcontract.Middleware, lastMiddlewares []httpcontract.Middleware) route.Route {
	return &FiberGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		lastMiddlewares:   lastMiddlewares,
	}
}

func (r *FiberGroup) Group(handler route.GroupFunc) {
	var middlewares []httpcontract.Middleware
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []httpcontract.Middleware{}
	prefix := pathToFiberPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	handler(NewFiberGroup(r.instance, prefix, middlewares, r.lastMiddlewares))
}

func (r *FiberGroup) Prefix(addr string) route.Route {
	r.prefix += "/" + addr

	return r
}

func (r *FiberGroup) Middleware(middlewares ...httpcontract.Middleware) route.Route {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

func (r *FiberGroup) Any(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().All(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Get(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Get(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Post(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Post(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Delete(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Delete(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Patch(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Patch(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Put(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Put(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Options(relativePath string, handler httpcontract.HandlerFunc) {
	r.getFiberRoutesWithMiddlewares().Options(pathToFiberPath(relativePath), r.handlerToFiberHandler(handler)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Resource(relativePath string, controller httpcontract.ResourceController) {
	r.getFiberRoutesWithMiddlewares().Get(pathToFiberPath(relativePath), r.handlerToFiberHandler(controller.Index)...)
	r.getFiberRoutesWithMiddlewares().Post(pathToFiberPath(relativePath), r.handlerToFiberHandler(controller.Store)...)
	r.getFiberRoutesWithMiddlewares().Get(pathToFiberPath(relativePath+"/{id}"), r.handlerToFiberHandler(controller.Show)...)
	r.getFiberRoutesWithMiddlewares().Put(pathToFiberPath(relativePath+"/{id}"), r.handlerToFiberHandler(controller.Update)...)
	r.getFiberRoutesWithMiddlewares().Patch(pathToFiberPath(relativePath+"/{id}"), r.handlerToFiberHandler(controller.Update)...)
	r.getFiberRoutesWithMiddlewares().Delete(pathToFiberPath(relativePath+"/{id}"), r.handlerToFiberHandler(controller.Destroy)...)
	r.clearMiddlewares()
}

func (r *FiberGroup) Static(relativePath, root string) {
	r.getFiberRoutesWithMiddlewares().Static(pathToFiberPath(relativePath), root)
	r.clearMiddlewares()
}

func (r *FiberGroup) StaticFile(relativePath, filepath string) {
	r.getFiberRoutesWithMiddlewares().Static(pathToFiberPath(relativePath), filepath)
	r.clearMiddlewares()
}

func (r *FiberGroup) StaticFS(relativePath string, fs http.FileSystem) {
	r.getFiberRoutesWithMiddlewares().Use(pathToFiberPath(relativePath), filesystem.New(filesystem.Config{
		Root: fs,
	}))
	r.clearMiddlewares()
}

func (r *FiberGroup) middlewaresList() []fiber.Handler {
	var middlewares []fiber.Handler
	ginOriginMiddlewares := middlewaresToFiberHandlers(r.originMiddlewares)
	ginMiddlewares := middlewaresToFiberHandlers(r.middlewares)
	ginLastMiddlewares := middlewaresToFiberHandlers(r.lastMiddlewares)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginLastMiddlewares...)

	return middlewares
}

func (r *FiberGroup) getFiberRoutesWithMiddlewares() fiber.Router {
	prefix := pathToFiberPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	return r.instance.Group(prefix)
}

func (r *FiberGroup) clearMiddlewares() {
	r.middlewares = []httpcontract.Middleware{}
}
