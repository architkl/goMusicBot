package internal

import (
	"architkl/goMusicBot/framework"
	"architkl/goMusicBot/pkg"
	"log"
	"strconv"
	"strings"
)

func WriteText(idx, curPos int, sb *strings.Builder, song framework.Song) {
	var txtLimit = 35

	if idx == curPos {
		if len(song.Title) > txtLimit {
			sb.WriteString("> " + strconv.Itoa(idx+1) + ") \"" + song.Title[:txtLimit] + "...\" '" + song.Duration + "'\n")
		} else {
			sb.WriteString("> " + strconv.Itoa(idx+1) + ") \"" + song.Title + "\" '" + song.Duration + "'\n")
		}
	} else {
		if len(song.Title) > txtLimit {
			sb.WriteString(strconv.Itoa(idx+1) + ") " + song.Title[:txtLimit] + "... '" + song.Duration + "'\n")
		} else {
			sb.WriteString(strconv.Itoa(idx+1) + ") " + song.Title + " '" + song.Duration + "'\n")
		}
	}
}

func DisplayQueue(ctx framework.Context) {

	if len(ctx.MediaPlayer.Queue) == 0 {
		ctx.ReplyEmbed("Queue", "```ml\nEmpty```", 0x5539CC)
		return
	}

	curPos, err := ctx.MediaPlayer.CurPosition()
	if err != nil {
		log.Println("DisplayQueue(): ", err.Error())
		ctx.ReplyEmbed("Oops!", "```prolog\nCouldn't get queue```", 0xEB5160)
		return
	}

	var qlen int = len(ctx.MediaPlayer.Queue)
	var queue strings.Builder
	queue.WriteString("```ml\n")

	// Avoid exceeding embed description length by checking queue size
	if qlen > 40 {
		// Show first 3 songs
		for idx := 0; idx < 3; idx++ {
			WriteText(idx, curPos, &queue, ctx.MediaPlayer.Queue[idx])
		}
		queue.WriteString(".\n.\n.\n")

		// Show 2 songs before and after current
		for idx := pkg.Max(3, curPos-2); idx < pkg.Min(curPos+3, qlen); idx++ {
			WriteText(idx, curPos, &queue, ctx.MediaPlayer.Queue[idx])
		}

		if curPos+3 < qlen {
			queue.WriteString(".\n.\n.\n")
		}

		// Show last 3 songs
		for idx := pkg.Max(curPos+3, qlen-3); idx < qlen; idx++ {
			WriteText(idx, curPos, &queue, ctx.MediaPlayer.Queue[idx])
		}
	} else {
		for idx, song := range ctx.MediaPlayer.Queue {
			WriteText(idx, curPos, &queue, song)
		}
	}
	queue.WriteString("```")

	ctx.ReplyEmbed("Queue", queue.String(), 0x5539CC)
}