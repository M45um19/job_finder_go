package idgen

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init initializes the Snowflake Node with a unique node ID (0-1023)
func Init(nodeID int64) {
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed to initialize snowflake ID generator: %v", err)
	}
}

// NextID generates a unique 64-bit Snowflake ID
func NextID() int64 {
	if node == nil {
		// Default to node 1 if not initialized explicitly
		Init(1)
	}
	return node.Generate().Int64()
}
