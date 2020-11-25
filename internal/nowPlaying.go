package internal

import (
	"../framework"
)

func NowPlaying(ctx framework.Context) {

	if pos, err := ctx.MediaPlayer.CurPosition(); err != nil || pos == len(ctx.MediaPlayer.Queue) {
		ctx.ReplyEmbed("Oops!", "No song currently playing", 0xEB5160)
	} else {
		ctx.ReplyEmbed("Now Playing", "```ml\n"+ctx.MediaPlayer.Queue[pos].Title+" "+ctx.MediaPlayer.Queue[pos].Duration+"```", 0x228FD3)
	}
}