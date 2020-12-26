package internal

import (
	"architkl/goMusicBot/framework"
)

func Pause(ctx framework.Context) {
	ctx.MediaPlayer.Pause()
}