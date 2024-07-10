package message

import (
	"time"

	"github.com/emersonary/go-authentication/model/base"
	"github.com/emersonary/go-authentication/pck"
	"github.com/google/uuid"
)

type TMessage struct {
	base.TBase
	FromUserId  pck.UUID  `json:"fromuserid"`
	ToUserId    pck.UUID  `json:"touserid"`
	ReadAt      time.Time `json:"readat"`
	MessageText string    `json:"messagetext"`
}

func NewMessage(fromUserId pck.UUID, toUserId pck.UUID, messageText string) *TMessage {

	return &TMessage{
		TBase: base.TBase{
			Id: uuid.Nil,
		},
		FromUserId:  fromUserId,
		ToUserId:    toUserId,
		MessageText: messageText}
}

func NewMessageEmpty() *TMessage {

	return &TMessage{}
}
