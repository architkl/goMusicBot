package main

import (
	"../framework"
	"../internal"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Command line flags and global variables
var (
	fToken     = flag.String("t", "", "bot token")
	fPrefix    = flag.String("p", "!", "bot prefix")
	CmdHandler *framework.CommandHandler
	botId      string
	player     *framework.Player
)

func main() {
	flag.Parse()

	// channel to receive signal to shutdown bot
	sc := make(chan os.Signal, 1)
	CmdHandler = framework.NewCommandHandler()
	player = framework.NewMediaPlayer()
	registerCommands(sc)

	dg, err := discordgo.New("Bot " + *fToken)
	if err != nil {
		log.Fatal(err)
	}

	usr, err := dg.User("@me")
	if err != nil {
		log.Println("Error obtaining account details,", err)
		return
	}
	botId = usr.ID

	dg.AddHandler(commandHandler)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}

	content := message.Content
	if len(content) <= len(*fPrefix) {
		return
	}
	if content[:len(*fPrefix)] != *fPrefix {
		return
	}
	content = content[len(*fPrefix):]
	if len(content) < 1 {
		return
	}

	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		log.Println("Error getting channel,", err)
		return
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		log.Println("Error getting guild,", err)
		return
	}

	ctx := framework.NewContext(discord, guild, channel, user, message, CmdHandler, player)
	args := strings.Fields(content)
	ctx.Args = args[1:]

	name := strings.ToLower(args[0])
	middleware, command, found := CmdHandler.Get(name)
	if !found {
		return
	}

	m := *middleware
	if err := m(*ctx); err != nil {
		log.Println(err)
		return
	}

	c := *command
	c(*ctx)
}

func registerCommands(sc chan os.Signal) {
	CmdHandler.Register("ping", internal.Logging, internal.Ping, "respongs")
	CmdHandler.Register("avatar", internal.Logging, internal.Avatar, "returns user's avatar")
	CmdHandler.Register("user", internal.Logging, internal.Username, "returns user's username")
	CmdHandler.Register("pl", internal.CheckVoice, internal.PlayPlaylist, "play the given playlist")
	CmdHandler.Register("shutdown", internal.Logging, func(ctx framework.Context) {
		ctx.Reply("Bye!")
		sc <- os.Interrupt
	}, "shutdown the bot")
}