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
	SAYU_ID  = "474138503719157760"
	JERRY_ID = "576323553519992832"
	//LULUBELL_ID = "627705751460118539"

	DAIGO_COMMAND       = "第五"
	DEBUG_DAIGO_COMMAND = "第六"
	RANDOM_COMMAND      = "ランダム"
	HELP_COMMAND        = "教えて"
	LIST_COMMAND        = "一覧"
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
	call := callCommandChecker(command[0])

	// for debug ignore not develop channel
	if call == DEBUG_DAIGO_COMMAND {
		if m.ChannelID != DEBUG_CHANNEL_ID {
			return
		}

		// command routing
		if err := commandRouting(s, m, command); err != nil {
			log.Fatal(err)
		}
	}

	if call == DAIGO_COMMAND {
		// さゆ特別対応
		if m.Author.ID == SAYU_ID {
			s.ChannelMessageSend(m.ChannelID, data.SAYU_TEXT)
			return
		}
		// ジェリー特別対応
		if m.Author.ID == JERRY_ID {
			s.ChannelMessageSend(m.ChannelID, data.JERRY_TEXT)
			return
		}
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

// ユーザーが機械得意ではないので、呼び出しコマンドの数字の半角全角もまとめて変換する
// bytes.Readerの実装した方がよさそうだけどそこまでこだわるものでもないのでいつか
// ともかくあとで整理する
func callCommandChecker(cmd string) string {
	// とりあえず一旦ベタ書き
	if cmd == "第5" || cmd == "第５" {
		return DAIGO_COMMAND
	}
	return cmd
}
