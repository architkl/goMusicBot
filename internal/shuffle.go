package internal

import (
	"../framework"
)

func Shuffle(ctx framework.Context) {
	ctx.MediaPlayer.Shuffle()
	ctx.ReplyEmbed("Queue Shuffled!", "", 0x3365A0)
}