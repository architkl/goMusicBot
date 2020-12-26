package internal

import "architkl/goMusicBot/framework"

func Ping(ctx framework.Context) {
	ctx.ReplyEmbed("pong", "", 0x5539CC)
}