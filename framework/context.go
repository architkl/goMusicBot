package framework

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Context struct {
	Discord     *discordgo.Session
	Guild       *discordgo.Guild
	TextChannel *discordgo.Channel
	User        *discordgo.User
	Message     *discordgo.MessageCreate
	Args        []string
	MediaPlayer *Player
	CmdHandler  *CommandHandler
	SongIdList  *IdListHandler
}

func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, cmdHandler *CommandHandler,
	player *Player, songIdList *IdListHandler) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.CmdHandler = cmdHandler
	ctx.MediaPlayer = player
	ctx.SongIdList = songIdList
	return ctx
}

func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		log.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (ctx Context) ReplyEmbed(title, description string, colour int) {
	if len(title) > 256 {
		log.Println("ReplyEmbed(): Title exceeds 256 character limit")
		title = title[:256]
	}

	if len(description) > 2048 {
		log.Println("ReplyEmbed(): Description exceeds 2048 character limit")

		if description[:5] == "```ml" {
			description = description[:2045] + "```"
		} else {
			description = description[:2048]
		}
	}

	embed := discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       colour,
		Type:        discordgo.EmbedTypeRich,
	}

	_, err := ctx.Discord.ChannelMessageSendEmbed(ctx.TextChannel.ID, &embed)
	if err != nil {
		log.Println("Error whilst sending message,", err)
	}
}