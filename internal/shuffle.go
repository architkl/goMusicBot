package internal

import (
	"architkl/goMusicBot/framework"
)

func Shuffle(ctx framework.Context) {
	ctx.MediaPlayer.Shuffle()
	ctx.ReplyEmbed("Queue Shuffled!", "", 0x3365A0)
}