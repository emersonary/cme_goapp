package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	rediscache "github.com/emersonary/go-authentication/cache/redis"
	errorsclass "github.com/emersonary/go-authentication/error"
	"github.com/emersonary/go-authentication/handler/message"
	"github.com/emersonary/go-authentication/handler/user"
	model "github.com/emersonary/go-authentication/model/message"
	"github.com/emersonary/go-authentication/pck"
	"github.com/emersonary/go-authentication/webserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/gocql/gocql"
)

type MessageHttpHandler struct {
	HandlerMessage *message.MessageHandler
	userHandler    *user.UserHandler
}

func (h *MessageHttpHandler) UrlPrefix() string {
	return "message"
}

func NewHttpHandlerMessage(messageHandler *message.MessageHandler, userHandler *user.UserHandler) *MessageHttpHandler {
	return &MessageHttpHandler{HandlerMessage: messageHandler, userHandler: userHandler}
}

func CreateMessageDependency(
	interfaceSession *gocql.Session,
	tokenAuth *jwtauth.JWTAuth,
	userHandler *user.UserHandler) {

	messageHandler := message.NewMessageHandler(interfaceSession)
	messageHttpHandler := NewHttpHandlerMessage(messageHandler, userHandler)
	router := webserver.AddServerHandler(messageHttpHandler, tokenAuth)

	(*router).Post("/send", messageHttpHandler.HandleInsert)
	(*router).Get("/findbyid/{id}", messageHttpHandler.HandleFindById)
	(*router).Get("/inbox/{userid}", messageHttpHandler.HandleInboxByUserId)

	webserver.RootRouter.Post("/send", messageHttpHandler.HandleInsert)
	webserver.RootRouter.Post("/messages/{userid}", messageHttpHandler.HandleInboxByUserId)

}

func (h *MessageHttpHandler) HandleInsert(w http.ResponseWriter, r *http.Request) {

	message := model.NewMessageEmpty()

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		error := errorsclass.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Inserting Message " + message.FromUserId.String())

	user, err := h.userHandler.FindById(message.FromUserId)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		json.NewEncoder(w).Encode("FromUserId not found in database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err = h.userHandler.FindById(message.ToUserId)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		json.NewEncoder(w).Encode("ToUserId not found in database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messageInserted, err := h.HandlerMessage.Insert(message)

	if err != nil {
		error := errorsclass.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rediscache.RedisCtrl.AddCache(messageInserted)

	if err != nil {
		error := errorsclass.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messageInserted)
	w.WriteHeader(http.StatusCreated)

}

func (h *MessageHttpHandler) HandleFindById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		error := errorsclass.Error{Message: "Error retrieving ID param"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)

		return
	}

	idUUID, err := pck.ParseID(id)

	if err != nil {

		fmt.Println(err)
		fmt.Println("error")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())

		return

	}

	message, err := h.HandlerMessage.FindById(idUUID)

	if err != nil {

		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return

	}

	if message != nil {
		json.NewEncoder(w).Encode(message)
		w.WriteHeader(http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Message not Found!")
	}

}

func (h *MessageHttpHandler) HandleInboxByUserId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "userid")

	if id == "" {
		error := errorsclass.Error{Message: "Error retrieving ID param"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)

		return
	}

	idUUID, err := pck.ParseID(id)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())

		return

	}

	messages, err := rediscache.RedisCtrl.CachedMessages(idUUID)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())

		return

	}

	err = rediscache.RedisCtrl.DeleteFirstMessages(idUUID, len(*messages))

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())

		return

	}

	h.HandlerMessage.UpdateMessagesReadAt(messages, time.Now())

	json.NewEncoder(w).Encode(messages)
	w.WriteHeader(http.StatusFound)

}
