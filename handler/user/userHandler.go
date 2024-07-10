package user

import (
	"errors"
	"time"

	"github.com/emersonary/go-authentication/model/user"
	"github.com/emersonary/go-authentication/pck"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type UserHandler struct {
	session *gocql.Session
}

func (u *UserHandler) Insert(user *user.TUser) (*user.TUser, error) {

	existingUser, err := u.FindByName(user.Name)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("User name (" + user.Name + ") already exists!")
	}

	user.CreatedAt = time.Now()

	if user.Id == uuid.Nil {
		user.Id = pck.NewID()
	}

	u.session.SetConsistency(gocql.Quorum)
	if err := u.session.Query(`insert into tbl_User 
	                 ( id
									 , name
									 , email
									 , password
									 , phone
									, createdat )
									 values
									 ( ?, ?, ?, ?, ?,?) `,
		user.Id.String(),
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.CreatedAt).Exec(); err != nil {

		return nil, err

	}

	return user, nil

}

func (u *UserHandler) FindById(id uuid.UUID) (*user.TUser, error) {

	u.session.SetConsistency(gocql.Quorum)
	iter := u.session.Query(`select id, name, phone, email, password,createdat from tbl_User 
	                           where id  = ? `, id.String()).Iter()

	var idResult string
	var nameResult string
	var phone string
	var email string
	var password string
	var createdat time.Time

	if !iter.Scan(&idResult, &nameResult, &phone, &email, &password, &createdat) {

		return nil, nil

	}

	user, err := user.NewUser(nameResult, email, phone, "")

	if err != nil {

		return nil, err

	}

	user.Id, err = pck.ParseID(idResult)

	if err != nil {

		return nil, err

	}

	user.Password = password
	user.CreatedAt = createdat

	return user, nil

}

// Desenvolver
func (u *UserHandler) DeleteById(id uuid.UUID) error {

	u.session.SetConsistency(gocql.Quorum)
	iter := u.session.Query(`select id, name, phone, email, password,createdat from tbl_User 
	                           where id  = ? `, id.String()).Iter()

	var idResult string
	var nameResult string
	var phone string
	var email string
	var password string
	var createdat time.Time

	if !iter.Scan(&idResult, &nameResult, &phone, &email, &password, &createdat) {

		return nil

	}

	user, err := user.NewUser(nameResult, email, phone, "")

	if err != nil {

		return err

	}

	user.Id, err = pck.ParseID(idResult)

	if err != nil {

		return err

	}

	user.Password = password
	user.CreatedAt = createdat

	return nil

}

func (u *UserHandler) FindByName(name string) (*user.TUser, error) {

	u.session.SetConsistency(gocql.Quorum)
	iter := u.session.Query(`select id, name, phone, email, password from tbl_User 
	                           where name  = ? `, name).Iter()

	var idResult string
	var nameResult string
	var phone string
	var email string
	var password string

	if !iter.Scan(&idResult, &nameResult, &phone, &email, &password) {

		return nil, nil

	}

	user, err := user.NewUser(nameResult, email, phone, "")

	if err != nil {

		return nil, err

	}

	user.Id, err = pck.ParseID(idResult)

	if err != nil {

		return nil, err

	}

	user.Password = password

	return user, nil

}

func NewUserHandler(session *gocql.Session) *UserHandler {

	return &UserHandler{session: session}

}
