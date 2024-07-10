package webserver

import (
	"net/http"
)

type HandlerRestInterfaceUrlPrefix = interface {
	UrlPrefix() string
}

type HandlerRestInterface = interface {
	HandlerRestInterfaceUrlPrefix
	HandleInsert(http.ResponseWriter, *http.Request)
	HandleUpdate(http.ResponseWriter, *http.Request)
	HandleStore(http.ResponseWriter, *http.Request)
	HandleDelete(http.ResponseWriter, *http.Request)
	HandleListAll(http.ResponseWriter, *http.Request)
	HandleFindById(http.ResponseWriter, *http.Request)
	HandleFindByName(http.ResponseWriter, *http.Request)
}
