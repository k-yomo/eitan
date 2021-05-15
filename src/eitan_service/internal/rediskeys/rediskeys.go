package rediskeys

import "github.com/k-yomo/rediskey"

var toplavelNamespace = rediskey.NewNamespace("eitan", nil)
var WaitingRoomNamespace = rediskey.NewNamespace("waiting-room", toplavelNamespace)

func NewWaitingRoomUserKey(userID string) *rediskey.Key {
	return rediskey.NewKey("users", userID, WaitingRoomNamespace)
}