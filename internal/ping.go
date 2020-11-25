package internal

import "../framework"

func Ping(ctx framework.Context) {
	ctx.ReplyEmbed("pong", "", 0x5539CC)
}