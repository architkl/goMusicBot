package internal

import (
	"architkl/goMusicBot/framework"
)

func ClearQueue(ctx framework.Context) {
	ctx.MediaPlayer.ClearQueue()
	ctx.ReplyEmbed("Cleared!", "", 0x00C49A)
}