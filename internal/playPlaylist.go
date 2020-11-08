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
	args := strings.Fields(ctx.Message.Content)

	if len(args) < 2 {
		ctx.Reply("Enter the playlist name!")
		return
	}

	playlistName := args[1]

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
		// Find the channel that the message came from.
		textChannel, err := ctx.Discord.State.Channel(ctx.Message.ChannelID)
		if err != nil {
			// Could not find channel.
			log.Println("Could not find channel: ", err)
			return
		}

		ctx.MediaPlayer.StartPlaying(ctx.Discord, textChannel.GuildID, ctx.Message.Author.ID)
	}

	// videoId, title := Search(query)
}