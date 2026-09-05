package command

import "rodis/pkg/resp"

type ConfigCommand struct{}

func (c *ConfigCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	temp := make([]resp.Payload, 0)
	return resp.NewArray(temp)
}
