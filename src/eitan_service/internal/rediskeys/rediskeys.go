package rediskeys

import "github.com/k-yomo/rediskey"

var toplavelNamespace = rediskey.NewNamespace("eitan", nil)
var WaitingRoomNamespace = rediskey.NewNamespace("waiting-room", toplavelNamespace)

func NewWaitingRoomPlayerKey(playerID string) *rediskey.Key {
	return rediskey.NewKey("players", playerID, WaitingRoomNamespace)
}