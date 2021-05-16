package uuid

import (
	"github.com/lithammer/shortuuid/v3"
)

var Generate = generate

func generate() string {
	return shortuuid.New()
}
