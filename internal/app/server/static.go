package server

import (
	"encoding/json"
	"errors"
	"github.com/clnbs/autorace/internal/app/models"
	"os"
	"time"

	"github.com/clnbs/autorace/internal/pkg/container"
	"github.com/clnbs/autorace/internal/pkg/database"
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/clnbs/autorace/pkg/logger"
	"github.com/clnbs/autorace/pkg/systool"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

//StaticServer handle creation request and the list of ongoing parties who had not launched yet.
// Every time a party is created, it start a new DynamicPartyServer instance
type StaticServer struct {
	rabbitConnection *messaging.RabbitConnection
	redisConnection  *database.RedisClient
}

// NewStaticServer generate a StaticServer (aka static server) with a particular RabbitMQ address
func NewStaticServer(rabbitConfig messaging.RabbitConnectionConfiguration) (*StaticServer, error) {
	newCreatorServer := new(StaticServer)
	var err error
	newCreatorServer.redisConnection = database.NewRedisClient()
	newCreatorServer.rabbitConnection, err = messaging.NewRabbitConnection(rabbitConfig)
	return newCreatorServer, err
}

//ReceivePlayerCreation create a player instance and store it in a Redis database
// Player creation use case
func (staticServer *StaticServer) ReceivePlayerCreation(readyToReceive chan bool) error {
	return staticServer.rabbitConnection.ReceiveMessageOnTopicWithCallback(
		"autocar.player.creation",                        //topic
		"autocar.player.creation",                        //response topic
		staticServer.rabbitConnection.SendMessageOnTopic, //callback
		staticServer.playerCreator,                       // response creator
		readyToReceive,                                   // ready to receive chan
	)
}

func (staticServer *StaticServer) playerCreator(msg amqp.Delivery) (interface{}, string) {
	defer logger.Trace(systool.TimeTrack(time.Now(), "playerCreator"))
	var playerCreationToken models.PlayerCreationToken
	err := json.Unmarshal(msg.Body, &playerCreationToken)
	if err != nil {
		return err, "." + playerCreationToken.SessionUUID.String()
	}
	newPlayer := models.NewPlayer(playerCreationToken.PlayerName)
	err = staticServer.redisConnection.SetPlayer(newPlayer)
	if err != nil {
		return err, "." + playerCreationToken.SessionUUID.String()
	}
	return newPlayer, "." + playerCreationToken.SessionUUID.String()
}

//ReceivePartyCreation store a party configuration in a Redis data and start a
// DynamicServer with a generated UUID
// Party creation use case
func (staticServer *StaticServer) ReceivePartyCreation(readyToReceive chan bool) error {
	err := staticServer.rabbitConnection.ReceiveMessageOnTopicWithHandler(
		"autocar.party.creation",  //topic
		staticServer.partyCreator, // response creator
		readyToReceive,            // ready to receive chan
	)
	return err
}

func (staticServer *StaticServer) partyCreator(msg amqp.Delivery) {
	defer logger.Trace(systool.TimeTrack(time.Now(), "partyCreator"))
	var partyCreationToken models.PartyCreationToken
	err := json.Unmarshal(msg.Body, &partyCreationToken)
	if err != nil {
		logger.Error("error while decoding party creation token :", err)
	}
	newPartyUUID := uuid.New()
	err = staticServer.redisConnection.SetPartyConfiguration(newPartyUUID.String(), partyCreationToken)
	if err != nil {
		staticServer.rabbitConnection.SendMessageOnTopic(errors.New("unable to register party with current configuration"), "autocar.parties.creation."+partyCreationToken.ClientID)
		logger.Error("unable to register party :", err)
	}
	envConfig := []string{
		"FLUENTD_HOST=" + os.Getenv("FLUENTD_HOST"),
		"FLUENTD_PORT=" + os.Getenv("FLUENTD_PORT"),
		"LOG_LEVEL=" + os.Getenv("LOG_LEVEL"),
		"RABBITMQ_HOST=" + os.Getenv("RABBITMQ_HOST"),
		"RABBITMQ_PORT=" + os.Getenv("RABBITMQ_PORT"),
		"RABBITMQ_USER=" + os.Getenv("RABBITMQ_USER"),
		"RABBITMQ_PASS=" + os.Getenv("RABBITMQ_PASS"),
	}
	err = container.CreateDynamicServer(newPartyUUID.String(), envConfig)
	if err != nil {
		staticServer.rabbitConnection.SendMessageOnTopic(errors.New("unable to start a party server"), "autocar.parties.creation."+partyCreationToken.ClientID)
		logger.Error("unable to start a party container :", err)
	}
}

//PartyListRequest generate a party list from registered party in a Redis database
// List all current parties use case
func (staticServer *StaticServer) PartyListRequest(readyToReceive chan bool) error {
	defer logger.Trace(systool.TimeTrack(time.Now(), "PartyListRequest"))
	err := staticServer.rabbitConnection.ReceiveMessageOnTopicWithCallback(
		"autocar.party.list",                             //topic
		"autocar.party.list",                             //response topic
		staticServer.rabbitConnection.SendMessageOnTopic, //callback
		staticServer.partyListCreator,                    // response creator
		readyToReceive,                                   // ready to receive chan
	)
	return err
}

func (staticServer *StaticServer) partyListCreator(msg amqp.Delivery) (interface{}, string) {
	var partyList []string
	var clientID string
	err := json.Unmarshal(msg.Body, &clientID)
	if err != nil {
		logger.Error("could not unmarshal client ID :", err)
		return err, ""
	}
	partyList, err = staticServer.redisConnection.GetPartyList()
	if err != nil {
		return err, "." + clientID
	}
	return partyList, "." + clientID
}

// Close is used to terminate ongoing connection with Redis and RabbitMQ
func (staticServer *StaticServer) Close() error {
	err := staticServer.redisConnection.Close()
	if err != nil {
		return err
	}
	return staticServer.rabbitConnection.Close()
}

// END list all current parties use case
