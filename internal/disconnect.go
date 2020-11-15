package internal

import (
	"../framework"
)

func Disconnect(ctx framework.Context) {
	ctx.MediaPlayer.DisconnectVoice()
}

func Connect(ctx framework.Context) {
	ctx.MediaPlayer.ConnectVoice(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
}