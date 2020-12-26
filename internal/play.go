package internal

import (
	"architkl/goMusicBot/framework"
	"log"
	"strings"
	"time"
)

// Play given song or resume if paused
func PlaySong(ctx framework.Context) {

	args := strings.Join(ctx.Args, " ")

	// Resume playing if no args
	if args == "" {
		if !ctx.MediaPlayer.IsConnected {
			if len(ctx.MediaPlayer.Queue) == 0 {
				ctx.ReplyEmbed("Oops!", "No song to play!", 0xEB5160)
			} else {
				ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
			}
		} else if err := CheckSameChannel(ctx); err != nil {
			log.Println(err)
		} else if err := ctx.MediaPlayer.Resume(); err != nil {
			ctx.ReplyEmbed("Oops!", "No song to play!", 0xEB5160)
		}

		return
	}

	// Get song details
	var song framework.Song
	song.Id, song.Title, song.Duration = GetMetaData(args)

	if song.Id == "" {
		ctx.ReplyEmbed("Oops!", args+" not found", 0xEB5160)
		return
	}

	// Check if duration is too long
	if d, err := time.ParseDuration(song.Duration); err != nil {
		ctx.ReplyEmbed("Oops!", "Encountered an error", 0xEB5160)
		return
	} else {
		lmt, _ := time.ParseDuration("1h")

		if d > lmt {
			ctx.ReplyEmbed("Oops!", "Song length too long", 0xEB5160)
			return
		}
	}

	// Check if song present locally
	_, ok := ctx.SongIdList.IdList[song.Id]
	if !ok {
		// Get the song in dca format
		if err := GetSong(song.Id); err != nil {
			log.Println("PlaySong(): " + err.Error())
			ctx.ReplyEmbed("Oops!", err.Error(), 0xEB5160)
			return
		}

		// Update song in program and file system
		if err := ctx.SongIdList.UpdateSongs(song); err != nil {
			log.Println("PlaySong(): " + err.Error())
			ctx.ReplyEmbed("Oops!", "Unable to load song", 0xEB5160)
			return
		}
	}

	// Add song to queue
	ctx.MediaPlayer.AddSongs(song)

	// start the player if its not running
	if !ctx.MediaPlayer.IsConnected {
		go ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
	} else if ctx.MediaPlayer.IsPaused {
		ctx.MediaPlayer.Resume()
	}

	ctx.ReplyEmbed(song.Title+" queued!", "", 0xFFD900)
}