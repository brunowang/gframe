package gfid

import (
	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
	"math/rand"
)

func GenUUID() string {
	v4, err := uuid.NewV4()
	if err != nil {
		// unexpect error
		panic(err)
	}
	return v4.String()
}

func GenID() string {
	node, err := snowflake.NewNode(rand.Int63n(1024))
	if err != nil {
		// unexpect error
		panic(err)
	}
	return node.Generate().Base58()
}
