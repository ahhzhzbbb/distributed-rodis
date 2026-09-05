package command

import (
	"rodis/pkg/resp"
	"rodis/storage/internal/engine"
)

type CommandContext struct {
	k *engine.KeyValue
}

func NewCommandContext(k *engine.KeyValue) *CommandContext {
	return &CommandContext{
		k: k,
	}
}

type Command interface {
	Execute(args []resp.Payload, ctx *CommandContext) resp.Payload
}
