package database

import (
	"context"
	"encoding/json"

	"github.com/clnbs/autorace/internal/app/models"

	"github.com/go-redis/redis/v8"
)

// RedisClient hold connection to three separate database :
// - one to store party configuration related object
// - one to store player configuration related object
// - one to store running party related object
type RedisClient struct {
	redisPartyConfigConnection  *redis.Client
	redisPlayerConfigConnection *redis.Client
	redisRunningConnection      *redis.Client
}

//NewRedisClient create a RedisClient and connect each endpoint
//TODO get configuration from config file
//TODO test connectivity with a ping first
func NewRedisClient() *RedisClient {
	redisClient := new(RedisClient)
	redisClient.redisPartyConfigConnection = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	redisClient.redisPlayerConfigConnection = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       1,
	})
	redisClient.redisRunningConnection = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       2,
	})
	return redisClient
}

// SetPartyConfiguration store a party configuration a Redis database
func (rdsClient *RedisClient) SetPartyConfiguration(partyID string, partyOption models.PartyCreationToken) error {
	ctx := context.Background()
	stringifyPartyOption, err := json.Marshal(partyOption)
	if err != nil {
		return err
	}
	status := rdsClient.redisPartyConfigConnection.Set(ctx, partyID, string(stringifyPartyOption), 0)
	return status.Err()
}

// GetPartyList returns a list of all keys store in the Party Config database. Keys are Party UUID
func (rdsClient *RedisClient) GetPartyList() ([]string, error) {
	ctx := context.Background()
	result, err := rdsClient.redisPartyConfigConnection.Do(ctx, "KEYS", "*").Result()
	if err != nil {
		return nil, err
	}
	partyList := make([]string, len(result.([]interface{})))
	for index, r := range result.([]interface{}) {
		partyList[index] = r.(string)
	}
	return partyList, nil
}

// GetPartyCreationToken returns a party config store in a redis database
func (rdsClient *RedisClient) GetPartyCreationToken(partyID string) (models.PartyCreationToken, error) {
	ctx := context.Background()
	stringifyPartyToken, err := rdsClient.redisPartyConfigConnection.Get(ctx, partyID).Result()
	if err != nil {
		return models.PartyCreationToken{}, err
	}
	var partyToken models.PartyCreationToken
	err = json.Unmarshal([]byte(stringifyPartyToken), &partyToken)
	if err != nil {
		return models.PartyCreationToken{}, err
	}
	return partyToken, nil
}

// RemovePartyCreationToken remove a party configuration in a Redis database
func (rdsClient *RedisClient) RemovePartyCreationToken(partyID string) error {
	ctx := context.Background()
	_, err := rdsClient.redisPartyConfigConnection.Del(ctx, partyID).Result()
	return err
}

// SetPlayer store a player in a Redis database bind on its UUID
func (rdsClient *RedisClient) SetPlayer(player *models.Player) error {
	ctx := context.Background()
	stringifyPlayer, err := json.Marshal(player)
	if err != nil {
		return err
	}
	status := rdsClient.redisPlayerConfigConnection.Set(ctx, player.PlayerUUID.String(), string(stringifyPlayer), 0)
	return status.Err()
}

// GetPlayer get a player from player's UUID
func (rdsClient *RedisClient) GetPlayer(playerID string) (*models.Player, error) {
	ctx := context.Background()
	stringifyPlayer, err := rdsClient.redisPlayerConfigConnection.Get(ctx, playerID).Result()
	if err != nil {
		return nil, err
	}
	player := new(models.Player)
	err = json.Unmarshal([]byte(stringifyPlayer), player)
	return player, err
}

// SetPlayerOnParty store a player UUID bind on a Party UUID
func (rdsClient *RedisClient) SetPlayerOnParty(partyID, playerID string) error {
	ctx := context.Background()
	var alreadyRegisteredPlayer []string
	exists, err := rdsClient.redisRunningConnection.Exists(ctx, partyID).Result()
	if err != nil {
		return err
	}
	if exists != 0 {
		alreadyRegisteredPlayer, err = rdsClient.GetPlayersOnParty(partyID)
		if err != nil {
			return err
		}
	}
	alreadyRegisteredPlayer = append(alreadyRegisteredPlayer, playerID)
	stringifyPlayer, err := json.Marshal(alreadyRegisteredPlayer)
	if err != nil {
		return err
	}
	return rdsClient.redisRunningConnection.Set(ctx, partyID, string(stringifyPlayer), 0).Err()
}

// GetPlayersOnParty returns players UUID bind by a party UUID
func (rdsClient *RedisClient) GetPlayersOnParty(partyID string) ([]string, error) {
	ctx := context.Background()
	stringifyPlayers, err := rdsClient.redisRunningConnection.Get(ctx, partyID).Result()
	if err != nil {
		return nil, err
	}
	var player []string
	err = json.Unmarshal([]byte(stringifyPlayers), &player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// Close terminate Redis connection
func (rdsClient *RedisClient) Close() error {
	err := rdsClient.redisRunningConnection.Close()
	if err != nil {
		return err
	}
	err = rdsClient.redisPartyConfigConnection.Close()
	if err != nil {
		return err
	}
	return rdsClient.redisPlayerConfigConnection.Close()
}
