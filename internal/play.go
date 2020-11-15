package internal

import (
	"../framework"
	"log"
	"strings"
)

// Play given song or resume if paused
func PlaySong(ctx framework.Context) {

	args := strings.Join(ctx.Args, " ")

	// Resume playing if no args
	if args == "" {
		if err := CheckSameChannel(ctx); err != nil {
			log.Println(err)
		} else if err := ctx.MediaPlayer.Resume(); err != nil {
			ctx.Reply("No song to play!")
		}

		return
	}

	// search for the song online
	videoId, title := Search(args)

	if videoId == "" {
		ctx.Reply("Song not found")
		return
	}

	// check if song present locally
	storedTitle, ok := ctx.SongIdList.IdList[videoId]
	if !ok {
		// get song from youtube
		if err := Get(videoId); err != nil {
			ctx.Reply("Song not found")
			return
		}

		// convert to mp3
		if err := ConvertMp3(videoId); err != nil {
			ctx.Reply("Song not found")
			return
		}

		// convert to dca
		if err := ConvertDca(videoId); err != nil {
			ctx.Reply("Song not found")
			return
		}

		// update song in program and file system
		storedTitle = title
		if err := ctx.SongIdList.UpdateSongs(videoId, title); err != nil {
			ctx.Reply("Unable to load song")
			return
		}
	}

	// Create Song struct
	song := framework.Song{
		Id:    videoId,
		Title: storedTitle,
	}

	// Add songs to queue
	ctx.MediaPlayer.AddSongs(song)

	// start the player if its not running
	if !ctx.MediaPlayer.IsConnected {
		go ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
	} else if ctx.MediaPlayer.IsPaused {
		ctx.MediaPlayer.Resume()
	}
}