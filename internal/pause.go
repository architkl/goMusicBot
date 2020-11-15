package internal

import (
	"../framework"
)

func Pause(ctx framework.Context) {
	ctx.MediaPlayer.Pause()
}