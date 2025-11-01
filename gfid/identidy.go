package gfid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"math/rand"
)

func GenUUID() string {
	v4, err := uuid.NewV7()
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
