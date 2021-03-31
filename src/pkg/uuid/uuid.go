package uuid

import (
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/oklog/ulid/v2"
	"math/rand"
)

var Generate = generate

func generate() string {
	now := clock.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(now), entropy).String()
}
