package internal

import (
	"../framework"
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

func AddSong(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	} else if len(args) < 2 {
		ctx.ReplyEmbed("Oops!", "Enter the song name!", 0xEB5160)
		return
	}

	playlistName := args[0]
	songName := strings.Join(args[1:], " ")

	// Open stored playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Println("AddSong(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}
	defer file.Close()

	// Get song details
	var song framework.Song
	song.Id, song.Title, song.Duration = GetMetaData(songName)

	if song.Id == "" {
		ctx.ReplyEmbed("Oops!", songName+" not found", 0xEB5160)
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
			log.Println("AddSong(): " + err.Error())
			ctx.ReplyEmbed("Oops!", err.Error(), 0xEB5160)
			return
		}

		// Update song in program and file system
		if err := ctx.SongIdList.UpdateSongs(song); err != nil {
			log.Println("AddSong(): " + err.Error())
			ctx.ReplyEmbed("Oops!", "Unable to load song", 0xEB5160)
			return
		}
	} else {
		// Check if song already in playlist
		var duplicate bool
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			songId := scanner.Text()

			if songId == song.Id {
				duplicate = true
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("AddSong(): " + err.Error())
			return
		}

		if duplicate {
			ctx.ReplyEmbed("Oops!", "Song already present!", 0xEB5160)
			return
		}
	}

	// Add song to playlist
	if _, err := file.WriteString(song.Id + "\n"); err != nil {
		log.Println("AddSong(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "```prolog\nCouldn't add song to "+playlistName+"```", 0xEB5160)
		return
	}

	ctx.ReplyEmbed("Song added successfully", "```ini\n[ "+song.Title+" ]```", 0x00C49A)
}