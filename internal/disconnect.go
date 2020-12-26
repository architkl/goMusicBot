package internal

import (
	"architkl/goMusicBot/framework"
)

func Disconnect(ctx framework.Context) {
	ctx.MediaPlayer.DisconnectVoice()
}