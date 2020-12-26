package internal

import "architkl/goMusicBot/framework"

func Username(ctx framework.Context) {
	ctx.ReplyEmbed("Username", ctx.Message.Author.Username, 0x00C49A)
}