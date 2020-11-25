package internal

import (
	"../framework"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func RemoveSong(ctx framework.Context) {

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

	// Read playlist
	contents, err := ioutil.ReadFile("./docs/playlists/" + playlistName + ".txt")
	if err != nil {
		log.Println("RemoveSong(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}

	// Search for song on youtube
	videoId := Search(songName)

	if videoId == "" {
		ctx.ReplyEmbed("Oops!", "Song not found", 0xEB5160)
		return
	}

	// Check if song present locally
	if _, ok := ctx.SongIdList.IdList[videoId]; !ok {
		ctx.ReplyEmbed("Oops!", "Song not found", 0xEB5160)
	}

	// Check if song is in playlist
	var oldList []string = strings.Split(string(contents), "\n")
	var newList strings.Builder
	var removed bool
	for _, songId := range oldList {
		if songId != videoId {
			newList.WriteString(songId + "\n")
		} else {
			removed = true
		}
	}

	if !removed {
		ctx.ReplyEmbed("Oops!", "Song not found", 0xEB5160)
		return
	}

	// Rewrite the playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("RemoveSong(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(newList.String()); err != nil {
		log.Println("RemoveSong(): " + err.Error())
		log.Println("OLD FILE\n")
		log.Println(oldList)
		log.Println("NEW FILE\n")
		log.Println(newList)
		ctx.ReplyEmbed("Oops!", "Couldn't remove song", 0xEB5160)
		return
	}

	ctx.ReplyEmbed("Song removed successfully", "", 0x00C49A)
}