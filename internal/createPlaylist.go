package internal

import (
	"../framework"
	"log"
	"os"
)

func CreatePlaylist(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	}

	playlistName := args[0]

	// Create playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_CREATE, 0777)
	if err != nil {
		log.Println("CreatePlaylist(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Couldn't create playlist", 0xEB5160)
		return
	}
	file.Close()

	ctx.ReplyEmbed(playlistName+" created successfully!", "", 0x00C49A)
}