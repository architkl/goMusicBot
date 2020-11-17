package internal

import (
	"../framework"
)

func Disconnect(ctx framework.Context) {
	ctx.MediaPlayer.DisconnectVoice()
}