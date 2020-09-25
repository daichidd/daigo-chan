package handler

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/daichidd/daigo-chan/data"
	"github.com/daichidd/daigo-chan/usecase"
)

const (
	// general(一般)
	DEBUG_CHANNEL_ID = "756882537485434952"
	//DAIGO_CHAN_ID    = "756783428023746613"

	DAIGO_COMMAND  = "第五"
	RANDOM_COMMAND = "ランダム"
	HELP_COMMAND   = "教えて"
	LIST_COMMAND   = "一覧"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "" {
		return
	}

	// メンバーが機械に弱いのでregexpで半角全角空白に対応する
	command := regexp.MustCompile("[/ /　]").Split(m.Content, -1)

	// for debug ignore not develop channel
	if command[0] == "第六" {
		if m.ChannelID != DEBUG_CHANNEL_ID {
			return
		}

		// command routing
		if err := commandRouting(s, m, command); err != nil {
			log.Fatal(err)
		}
	}

	if command[0] == "第五" {
		// command routing
		if err := commandRouting(s, m, command); err != nil {
			log.Fatal(err)
		}
	}
}

func commandRouting(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) error {
	switch cmd[1] {
	case HELP_COMMAND:
		s.ChannelMessageSend(m.ChannelID, data.HELP_TEXT)
	case LIST_COMMAND:
		s.ChannelMessageSend(m.ChannelID, data.LIST_TEXT)
	case RANDOM_COMMAND:
		if err := usecase.RandomPicker(s, m, cmd); err != nil {
			return err
		}
	}

	return nil
}
