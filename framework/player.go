package framework

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	Queue       []Song
	idx         int
	Current     *Song
	voice       *discordgo.VoiceConnection
	IsConnected bool
	IsPaused    bool
	mu          sync.Mutex // lock for playing audio
}

// create new player object
func NewMediaPlayer() *Player {
	player := new(Player)
	player.IsPaused = true
	player.Queue = make([]Song, 0)
	player.mu.Lock()
	return player
}

// Connect to voice channel
func (player *Player) ConnectVoice(session *discordgo.Session, guild *discordgo.Guild, authorID string) {

	// check if already connected
	if player.IsConnected {
		log.Println("Player.ConnectVoice(): Already connected")
		return
	}

	// Look for the message sender in that guild's current voice states.
	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == authorID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	// Join the provided voice channel.
	vc, err := session.ChannelVoiceJoin(guild.ID, voiceChannelID, false, true)
	if err != nil {
		log.Println(err)
		return
	}

	// assign voice channel to player
	player.voice = vc

	// set the player to connected
	player.IsConnected = true
}

// Disconnect from voice channel
func (player *Player) DisconnectVoice() {

	// Check if not connected
	if !player.IsConnected {
		log.Println("Player.DisconnectVoice(): Not connected")
		return
	}

	player.IsConnected = false

	if player.IsPaused {
		// Resume so that playSound can return
		player.Resume()
	}
}

// Add songs to the player
func (player *Player) AddSongs(songs ...Song) {
	player.Queue = append(player.Queue, songs...)
}

// Start playing audio
func (player *Player) StartPlaying(session *discordgo.Session, guild *discordgo.Guild, authorID string) {

	// connect to the voice channel if not connected
	if !player.IsConnected {
		player.ConnectVoice(session, guild, authorID)
		if player.voice == nil {
			return
		}
	}

	// check if index is at end of queue
	if player.idx == len(player.Queue) {
		player.idx = 0
	}

	// Start speaking.
	player.voice.Speaking(true)
	player.mu.Unlock()
	player.IsPaused = false

	// loop through the queue
	for idx := &player.idx; *idx < len(player.Queue); *idx++ {
		song := player.Queue[*idx]
		var buffer = make([][]byte, 0)

		// load song data from file system
		err := loadSound(song.Id, &buffer)
		if err != nil {
			log.Println("Player.StartPlaying(): Error loading sound: ", err)
			continue
		}

		player.playSound(buffer)

		// check if player still connected
		if !player.IsConnected {
			break
		}
	}

	// Stop speaking
	player.voice.Speaking(false)
	player.IsPaused = true
	player.mu.Lock()

	// Disconnect from the provided voice channel.
	player.voice.Disconnect()
}

// playSound plays the current buffer to the provided channel.
func (player *Player) playSound(buffer [][]byte) {

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Send the buffer data.
	for _, buff := range buffer {
		player.mu.Lock()
		if !player.IsConnected {
			// return due to disconnection
			player.mu.Unlock()
			return
		}

		player.voice.OpusSend <- buff
		player.mu.Unlock()
	}

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)
}

// Resume playing
func (player *Player) Resume() error {

	// Check if already playing
	if !player.IsPaused {
		log.Println("Player.Play(): Already playing")
		return nil
	}

	if player.idx == len(player.Queue) {
		log.Println("Player.Play(): At end of queue")
		return errors.New("No song to play!")
	}

	player.voice.Speaking(true)
	player.IsPaused = false
	player.mu.Unlock()
	return nil
}

// Pause by locking mutex
func (player *Player) Pause() {

	// Check if already paused
	if player.IsPaused {
		log.Println("Player.Pause(): Already paused")
		return
	}

	player.mu.Lock()
	player.voice.Speaking(false)
	player.IsPaused = true
}
