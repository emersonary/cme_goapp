package message

import (
	"fmt"
	"time"

	"github.com/emersonary/go-authentication/model/message"
	"github.com/emersonary/go-authentication/pck"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type MessageHandler struct {
	session *gocql.Session
}

func (m *MessageHandler) Insert(message *message.TMessage) (*message.TMessage, error) {

	message.CreatedAt = time.Now()

	if message.Id == uuid.Nil {
		message.Id = pck.NewID()
	}

	if err := m.session.Query(`insert into tbl_Message
	                 ( id
									 , fromid
									 , toid
									 , messagetext
									 , createdat
									 , readat )
									 values
									 ( ?, ?, ?, ?, ?, ?) `,
		message.Id.String(),
		message.FromUserId.String(),
		message.ToUserId.String(),
		message.MessageText,
		message.CreatedAt,
		message.ReadAt).Exec(); err != nil {

		return nil, err

	}

	return m.FindById(message.Id)

}

func (m *MessageHandler) UpdateMessagesReadAt(messages *[]message.TMessage, time time.Time) error {

	for i, message := range *messages {

		message.ReadAt = time
		(*messages)[i].ReadAt = time

		if err := m.session.Query(`Update tbl_Message
	                 set readat = ? 
									 where id = ?`,
			message.ReadAt,
			message.Id.String()).Exec(); err != nil {

			return err

		}

	}

	return nil

}
func (m *MessageHandler) FindById(id uuid.UUID) (*message.TMessage, error) {

	fmt.Println("FindById")
	iter := m.session.Query(`select id, fromid, toid, messagetext, createdat,readat from tbl_Message 
	                           where id  = ? `, id.String()).Iter()

	var idResult string
	var fromUserId string
	var toUserId string
	var createdat time.Time
	var readat time.Time
	var messageText string

	if !iter.Scan(&idResult, &fromUserId, &toUserId, &messageText, &createdat, &readat) {

		return nil, nil

	}

	fromUserUUID, err := pck.ParseID(fromUserId)

	if err != nil {

		return nil, err

	}

	toUserUUID, err := pck.ParseID(toUserId)

	if err != nil {

		return nil, err

	}

	message := message.NewMessage(fromUserUUID, toUserUUID, messageText)

	message.CreatedAt = createdat
	message.ReadAt = readat

	message.Id, err = pck.ParseID(idResult)

	if err != nil {

		return nil, err

	}

	return message, nil

}

func NewMessageHandler(session *gocql.Session) *MessageHandler {

	return &MessageHandler{session: session}

}
