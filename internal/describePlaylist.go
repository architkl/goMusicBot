package internal

import (
	"../framework"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func DescribePlaylist(ctx framework.Context) {

	args := ctx.Args

	if len(args) < 1 {
		ctx.ReplyEmbed("Oops!", "Enter the playlist name!", 0xEB5160)
		return
	}

	playlistName := args[0]

	// Open stored playlist
	file, err := os.OpenFile("./docs/playlists/"+playlistName+".txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Println("DescribePlaylist(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Playlist not found", 0xEB5160)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var songs []string
	var buf strings.Builder
	var idx int = 1
	buf.WriteString("```ml\n")
	for scanner.Scan() {
		songId := scanner.Text()
		song := ctx.SongIdList.IdList[songId]

		if len(song.Title) > 40 {
			buf.WriteString(strconv.Itoa(idx) + ") " + song.Title[:40] + "... '" + song.Duration + "'\n")
		} else {
			buf.WriteString(strconv.Itoa(idx) + ") " + song.Title + " '" + song.Duration + "'\n")
		}

		idx++

		if len(buf.String()) > 2000 {
			buf.WriteString("```")
			songs = append(songs, buf.String())
			buf.Reset()
			buf.WriteString("```ml\n")
		}
	}
	buf.WriteString("```")
	songs = append(songs, buf.String())

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
		return
	}

	if songs[0] == "```ml\n```" {
		ctx.ReplyEmbed(playlistName, "```ml\nEmpty Playlist```", 0x228FD3)
		return
	}

	for _, text := range songs {
		ctx.ReplyEmbed(playlistName, text, 0x228FD3)
	}
}