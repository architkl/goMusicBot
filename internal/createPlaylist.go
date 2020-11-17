package internal

import (
	"../framework"
	"log"
	"os"
)

func CreatePlaylist(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.Reply("Enter the playlist name!")
		return
	}

	playlistName := args[0]

	// Create playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_CREATE, 0777)
	if err != nil {
		log.Println("CreatePlaylist(): " + err.Error())
		ctx.Reply("Couldn't create playlist")
		return
	}
	file.Close()

	ctx.Reply(playlistName + " created successfully!")
}