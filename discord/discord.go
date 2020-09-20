package discord

import "github.com/bwmarrin/discordgo"

func SendMessage(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	s.ChannelMessageSend(m.ChannelID, content)
}
