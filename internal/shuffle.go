package internal

import (
	"../framework"
)

func Shuffle(ctx framework.Context) {
	ctx.MediaPlayer.Shuffle()
}