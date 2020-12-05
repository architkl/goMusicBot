package internal

import (
	"../framework"
	"bufio"
	"log"
	"os"
)

// Play the given playlist
func PlayPlaylist(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	}

	playlistName := args[0]

	// Open stored playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}

	scanner := bufio.NewScanner(file)
	var songs []framework.Song
	for scanner.Scan() {
		songId := scanner.Text()

		songs = append(songs, ctx.SongIdList.IdList[songId])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
		return
	}

	file.Close()

	ctx.MediaPlayer.AddSongs(songs...)

	// start the player if its not running
	if !ctx.MediaPlayer.IsConnected {
		go ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
	}

	ctx.ReplyEmbed(playlistName+" queued!", "", 0xFFD900)
}