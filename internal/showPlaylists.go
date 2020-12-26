package internal

import (
	"architkl/goMusicBot/framework"
	"io/ioutil"
	"log"
	"strings"
)

func ShowPlaylists(ctx framework.Context) {

	// Read playlist directory
	files, err := ioutil.ReadDir("./docs/playlists/")
	if err != nil {
		log.Println("ShowPlaylists(): " + err.Error())
		ctx.ReplyEmbed("Oops!", "Couldn't find playlists", 0xEB5160)
		return
	}

	var playlists strings.Builder
	for _, f := range files {
		// Remove .txt from name
		playlists.WriteString(f.Name()[:len(f.Name())-4] + "\n")
	}

	ctx.ReplyEmbed("Playlists", "```css\n"+playlists.String()+"```", 0x228FD3)
}