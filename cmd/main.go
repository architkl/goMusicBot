package main

import (
	"architkl/goMusicBot/framework"
	"architkl/goMusicBot/internal"
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
	fPrefix    = flag.String("p", "+", "bot prefix")
	CmdHandler = framework.NewCommandHandler()
	player     = framework.NewMediaPlayer()
	sc         = make(chan os.Signal, 1) // channel to receive signal to shutdown bot
	songIdList = framework.NewIdListHandler()
	botId      string
)

func main() {

	flag.Parse()
	songIdList.LoadSongs()
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

// Check user messages for valid commands
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// Recover in case a command panics
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from: ", r)
		}
	}()

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

	ctx := framework.NewContext(discord, guild, channel, user, message, CmdHandler, player, songIdList)
	args := strings.Fields(content)
	ctx.Args = args[1:]

	name := strings.ToLower(args[0])
	middleware, command, found := CmdHandler.Get(name)
	if !found {
		return
	}

	// run middleware
	m := *middleware
	if err := m(*ctx); err != nil {
		log.Println(err)
		return
	}

	c := *command
	c(*ctx)
}

// Register user commands
func registerCommands(sc chan os.Signal) {

	// Basic commands
	CmdHandler.Register("ping", internal.Logging, internal.Ping, "Respongs")
	CmdHandler.Register("avatar", internal.Logging, internal.Avatar, "Show avatar")
	CmdHandler.Register("user", internal.Logging, internal.Username, "Show username")
	CmdHandler.Register("help", internal.Logging, internal.Help, "Show description for each command")

	// Media player commands
	CmdHandler.Register("p", internal.CanPlay, internal.PlaySong, "Play the given song")
	CmdHandler.Register("pl", internal.CanPlay, internal.PlayPlaylist, "Play the given playlist")
	CmdHandler.Register("np", internal.CheckSameChannel, internal.NowPlaying, "Show current song")
	CmdHandler.Register("pause", internal.CheckSameChannel, internal.Pause, "Pause")
	CmdHandler.Register("shuffle", internal.CheckSameChannel, internal.Shuffle, "Shuffle queue")
	CmdHandler.Register("n", internal.CheckSameChannel, internal.NextSong, "Play the next song")
	CmdHandler.Register("r", internal.CheckSameChannel, internal.RemoveFromQueue, "Remove song from queue")
	CmdHandler.Register("m", internal.CheckSameChannel, internal.MovePosition, "Move song position by index")
	CmdHandler.Register("cl", internal.CanPlay, internal.ClearQueue, "Clear the queue")
	CmdHandler.Register("dc", internal.CheckSameChannel, internal.Disconnect, "Disconnect the player")

	// Playlist commands
	CmdHandler.Register("newpl", internal.Logging, internal.CreatePlaylist, "Create a new playlist")
	// CmdHandler.Register("removepl", internal.Logging, internal.RemovePlaylist, "Remove playlist")
	CmdHandler.Register("showpl", internal.Logging, internal.ShowPlaylists, "Show playlists")
	CmdHandler.Register("despl", internal.Logging, internal.DescribePlaylist, "Describe a playlist")
	CmdHandler.Register("ads", internal.Logging, internal.AddSong, "Add song to playlist")
	CmdHandler.Register("rs", internal.Logging, internal.RemoveSongByIndex, "Remove song by index from playlist")
	CmdHandler.Register("rsn", internal.Logging, internal.RemoveSongByName, "Remove song by name from playlist")
	CmdHandler.Register("q", internal.Logging, internal.DisplayQueue, "Show songs in queue")

	// Shutdown
	CmdHandler.Register("shutdown", internal.Logging, func(ctx framework.Context) {
		ctx.ReplyEmbed("Bye!", "", 0x00C49A)
		sc <- os.Interrupt
	}, "Shutdown the bot")
}