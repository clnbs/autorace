package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/clnbs/autorace/internal/app/engine"
	"github.com/clnbs/autorace/internal/app/models"
	"github.com/clnbs/autorace/pkg/logger"

	"github.com/faiface/pixel/pixelgl"
)

func init() {
	logger.SetStdLogger("trace", "stdout")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	config := engine.NewWindowConfiguration()
	fmt.Print("Enter your player name : ")
	playerName, _ := reader.ReadString('\n')
	playerName = strings.Replace(playerName, "\n", "", -1)
	playerName = strings.Replace(playerName, "\r", "", -1)
	fmt.Print("Enter server address : ")
	rabbitMQAddr, _ := reader.ReadString('\n')
	rabbitMQAddr = strings.Replace(rabbitMQAddr, "\n", "", -1)
	rabbitMQAddr = strings.Replace(rabbitMQAddr, "\r", "", -1)
	var rabbitMQPort int
	fmt.Print("Enter server port : ")
	_, err := fmt.Scanf("%d", &rabbitMQPort)
	if rabbitMQPort == 0 {
		rabbitMQPort = 5672
	}
	mainWindow, err := engine.NewMainGameWindow(config, playerName, rabbitMQAddr, rabbitMQPort)
	if err != nil {
		logger.Error("error while setting up mainWindow :", err)
		return
	}
	fmt.Print("Do you want to create a game ? (y/n) : ")
	var createOrNot string
	for createOrNot != "y" && createOrNot != "yes" && createOrNot != "n" && createOrNot != "no" {
		createOrNot, _ = reader.ReadString('\n')
		createOrNot = strings.Replace(createOrNot, "\n", "", -1)
		createOrNot = strings.Replace(createOrNot, "\r", "", -1)
		createOrNot = strings.ToLower(createOrNot)
	}
	if createOrNot == "y" || createOrNot == "yes" {
		err = createParty(mainWindow)
		if err != nil {
			panic(err)
		}
		return
	}
	err = joinParty(mainWindow)
	if err != nil {
		panic(err)
	}
}

func joinParty(mainWindow *engine.MainGameWindow) error {
	readyToReceive := make(chan bool)
	partyList, err := mainWindow.GameInfo.GetPartyList()
	if err != nil {
		logger.Error("error while getting party list :", err)
		return err
	}
	fmt.Println("Choose a party from the list below : ")
	for index, party := range partyList {
		fmt.Println("party", index+1, ":", party)
	}

	chosenParty := -1

	for chosenParty < 0 || chosenParty > len(partyList) {
		_, err = fmt.Scanf("%d", &chosenParty)
		if err != nil {
			return err
		}
		chosenParty--
	}
	fmt.Println("chosen party", chosenParty, ":", partyList[chosenParty])

	go func() {
		logger.Debug("about to start communication daemon")
		err = mainWindow.StartCommunicationDaemonWithPartyID(partyList[chosenParty], readyToReceive)
		if err != nil {
			logger.Error("error while starting communication daemon :", err)
			panic(err)
		}
	}()
	logger.Debug("waiting for readyToReceive ..")
	if !<-readyToReceive {
		logger.Error("unable to start communication daemon")
		return errors.New("unable to start communication daemon")
	}
	logger.Trace("about to add player in party")
	err = mainWindow.GameInfo.AddPlayerToAParty(partyList[chosenParty])
	if err != nil {
		logger.Error("could not add player in party :", err)
		return err
	}
	pixelgl.Run(mainWindow.Run)
	err = mainWindow.GameInfo.Close()
	if err != nil {
		logger.Error("while closing ongoing connection :", err)
		return err
	}
	return nil
}

func createParty(mainWindow *engine.MainGameWindow) error {
	partyToken := models.PartyCreationToken{
		ClientID:  mainWindow.GameInfo.ActorPlayer.Player.PlayerUUID.String(),
		Seed:      0,
		PartyName: "totos party",
		CircuitConfig: models.CircuitMapConfig{
			Seed:     0,
			MaxPoint: 100,
			MinPoint: 50,
			XSize:    4000,
			YSize:    4000,
		},
	}
	err := mainWindow.GameInfo.GetNewParty(partyToken)
	if err != nil {
		logger.Error("error while requesting track :", err)
		return err
	}

	err = mainWindow.StartCommunicationDaemon()
	if err != nil {
		logger.Error("error while starting communication daemon :", err)
		panic(err)
	}
	pixelgl.Run(mainWindow.Run)
	err = mainWindow.GameInfo.Close()
	if err != nil {
		logger.Error("while closing ongoing connection :", err)
		return err
	}
	return nil
}
