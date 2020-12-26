package internal

import (
	"architkl/goMusicBot/framework"
	"log"
	"os"
)

func RemovePlaylist(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	}

	playlistName := args[0]

	// Delete playlist
	if err := os.Remove("./docs/playlists/" + playlistName + ".txt"); err != nil {
		log.Println("RemovePlaylist(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Couldn't remove playlist", 0xEB5160)
		return
	}

	ctx.ReplyEmbed(playlistName+" removed successfully!", "", 0x00C49A)
}