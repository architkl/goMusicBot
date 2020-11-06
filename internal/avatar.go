package internal

import "../framework"

func Avatar(ctx framework.Context) {
	ctx.Reply(ctx.Message.Author.AvatarURL("2048"))
}