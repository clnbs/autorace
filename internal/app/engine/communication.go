package engine

import (
	"github.com/clnbs/autorace/internal/app/client"
	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/app/server"
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/clnbs/autorace/pkg/logger"
	"strconv"

	"github.com/faiface/pixel"
)

// GameCommunication is a client.AutoraceClient wrapper who handle communication
// between game interface and servers. It is use to feed actors and party content.
// It can dialogue with main game interface via an event channel if needed.
type GameCommunication struct {
	Client      *client.AutoraceClient
	Party       *models.Party
	ActorPlayer *models.MainActor
	Competitors map[string]*models.CompetitorActor
	CheckPoints []*models.Checkpoint // unused for now
	events      chan models.Event
}

// NewGameCommunication create game communication handler by feeding some of the main
// GameCommunication content.
func NewGameCommunication(name, rabbitAddr string, rabbitPort int, events chan models.Event) (*GameCommunication, error) {
	newClient := new(GameCommunication)
	var err error
	rabbitMQConfig := messaging.RabbitConnectionConfiguration{
		Host:     rabbitAddr,
		Port:     strconv.FormatInt(int64(rabbitPort), 10),
		User:     "guest",
		Password: "guest",
	}
	newClient.Client, err = client.NewAutoraceClient(name, rabbitMQConfig)
	if err != nil {
		return nil, err
	}
	newClient.events = events
	newClient.Party = new(models.Party)
	newClient.ActorPlayer = new(models.MainActor)
	newClient.ActorPlayer.Act = new(models.Actor)
	newClient.Competitors = make(map[string]*models.CompetitorActor)
	newClient.CheckPoints = make([]*models.Checkpoint, 0)
	return newClient, nil
}

// GetNewPlayer handle player registration to servers via a RabbitMQ connection.
// Under the hood, a static server instance create an player object, register it in a
// Redis database and send it back to client.
func (gameCommunication *GameCommunication) GetNewPlayer() error {
	readyToReceive := make(chan bool)
	var err error
	gameCommunication.ActorPlayer.Player, err = gameCommunication.Client.RequestPlayerCreation(readyToReceive)
	return err
}

// GetNewParty handle party creation by registering it and receiving it back.
// Under the hood, a static server instance register party's configuration in a redis
// database and start a dynamic server instance.
// The dynamic server instance get the party's configuration, create it and send it
// back to the client.
func (gameCommunication *GameCommunication) GetNewParty(partyConfiguration models.PartyCreationToken) error {
	readyToReceive := make(chan bool)
	var err error
	gameCommunication.Party, err = gameCommunication.Client.RequestPartyCreation(partyConfiguration, readyToReceive)
	return err
}

// ReceiveParty receive a party instance after asking for it. It is used when the asked party is already
// created and a dynamic server instance is already created
func (gameCommunication *GameCommunication) ReceiveParty(partyID string, readyToReceive chan bool) error {
	var err error
	gameCommunication.Party, err = gameCommunication.Client.ReceiveParty(partyID, readyToReceive)
	if err != nil {
		return err
	}
	return nil
}

// HandleSync receive sync message from a dynamic server instance. Sync message are
// use to update players (main actor and competitors) position in game
func (gameCommunication *GameCommunication) HandleSync(partyID string, readyToReceive chan bool) error {
	syncMessages := make(chan *server.SyncMessageContent)
	// start ReceiveSync from client and communicate sync message received to
	// this interface via a chan
	go func() {
		err := gameCommunication.Client.ReceiveSync(partyID, readyToReceive, syncMessages)
		if err != nil {
			logger.Error("something went wrong while listening to sync message :", err)
			return
		}
	}()
	readyToReceive <- true
	for {
		syncMessage := <-syncMessages
		reloadCar := false
		if len(gameCommunication.Competitors) != len(syncMessage.Competitors) {
			reloadCar = true
		}
		gameCommunication.assignSyncMessageToActors(syncMessage)
		gameCommunication.Party.SetState(syncMessage.PartyState)
		if reloadCar {
			gameCommunication.events <- models.AddCar{}
		}
	}
}

