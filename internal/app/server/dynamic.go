package server

import (
	"encoding/json"
	"github.com/clnbs/autorace/internal/pkg/messaging/rabbit"
	"time"

	"github.com/streadway/amqp"

	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/internal/pkg/database"
	"github.com/clnbs/autorace/pkg/logger"
	"github.com/clnbs/autorace/pkg/systool"
)

// SyncMessageContent are send to every players every server's tick
type SyncMessageContent struct {
	PartyState  models.State `json:"party_state"`
	Competitors []*models.CompetitorActor
	MainActor   *models.MainActor
}

// DynamicPartyServer hold logic to run a party from the generation of the racetrack to the end of it.
// DynamicPartyServer also hold connection with clients.
type DynamicPartyServer struct {
	rabbitConnection           *rabbit.RabbitConnection
	redisConnection            *database.RedisClient
	party                      *models.Party
	closestRacetrackPointIndex map[string]int
	tickPerSecond              uint
}

// NewDynamicPartyServer create a DynamicPartyServer instance.
// NewDynamicPartyServer generate a party from a stored party configuration in a Redis database
// and send it back to the player who ask for its creation.
func NewDynamicPartyServer(partyID string, rabbitConfig rabbit.RabbitConnectionConfiguration) (*DynamicPartyServer, error) {
	dServer := new(DynamicPartyServer)
	dServer.tickPerSecond = 120
	dServer.closestRacetrackPointIndex = make(map[string]int)
	var err error
	dServer.rabbitConnection, err = rabbit.NewRabbitConnection(rabbitConfig)
	if err != nil {
		return nil, err
	}
	dServer.redisConnection = database.NewRedisClient()
	partyConfiguration, err := dServer.redisConnection.GetPartyCreationToken(partyID)
	if err != nil {
		return nil, err
	}
	dServer.party, err = models.NewParty(partyConfiguration, partyID)
	if err != nil {
		return nil, err
	}
	dServer.party.MapCircuit.MapGeneration(dServer.party.CircuitConfig)
	player, err := dServer.redisConnection.GetPlayer(partyConfiguration.ClientID)
	if err != nil {
		return nil, err
	}
	dServer.party.Players[player.PlayerUUID.String()] = player
	dServer.SendCreatedParty(player.PlayerUUID.String())
	return dServer, nil
}

// SendCreatedParty is used to send a just-created party back to the client who asked for it
func (dServer *DynamicPartyServer) SendCreatedParty(playerID string) {
	dServer.rabbitConnection.SendMessageOnTopic(dServer.party, "autocar.party.creation."+playerID)
}

// SendPartyToOnePlayer is used to send the party to a player who asked for it
func (dServer *DynamicPartyServer) SendPartyToOnePlayer(playerID string) {
	dServer.rabbitConnection.SendMessageOnTopic(dServer.party, "autocar.party."+dServer.party.PartyUUID.String()+".map."+playerID)
}

//ReceiveAddPlayer handle request to add a player in the ongoing party
func (dServer *DynamicPartyServer) ReceiveAddPlayer(readyToReceive chan bool) error {
	err := dServer.rabbitConnection.ReceiveMessageOnTopicWithHandler(
		"autocar.party."+dServer.party.PartyUUID.String()+".addPlayer", // topic
		dServer.addPlayerHandler, // handler
		readyToReceive,           // ready to receive chan
	)
	if err != nil {
		return err
	}
	return nil
}

// TODO send a error message if the Party is already started
func (dServer *DynamicPartyServer) addPlayerHandler(msg amqp.Delivery) {
	var addPlayerToken models.PlayerToken
	err := json.Unmarshal(msg.Body, &addPlayerToken)
	if err != nil {
		logger.Error("unable to unmarshal message from client :", err)
		return
	}
	newPlayer, err := dServer.redisConnection.GetPlayer(addPlayerToken.ClientID)
	if err != nil {
		logger.Error("while trying to register new player in party :", err)
		return
	}
	err = dServer.party.AddPlayer(newPlayer)
	if err != nil {
		logger.Error("while adding player in a party :", err)
	}
	dServer.closestRacetrackPointIndex[newPlayer.PlayerUUID.String()] = 0
	dServer.SendPartyToOnePlayer(newPlayer.PlayerUUID.String())
	dServer.SyncParty()
}

//SyncParty send a Sync Message to all players in the party
func (dServer *DynamicPartyServer) SyncParty() {
	for _, player := range dServer.party.Players {
		dServer.SyncPartyForOnePlayer(player.PlayerUUID.String())
	}
}

// SyncPartyForOnePlayer send a Sync Message to one player only
func (dServer *DynamicPartyServer) SyncPartyForOnePlayer(clientID string) {
	// creating the actual Sync Message
	syncMessage := &SyncMessageContent{
		PartyState: dServer.party.GetState(),
	}
	for playerID, player := range dServer.party.Players {
		if playerID == clientID {
			syncMessage.MainActor = &models.MainActor{
				Act: &models.Actor{
					Name: player.PlayerName,
					Rank: 0,
				},
				Player: player,
			}
			continue
		}
		syncMessage.Competitors = append(syncMessage.Competitors, &models.CompetitorActor{
			Act: &models.Actor{
				Name: player.PlayerName,
				Rank: 0,
			},
			Position:  player.Position,
			ActorUUID: player.PlayerUUID,
		})
	}
	//Sending it
	dServer.rabbitConnection.SendMessageOnTopic(
		syncMessage, // object to send
		"autocar.party."+dServer.party.PartyUUID.String()+".sync."+clientID, // topic
	)
}

