package internal

import (
	"../framework"
	"strings"
)

func Play(ctx framework.Context) {
	args := strings.Join(ctx.Args, " ")

	if args == "" {
		ctx.Reply("Enter the song name!")
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

	song := framework.Song{
		Id:    videoId,
		Title: storedTitle,
	}

	ctx.MediaPlayer.AddSongs(song)

	// start the player if its not running
	if !ctx.MediaPlayer.IsRunning {
		ctx.MediaPlayer.StartPlaying(ctx.Discord, ctx.Guild, ctx.Message.Author.ID)
	}
}