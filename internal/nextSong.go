package internal

import (
	"architkl/goMusicBot/framework"
)

func NextSong(ctx framework.Context) {
	ctx.MediaPlayer.Skip()
}