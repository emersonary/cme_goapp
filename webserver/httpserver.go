package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/emersonary/go-authentication/config"
	"github.com/emersonary/go-authentication/event/prometheus"
	webserver "github.com/emersonary/go-authentication/webserver/restapi"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/go-chi/chi/middleware"
)

var RootRouter *chi.Mux

func AddServerHandler(handleInterfaceUrlPrefix webserver.HandlerRestInterfaceUrlPrefix, tokenAuth *jwtauth.JWTAuth) *chi.Router {

	var result *chi.Router = nil

	RootRouter.Route("/"+handleInterfaceUrlPrefix.UrlPrefix(), func(router chi.Router) {

		if tokenAuth != nil {
			router.Use(jwtauth.Verifier(tokenAuth))
			router.Use(jwtauth.Authenticator)
		}

		if handleInterface, ok := handleInterfaceUrlPrefix.(webserver.HandlerRestInterface); ok {

			router.Post("/insert", handleInterface.HandleInsert)
			router.Get("/findbyid/{id}", handleInterface.HandleFindById)
			router.Get("/findbyname/{name}", handleInterface.HandleFindByName)

		}

		result = &router

	})

	return result

}

func InitMiddlewares(config *config.Conf) {

	RootRouter.Use(middleware.WithValue("jwt", config.TokenAuth))
	RootRouter.Use(middleware.WithValue("jwtExpiresIn", config.JWTExpiresIn))

}

func StartWebServer(config *config.Conf) error {

	fmt.Println("Starting Web Server at port " + strconv.Itoa(config.WebServerPort))
	return http.ListenAndServe(":"+strconv.Itoa(config.WebServerPort), RootRouter)

}

func init() {

	RootRouter = chi.NewRouter()
	RootRouter.Use(middleware.Logger)
	RootRouter.Use(middleware.Recoverer)
	RootRouter.Use(prometheus.Logger)

}