// ReceiveSyncRequest handle sync request from one player
func (dServer *DynamicPartyServer) ReceiveSyncRequest(readyToReceive chan bool) error {
	received := make(chan interface{})
	err := dServer.rabbitConnection.ReceiveMessageOnTopicWithHeader(
		"autocar.party."+dServer.party.PartyUUID.String()+".sync", // topic
		dServer.syncRequestHandler,                                // handler function
		received,                                                  // received object chan
		readyToReceive,                                            // ready to receive chan
	)
	if err != nil {
		return err
	}
	for {
		msg := <-received
		switch msg.(type) {
		case models.PlayerToken:
			dServer.SyncPartyForOnePlayer(msg.(models.PlayerToken).ClientID)
		default:
			logger.Error("while receiving sync request :", msg.(error))
		}
	}
}

func (dServer *DynamicPartyServer) syncRequestHandler(msg amqp.Delivery) interface{} {
	msgBody := msg.Body
	var playerToken models.AddPlayerToken
	err := json.Unmarshal(msgBody, &playerToken)
	if err != nil {
		return err
	}
	return playerToken
}

// ReceiveNewState handle changing game state request
func (dServer *DynamicPartyServer) ReceiveNewState(readyToReceive chan bool) error {
	defer logger.Trace(systool.TimeTrack(time.Now(), "ReceiveNewState"))
	//received := make(chan interface{})
	err := dServer.rabbitConnection.ReceiveMessageOnTopicWithCallback(
		"autocar.party."+dServer.party.PartyUUID.String()+".state",
		"autocar.party."+dServer.party.PartyUUID.String()+".state",
		dServer.rabbitConnection.SendMessageOnTopic,
		dServer.newStateRequestHandler,
		readyToReceive,
	)
	if err != nil {
		logger.Error("while receiving new game state :", err)
		return err
	}
	return nil
}

func (dServer *DynamicPartyServer) newStateRequestHandler(msg amqp.Delivery) (interface{}, string) {
	msgBody := msg.Body
	var stateRequest models.ChangeStateToken
	err := json.Unmarshal(msgBody, &stateRequest)
	if err != nil {
		return err, ""
	}
	dServer.party.SetState(stateRequest.DesiredState)
	newState := models.ChangeStateAck{
		PartyID:      stateRequest.PlayerToken.PartyID,
		DesiredState: stateRequest.DesiredState,
		NewState:     dServer.party.GetState(),
		Message:      "OK",
	}
	return newState, "." + stateRequest.PlayerToken.ClientID
}

// ReceivePlayersInput receive and store locally players' inputs
func (dServer *DynamicPartyServer) ReceivePlayersInput(readyToReceive chan bool) error {
	received := make(chan interface{})
	go func() {
		err := dServer.rabbitConnection.ReceiveMessageOnTopic(
			"autocar.party."+dServer.party.PartyUUID.String()+".input",
			dServer.handlePlayerInput,
			received,
			readyToReceive,
		)
		if err != nil {
			logger.Error("error while listening to player input :", err)
			return
		}
	}()
	for {
		msg := <-received
		switch msg.(type) {
		case models.PlayerInput:
			newPlayerInput := new(models.PlayerInput)
			newPlayerInput.Acceleration = msg.(models.PlayerInput).Acceleration
			newPlayerInput.MessageNumber = msg.(models.PlayerInput).MessageNumber
			newPlayerInput.Timestamp = msg.(models.PlayerInput).Timestamp
			newPlayerInput.Turning = msg.(models.PlayerInput).Turning
			dServer.party.Players[msg.(models.PlayerInput).PlayerUUID.String()].Input = newPlayerInput
		default:
			logger.Error("got error while receiving player input :", msg.(error))
			return msg.(error)
		}
	}
}

func (dServer *DynamicPartyServer) handlePlayerInput(msg []byte) interface{} {
	var playerInput models.PlayerInput
	err := json.Unmarshal(msg, &playerInput)
	if err != nil {
		return err
	}
	return playerInput
}

// Close terminate connection with Redis and RabbitMQ
func (dServer *DynamicPartyServer) Close() error {
	return dServer.rabbitConnection.Close()
}

// Run start the actual game loop until the party is over
func (dServer *DynamicPartyServer) Run() {
	tickDuration := time.Duration((1.0/float64(dServer.tickPerSecond))*1000000000) * time.Nanosecond
	ticker := time.NewTicker(tickDuration)
	last := time.Now()
	tick := 0
	go func() {
		second := time.NewTicker(time.Second)
		for {
			<-second.C
			if uint(tick) < dServer.tickPerSecond {
				logger.Warning("dynamic server's tick is too low :", tick)
			}
			tick = 0
		}
	}()
	for {
		<-ticker.C
		deltaTime := time.Since(last).Seconds()
		last = time.Now()
		switch dServer.party.GetState() {
		case models.LOBBY:
			dServer.setCarAtStart()
			dServer.SyncParty()
		case models.END:
			//TODO end game and self destruct and remove container as well
			ticker.Stop()
			dServer.SyncParty()
			time.Sleep(1 * time.Second)
			return
		case models.PAUSE:
			dServer.SyncParty()
		case models.RUN:
			dServer.computeNewPosition(deltaTime)
			dServer.SyncParty()
		}
		tick++
	}
}