func (gameCommunication *GameCommunication) assignSyncMessageToActors(syncMessage *server.SyncMessageContent) {
	//Main actor
	gameCommunication.ActorPlayer.Act.Rank = syncMessage.MainActor.Act.Rank
	gameCommunication.ActorPlayer.Act.Car.Position = pixel.Vec{
		X: syncMessage.MainActor.Player.Position.CurrentPosition.X,
		Y: syncMessage.MainActor.Player.Position.CurrentPosition.Y,
	}
	gameCommunication.ActorPlayer.Act.Car.Angle = syncMessage.MainActor.Player.Position.CurrentAngle
	for _, c := range syncMessage.Competitors {
		if _, ok := gameCommunication.Competitors[c.ActorUUID.String()]; !ok {
			gameCommunication.Competitors[c.ActorUUID.String()] = &models.CompetitorActor{
				Act: &models.Actor{
					Car:  new(models.Car),
					Name: "",
					Rank: 0,
				},
				ActorUUID: c.ActorUUID,
				Position:  new(models.PlayerPosition),
			}
		}
		gameCommunication.Competitors[c.ActorUUID.String()].Act.Car.Position = pixel.Vec{
			X: c.Position.CurrentPosition.X,
			Y: c.Position.CurrentPosition.Y,
		}
		gameCommunication.Competitors[c.ActorUUID.String()].Act.Car.Angle = c.Position.CurrentAngle
		gameCommunication.Competitors[c.ActorUUID.String()].Act.Rank = c.Act.Rank
	}
}

// SendPlayerInput send player input to dynamic server instance
func (gameCommunication *GameCommunication) SendPlayerInput() {
	gameCommunication.Client.SendPlayerInput(gameCommunication.ActorPlayer.Player.Input)
}

// GetPartyList request joinable party list from static sever instance
func (gameCommunication *GameCommunication) GetPartyList() ([]string, error) {
	readyToReceive := make(chan bool)
	return gameCommunication.Client.RequestPartyList(readyToReceive)
}

// AddPlayerToAParty send request to add the player to a party. The party can only be join
// if the party is not started
func (gameCommunication *GameCommunication) AddPlayerToAParty(partyID string) error {
	return gameCommunication.Client.AddPlayerRequest(partyID)
}

// Sync request a sync message from server. /!\ it can only be trigger if HandleSync
// is already started and able to handle messages
func (gameCommunication *GameCommunication) Sync() {
	gameCommunication.Client.SendSyncRequest()
}

// SendState send a game state request to a dynamic server instance. /!\ it can only
// be trigger if HandleGameState is started and able to handle messages
func (gameCommunication *GameCommunication) SendState(newState models.State) error {
	stateToken := models.ChangeStateToken{
		PlayerToken: models.PlayerToken{
			ClientID: gameCommunication.ActorPlayer.Player.PlayerUUID.String(),
			PartyID:  gameCommunication.Party.PartyUUID.String(),
		},
		DesiredState: newState,
	}
	gameCommunication.Client.SendGameState(stateToken)
	return nil
}

// HandleGameState handle changing state message from server
func (gameCommunication *GameCommunication) HandleGameState(partyID string, readyToReceive chan bool) error {
	gameState := make(chan models.State)
	go func() {
		err := gameCommunication.Client.ReceiveGameState(partyID, readyToReceive, gameState)
		if err != nil {
			logger.Error("while listening to game state :", err)
			return
		}
	}()
	readyToReceive <- true
	for {
		newState := <-gameState
		gameCommunication.Party.SetState(newState)
	}
}

//Close terminate ongoing connection
func (gameCommunication *GameCommunication) Close() error {
	return gameCommunication.Client.Close()
}
