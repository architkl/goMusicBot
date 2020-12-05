package internal

import (
	"../framework"
)

func ClearQueue(ctx framework.Context) {
	ctx.MediaPlayer.ClearQueue()
	ctx.ReplyEmbed("Cleared!", "", 0x00C49A)
}