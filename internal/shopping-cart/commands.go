package shoppingcart

import (
	"time"

	"github.com/google/uuid"
)

type Command interface {
	isEvent()
	GetType() string
}

func (e OpenShoppingCartCommand) isEvent()   {}
func (e CancelShoppingCartCommand) isEvent() {}

func (e OpenShoppingCartCommand) GetType() string {
	return e.Type
}
func (e CancelShoppingCartCommand) GetType() string {
	return e.Type
}

type CommandHeader struct {
	EventId   uuid.UUID `json:"event_id"`
	Type      string    `json:"type"`
	Kind      string    `json:"kind"`
	Timestamp time.Time `json:"timestamp"`
	StreamId  uuid.UUID `json:"stream_id"`
}

type CommandMetadata struct {
	IssuedBy    string `json:"issued_by"`
	IsTestEvent bool   `json:"is_test_event"`
}

type OpenShoppingCartCommandBody struct {
	Id       uuid.UUID `json:"id"`
	ClientId uuid.UUID `json:"client_id"`
}

type OpenShoppingCartCommand struct {
	*CommandHeader
	*CommandMetadata
	Data OpenShoppingCartCommandBody `json:"data"`
}

func NewOpenShoppingCartCommand() OpenShoppingCartCommand {
	eventId, _ := uuid.NewUUID()
	streamId, _ := uuid.NewUUID()
	clientId, _ := uuid.NewUUID()

	return OpenShoppingCartCommand{
		CommandHeader: &CommandHeader{
			EventId:   eventId,
			Type:      "open_new_shopping_cart",
			Kind:      "command",
			Timestamp: time.Now(),
			StreamId:  streamId,
		},
		CommandMetadata: &CommandMetadata{
			IsTestEvent: false,
			IssuedBy:    clientId.String(),
		},
		Data: OpenShoppingCartCommandBody{
			Id:       streamId,
			ClientId: clientId,
		},
	}
}

type CancelShoppingCartCommandBody struct {
	Id       uuid.UUID `json:"id"`
	ClientId uuid.UUID `json:"client_id"`
}

type CancelShoppingCartCommand struct {
	*CommandHeader
	*CommandMetadata
	Data CancelShoppingCartCommandBody `json:"data"`
}

func NewCancelShoppingCartCommand(id uuid.UUID) CancelShoppingCartCommand {
	eventId, _ := uuid.NewUUID()
	clientId, _ := uuid.NewUUID()

	return CancelShoppingCartCommand{
		CommandHeader: &CommandHeader{
			EventId:   eventId,
			Type:      "cancel_shopping_cart",
			Kind:      "command",
			Timestamp: time.Now(),
			StreamId:  id,
		},
		CommandMetadata: &CommandMetadata{
			IsTestEvent: false,
			IssuedBy:    clientId.String(),
		},
		Data: CancelShoppingCartCommandBody{
			Id:       id,
			ClientId: clientId,
		},
	}
}
