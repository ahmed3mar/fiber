package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/goravel/framework/contracts/exception"
	"github.com/goravel/framework/contracts/translation"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/validation"
)

const HttpBinding = "goravel.fiber.http"
const RouteBinding = "goravel.fiber.route"

var App foundation.Application

var (
	ConfigFacade      config.Config
	CacheFacade       cache.Cache
	LogFacade         log.Log
	RateLimiterFacade http.RateLimiter
	ValidationFacade  validation.Validation
	TranslationFacade translation.Translation
	ExceptionFacade   exception.Exception
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(HttpBinding, func(app foundation.Application) (any, error) {
		return NewFiberContext(&fiber.Ctx{}), nil
	})
	app.Bind(RouteBinding, func(app foundation.Application) (any, error) {
		return NewFiberRoute(app.MakeConfig()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	ConfigFacade = app.MakeConfig()
	CacheFacade = app.MakeCache()
	LogFacade = app.MakeLog()
	ValidationFacade = app.MakeValidation()
	TranslationFacade = app.MakeTranslation()
	ExceptionFacade = app.MakeException()
}
