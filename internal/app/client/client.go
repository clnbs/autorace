package client

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/app/server"
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/clnbs/autorace/pkg/logger"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

//AutoraceClient handle client connection to a RabbitMQ server
type AutoraceClient struct {
	SessionID        uuid.UUID
	playerName       string
	playerUUID       uuid.UUID
	partyUUID        uuid.UUID
	rabbitConnection *messaging.RabbitConnection
}

//NewAutoraceClient return a AutoraceClient with a given name
func NewAutoraceClient(name string, rabbitConfig messaging.RabbitConnectionConfiguration) (*AutoraceClient, error) {
	var err error
	arClient := new(AutoraceClient)
	arClient.SessionID = uuid.New()
	arClient.playerName = name
	arClient.rabbitConnection, err = messaging.NewRabbitConnection(rabbitConfig)
	return arClient, err
}

//RequestPlayerCreation handle player creation with server via a RabbitMQ connection
func (arClient *AutoraceClient) RequestPlayerCreation(readyToReceive chan bool) (*models.Player, error) {
	//object received by the handler are pass through this "received" chan
	received := make(chan interface{})
	playerRequest := models.PlayerCreationToken{
		SessionUUID: arClient.SessionID,
		PlayerName:  arClient.playerName,
	}
	// start object handler from server
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopic("autocar.player.creation."+arClient.SessionID.String(), arClient.computePlayerCreationResponse, received, readyToReceive)
		if err != nil {
			logger.Error("error while receiving Player creation from server :", err)
			return
		}
	}()
	// wait for receiver to be started
	if !<-readyToReceive {
		return nil, errors.New("could not received message on topic")
	}
	go func() {
		arClient.rabbitConnection.SendMessageOnTopic(playerRequest, "autocar.player.creation")
	}()
	response := <-received
	// check response type
	if response == nil {
		return nil, errors.New("unable to receive response to Player creation token")
	}
	if reflect.TypeOf(response) == reflect.TypeOf(messaging.ErrorResponse{}) {
		return nil, errors.New(response.(messaging.ErrorResponse).ErrorMessage)
	}
	// assign player UUID to this object for further usage
	arClient.playerUUID = response.(*models.Player).PlayerUUID
	return response.(*models.Player), nil
}

func (arClient *AutoraceClient) computePlayerCreationResponse(msg []byte) interface{} {
	playerInResponse := new(models.Player)
	err := json.Unmarshal(msg, playerInResponse)
	if err != nil {
		serverError := new(messaging.ErrorResponse)
		err = json.Unmarshal(msg, serverError)
		if err != nil {
			return messaging.ErrorResponse{ErrorMessage: "unable to read response from server :" + err.Error()}
		}
		return serverError
	}
	return playerInResponse
}

//RequestPartyCreation handle party creation with a static server instance via a RabbitMQ connection
func (arClient *AutoraceClient) RequestPartyCreation(partyConfig models.PartyCreationToken, readyToReceive chan bool) (*models.Party, error) {
	received := make(chan interface{})
	// start the handling function before sending a request
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopic("autocar.party.creation."+arClient.playerUUID.String(), arClient.computePartyCreationResponse, received, readyToReceive)
		if err != nil {
			logger.Error("error while receiving Party creation from server :", err)
			return
		}
	}()
	// waiting for listener to be started
	<-readyToReceive
	// send the actual request
	go func() {
		arClient.rabbitConnection.SendMessageOnTopic(partyConfig, "autocar.party.creation")
	}()
	// waiting response
	response := <-received
	// handle response
	if response == nil {
		return nil, errors.New("unable to receive response to Party creation token")
	}
	if reflect.TypeOf(response) == reflect.TypeOf(messaging.ErrorResponse{}) {
		return nil, errors.New(response.(messaging.ErrorResponse).ErrorMessage)
	}
	// store party UUID for further usage
	arClient.partyUUID = response.(*models.Party).PartyUUID
	return response.(*models.Party), nil
}

func (arClient *AutoraceClient) computePartyCreationResponse(msg []byte) interface{} {
	newParty := new(models.Party)
	err := json.Unmarshal(msg, newParty)
	if err != nil {
		newServerError := new(messaging.ErrorResponse)
		err = json.Unmarshal(msg, newServerError)
		if err != nil {
			return messaging.ErrorResponse{ErrorMessage: err.Error()}
		}
		return newServerError
	}
	return newParty
}

// ReceiveParty handle party sent by a dynamic server instance
func (arClient *AutoraceClient) ReceiveParty(partyID string, readyToReceive chan bool) (*models.Party, error) {
	received := make(chan interface{})
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopicWithHeader(
			"autocar.party."+partyID+".map."+arClient.playerUUID.String(),
			arClient.mapHandler,
			received,
			readyToReceive,
		)
		if err != nil {
			logger.Error("while receiving a party :", err)
			return
		}
	}()
	<-readyToReceive
	response := <-received
	switch response.(type) {
	case error:
		return nil, response.(error)
	case *models.Party:
		arClient.partyUUID = response.(*models.Party).PartyUUID
		return response.(*models.Party), nil
	}
	readyToReceive <- true
	return nil, nil
}

func (arClient *AutoraceClient) mapHandler(msg amqp.Delivery) interface{} {
	updatedParty := new(models.Party)
	err := json.Unmarshal(msg.Body, updatedParty)
	if err != nil {
		logger.Error("error while trying to unmarshal sync from server :", err)
		return err
	}
	return updatedParty
}

