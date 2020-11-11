package framework

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	Queue     []Song
	idx       int
	Current   *Song
	IsRunning bool
}

func NewMediaPlayer() *Player {
	player := new(Player)
	player.Queue = make([]Song, 0)
	return player
}

func (player *Player) AddSongs(songs ...Song) {
	player.Queue = append(player.Queue, songs...)
}

func (player *Player) StartPlaying(session *discordgo.Session, guild *discordgo.Guild, authorID string) {

	// Look for the message sender in that guild's current voice states.
	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == authorID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	// Join the provided voice channel.
	// vc - voice connection
	vc, err := session.ChannelVoiceJoin(guild.ID, voiceChannelID, false, true)
	if err != nil {
		log.Println(err)
		return
	}

	// set the player to running
	player.IsRunning = true

	// check if index is at end of queue
	if player.idx == len(player.Queue) {
		player.idx = 0
	}

	// loop through the queue
	for idx := &player.idx; *idx < len(player.Queue); *idx++ {
		song := player.Queue[*idx]
		var buffer = make([][]byte, 0)
		err := loadSound(song.Id, &buffer)
		if err != nil {
			log.Println("Error loading sound: ", err)
			continue
		}

		playSound(vc, buffer)
	}

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	player.IsRunning = false
}