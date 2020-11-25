package internal

import (
	"../framework"
)

func NextSong(ctx framework.Context) {
	ctx.MediaPlayer.Skip()
}