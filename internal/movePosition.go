package internal

import (
	"architkl/goMusicBot/framework"
	"log"
	"strconv"
)

func MovePosition(ctx framework.Context) {

	// Move song in queue
	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter original postion!", 0xEB5160)
		return
	} else if len(args) < 2 {
		ctx.ReplyEmbed("Oops!", "Enter new position!", 0xEB5160)
		return
	}

	src, convErr := strconv.Atoi(args[0])
	if convErr != nil {
		log.Println("MovePosition(): ", convErr.Error())
		ctx.ReplyEmbed("Oops!", "Invalid index!", 0xEB5160)
	}
	src--

	dest, convErr := strconv.Atoi(args[1])
	if convErr != nil {
		log.Println("MovePosition(): ", convErr.Error())
		ctx.ReplyEmbed("Oops!", "Invalid index!", 0xEB5160)
	}
	dest--

	if curPosition, err := ctx.MediaPlayer.CurPosition(); err != nil {
		log.Println("MovePosition(): ", err.Error())
		ctx.ReplyEmbed("Oops!", "Error in getting current position", 0xEB5160)
		return
	} else if src == curPosition || dest == curPosition {
		ctx.ReplyEmbed("Oops!", "Can't move the song currently playing", 0xEB5160)
		return
	}

	if err := ctx.MediaPlayer.MoveSongPosition(src, dest); err != nil {
		log.Println("MovePosition(): ", err.Error())
		ctx.ReplyEmbed("Oops!", err.Error(), 0xEB5160)
	}

	ctx.ReplyEmbed(ctx.MediaPlayer.Queue[dest].Title+" moved from "+args[0]+" to "+args[1], "", 0x3365A0)
}