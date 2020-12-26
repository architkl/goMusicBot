package internal

import (
	"architkl/goMusicBot/framework"
	"sort"
	"strings"
)

func Help(ctx framework.Context) {

	var desc []string
	for k, v := range ctx.CmdHandler.GetCmds() {
		desc = append(desc, "+"+k+":"+strings.Repeat(" ", 9-len(k))+v.GetHelp())
	}
	sort.Strings(desc)

	ctx.ReplyEmbed("Help", "```ml\n"+strings.Join(desc, "\n")+"```", 0x228FD3)
}