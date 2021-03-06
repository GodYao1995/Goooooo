package tools

import (
	"github.com/bwmarrin/snowflake"
)

const (
	_Node = 1
)

func SnowId() int64 {
	node, _ := snowflake.NewNode(_Node)
	return node.Generate().Int64()
}
