package internal

import (
	"../framework"
	"../pkg"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func PlayPlaylist(ctx framework.Context) {
	args := ctx.Args

	if len(args) < 1 {
		ctx.Reply("Enter the playlist name!")
		return
	}

	playlistName := args[0]

	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_RDONLY, 0755)
	if err != nil {
		log.Println(err)
		ctx.Reply("Playlist not found")
		return
	}

	scanner := bufio.NewScanner(file)
	var songs []framework.Song
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ",")

		fmt.Println(words[0] + " " + words[1])

		songs = append(songs, framework.Song{
			Id:    words[0],
			Title: words[1],
		})
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
		return
	}

	err = file.Close()
	pkg.HandleError(err, "")

	ctx.MediaPlayer.AddSongs(songs...)

	// start the player if its not running
	if !ctx.MediaPlayer.IsRunning {
		ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
	}
}