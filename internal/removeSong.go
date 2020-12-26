package internal

import (
	"../framework"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// Remove song from queue
func RemoveFromQueue(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter index of song!", 0xEB5160)
		return
	}

	if idx, err := strconv.Atoi(args[0]); err != nil {
		log.Println("RemoveFromQueue(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Invalid index!", 0xEB5160)
	} else {
		if songName, err := ctx.MediaPlayer.RemoveFromQueue(idx - 1); err != nil {
			log.Println("RemoveFromQueue(): " + err.Error())
			ctx.ReplyEmbed("Oops!", "Invalid index!", 0xEB5160)
		} else {
			ctx.ReplyEmbed(songName+" Removed!", "", 0x3365A0)
		}
	}
}

// Remove song by name from playlist
func RemoveSongByName(ctx framework.Context) {

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

	// Remove trailing new line
	if _, err := file.WriteString(strings.TrimSuffix(newList.String(), "\n")); err != nil {
		log.Println("RemoveSong(): " + err.Error())
		log.Println("OLD FILE\n")
		log.Println(oldList)
		log.Println("NEW FILE\n")
		log.Println(newList)
		ctx.ReplyEmbed("Oops!", "Couldn't remove song", 0xEB5160)
		return
	}

	ctx.ReplyEmbed(songName+" removed successfully", "", 0x00C49A)
}

// Remove song by index from playlist
func RemoveSongByIndex(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	} else if len(args) < 2 {
		ctx.ReplyEmbed("Oops!", "Enter the song index!", 0xEB5160)
		return
	}

	playlistName := args[0]

	idx, err := strconv.Atoi(args[1])

	if err != nil {
		log.Println("RemoveSongByIndex(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Invalid index!", 0xEB5160)
		return
	}

	// Read playlist
	contents, err := ioutil.ReadFile("./docs/playlists/" + playlistName + ".txt")
	if err != nil {
		log.Println("RemoveSong(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}

	// Check if song is in playlist
	var oldList []string = strings.Split(string(contents), "\n")
	var newList []string = oldList
	if idx < 0 || idx >= len(oldList) {
		ctx.ReplyEmbed("Oops!", "Index out of range", 0xEB5160)
		return
	}

	// Shift all elements after idx left
	if idx+1 != len(newList) {
		copy(newList[idx:], newList[idx+1:])
	}

	// Remove last element
	newList = newList[:len(newList)-1]

	// Rewrite the playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("RemoveSongByIndex(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(strings.Join(newList, "\n")); err != nil {
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