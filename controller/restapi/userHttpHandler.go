package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emersonary/go-authentication/dto"
	"github.com/emersonary/go-authentication/error"
	handler "github.com/emersonary/go-authentication/handler/user"
	"github.com/emersonary/go-authentication/model/user"
	"github.com/emersonary/go-authentication/pck"
	"github.com/emersonary/go-authentication/webserver"
	restapi "github.com/emersonary/go-authentication/webserver/restapi"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/gocql/gocql"
)

type UserHttpHandler struct {
	restapi.HandlerRestInterface
	HandlerUser *handler.UserHandler
}

type authHttpHandler struct {
	HandlerUser *handler.UserHandler
}

func (u *UserHttpHandler) UrlPrefix() string {
	return "user"
}

func (u *authHttpHandler) UrlPrefix() string {
	return "login"
}

func NewHttpHandlerUser(userHandler *handler.UserHandler) *UserHttpHandler {
	return &UserHttpHandler{HandlerUser: userHandler}
}

func newHttpHandlerAuth(userHandler *handler.UserHandler) *authHttpHandler {
	return &authHttpHandler{HandlerUser: userHandler}
}

func CreateUserDependency(interfaceSession *gocql.Session, tokenAuth *jwtauth.JWTAuth) *handler.UserHandler {

	userHandler := handler.NewUserHandler(interfaceSession)
	userHttpHandler := NewHttpHandlerUser(userHandler)
	webserver.AddServerHandler(userHttpHandler, tokenAuth)

	authHttpHandler := newHttpHandlerAuth(userHandler)
	router := webserver.AddServerHandler(authHttpHandler, nil)

	(*router).Post("/", authHttpHandler.HandleAuth)

	webserver.RootRouter.Post("/register", userHttpHandler.HandleInsert)

	return userHandler

}

func (h *UserHttpHandler) HandleInsert(w http.ResponseWriter, r *http.Request) {

	userDTO := user.NewUserDTOEmpty()

	err := json.NewDecoder(r.Body).Decode(&userDTO)

	if err != nil {
		error := error.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userDTO.TUserFrom()

	if err != nil {
		error := error.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userInserted, err := h.HandlerUser.Insert(user)
	if err != nil {
		error := error.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userSelected, err := h.HandlerUser.FindById(userInserted.Id)
	if err != nil {
		error := error.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userSelected.Password = ""

	json.NewEncoder(w).Encode(userSelected)
	w.WriteHeader(http.StatusCreated)

}

func (h *UserHttpHandler) HandleFindById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		error := error.Error{Message: "Error retrieving ID param"}
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

	user, err := h.HandlerUser.FindById(idUUID)

	if err != nil {

		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return

	}

	if user != nil {
		user.Password = ""
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not Found!")
	}

}

func (h *UserHttpHandler) HandleFindByName(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	if name == "" {

		error := error.Error{Message: "Error retrieving ID param"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)

		return
	}

	user, err := h.HandlerUser.FindByName(name)

	if err != nil {

		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return

	}

	if user != nil {
		user.Password = ""
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not Found!")
	}

}

func (h *authHttpHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {

	userAuth := &user.TUserAuth{}

	err := json.NewDecoder(r.Body).Decode(&userAuth)

	if err != nil {

		error := error.Error{Message: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)

		return
	}

	if userAuth.Name == "" || userAuth.Password == "" {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid body")
		return

	}

	user, err := h.HandlerUser.FindByName(userAuth.Name)

	if err != nil {

		error := error.Error{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(error)

		return
	}

	if user == nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid credentials")

		return

	}

	if !user.ValidatePassword(userAuth.Password) {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid credentials")

		return

	}

	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.Id.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&accessToken)

}