// SendPlayerInput is use to send input to a dynamic server instance.
func (arClient *AutoraceClient) SendPlayerInput(pInput *models.PlayerInput) {
	arClient.rabbitConnection.SendMessageOnTopic(pInput, "autocar.party."+arClient.partyUUID.String()+".input")
}

// AddPlayerRequest handle adding player. It send a request to a dynamic server instance.
// Dynamic server instance respond by sending to the client the already-registered-client list
// as competitor and the party's content. On the server side, the dynamic instance store this client
// as a game participant
func (arClient *AutoraceClient) AddPlayerRequest(partyID string) error {
	addPlayerToken := models.PlayerToken{
		ClientID: arClient.playerUUID.String(),
		PartyID:  partyID,
	}
	arClient.rabbitConnection.SendMessageOnTopic(addPlayerToken, "autocar.party."+partyID+".addPlayer")
	return nil
}

// ReceiveSync receive sync message from a dynamic server instance, filter error message and
// send sync message to a given chan
func (arClient *AutoraceClient) ReceiveSync(partyID string, readyToReceive chan bool, syncMessages chan *server.SyncMessageContent) error {
	received := make(chan interface{})
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopic(
			"autocar.party."+partyID+".sync."+arClient.playerUUID.String(),
			arClient.syncHandler,
			received,
			readyToReceive,
		)
		if err != nil {
			logger.Error("while receiving a sync message :", err)
			return
		}
	}()

	if !<-readyToReceive {
		return errors.New("something went wrong while listening to sync message")
	}
	for {
		response := <-received
		switch response.(type) {
		case error:
			logger.Error("error while decoding sync response :", response.(error))
		case *server.SyncMessageContent:
			syncMessages <- response.(*server.SyncMessageContent)
		default:
			logger.Error("sync message received but something wrong happened")
		}
	}
}

func (arClient *AutoraceClient) syncHandler(msg []byte) interface{} {
	syncResponse := new(server.SyncMessageContent)
	err := json.Unmarshal(msg, syncResponse)
	if err != nil {
		logger.Error("error while trying to unmarshal sync from server :", err)
		return err
	}
	return syncResponse
}

// SendSyncRequest send a sync request message to a dynamic server instance
func (arClient *AutoraceClient) SendSyncRequest() {
	playerToken := models.PlayerToken{
		ClientID: arClient.playerUUID.String(),
		PartyID:  arClient.partyUUID.String(),
	}
	arClient.rabbitConnection.SendMessageOnTopic(playerToken, "autocar.party."+arClient.partyUUID.String()+".sync")
}

// RequestPartyList send a request to a static server instance and receive a joinable party list by
// filtering received message and treat response as a slice of string
func (arClient *AutoraceClient) RequestPartyList(readyToReceive chan bool) ([]string, error) {
	received := make(chan interface{})
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopic("autocar.party.list."+arClient.playerUUID.String(), arClient.computePartyList, received, readyToReceive)
		if err != nil {
			logger.Error("error while trying to receive message on request party list :", err)
			readyToReceive <- false
			return
		}
	}()
	if !<-readyToReceive {
		return nil, errors.New("could not receive message on request party list")
	}
	go func() {
		arClient.rabbitConnection.SendMessageOnTopic(arClient.playerUUID.String(), "autocar.party.list")
	}()
	response := <-received
	switch response.(type) {
	case []string:
		return response.([]string), nil
	case string:
		return []string{response.(string)}, nil
	case error:
		return nil, response.(error)
	}
	return nil, nil
}

func (arClient *AutoraceClient) computePartyList(msg []byte) interface{} {
	var partyList []string
	var partyID string
	err := json.Unmarshal(msg, &partyList)
	if err != nil {
		err = json.Unmarshal(msg, &partyID)
		if err != nil {
			logger.Error("error while decoding server response :", string(msg))
			return err
		}
		return partyID
	}
	return partyList
}

// SendGameState is used to send a changing game's state request to a dynamic server instance
func (arClient *AutoraceClient) SendGameState(newState models.ChangeStateToken) {
	arClient.rabbitConnection.SendMessageOnTopic(newState, "autocar.party."+arClient.partyUUID.String()+".state")
}

func (arClient *AutoraceClient) computeNewGameState(msg []byte) interface{} {
	var newGameState models.ChangeStateAck
	err := json.Unmarshal(msg, &newGameState)
	if err != nil {
		return err
	}
	return newGameState
}

// ReceiveGameState receive a message from a dynamic server instance after a changing game state request
func (arClient *AutoraceClient) ReceiveGameState(partyID string, readyToReceive chan bool, gameState chan models.State) error {
	ready := make(chan bool)
	received := make(chan interface{})
	go func() {
		err := arClient.rabbitConnection.ReceiveMessageOnTopic(
			"autocar.party."+partyID+".state."+arClient.playerUUID.String(),
			arClient.computeNewGameState,
			received,
			ready,
		)
		if err != nil {
			logger.Error("error while trying to receive message on game state :", err)
			return
		}
	}()
	if !<-ready {
		readyToReceive <- false
		return errors.New("could not receive message on game state")
	}
	readyToReceive <- true
	for {
		response := <-received
		switch response.(type) {
		case models.ChangeStateAck:
			gameState <- response.(models.ChangeStateAck).NewState
		case error:
			return response.(error)
		}
	}
}

// Close terminate ongoing connection
func (arClient *AutoraceClient) Close() error {
	return arClient.rabbitConnection.Close()
}
