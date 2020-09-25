package usecase

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/daichidd/daigo-chan/discord"
	"github.com/daichidd/daigo-chan/entity"
	"github.com/daichidd/daigo-chan/model"
)

const (
	CSV_DIR       = "./data"
	HUNTER_BASE   = "hunter"
	SURVIVOR_BASE = "survivor"
	MAP_BASE      = "map"
	DATA_EXT      = ".csv"

	CUSTOM_COMMAND    = "カスタム"
	SURVIVOR_COMMAND  = "サバ"
	IS_RANK_COMMAND   = "殿堂"
	BAN_COMMAND       = "バン"
	FREQUENCY_COMMAND = "使用頻度"
	MAP_COMMAND       = "マップ"
)

var (
	hunters   = []*model.Hunter{}
	survivors = []*model.Survivor{}
	maps      = []*model.Map{}
)

type PickResult struct {
	HunterNickName    string
	SurvivorNickNames []string
	MapName           string
	Count             int
	Type              string
}

func dirWalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			f, err := dirWalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, f...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}

func csvReader(path string) (*csv.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return csv.NewReader(f), nil
}

func setData() error {
	filePaths, err := dirWalk(CSV_DIR)
	if err != nil {
		return err
	}

	for _, filePath := range filePaths {
		basename := filepath.Base(filePath)
		switch basename {
		case HUNTER_BASE + DATA_EXT:
			f, err := csvReader(filePath)
			if err != nil {
				return err
			}
			hunterEntities, err := entity.SetCharacters(f, HUNTER_BASE)
			if err != nil {
				return err
			}
			hunters = model.NewHunters(hunterEntities)

		case SURVIVOR_BASE + DATA_EXT:
			f, err := csvReader(filePath)
			if err != nil {
				return err
			}
			survivorEntities, err := entity.SetCharacters(f, SURVIVOR_BASE)
			if err != nil {
				return err
			}
			survivors = model.NewSurvivors(survivorEntities)

		case MAP_BASE + DATA_EXT:
			f, err := csvReader(filePath)
			if err != nil {
				return err
			}
			mapEntities, err := entity.SetMaps(f)
			if err != nil {
				return err
			}
			maps = model.NewMaps(mapEntities)
		}
	}

	if len(hunters) == 0 || len(survivors) == 0 || len(maps) == 0 {
		return nil
	}

	return nil
}

func (m *PickResult) buildMessage() string {
	mapStr := "マップ: %s \n"
	hunterStr := "ハンター: %s \n"
	survivor1 := "サバイバー1: %s \n"
	survivor2 := "サバイバー2: %s \n"
	survivor3 := "サバイバー3: %s \n"
	survivor4 := "サバイバー4: %s \n"
	frequencyStr := "使用頻度: %d"

	survivorStr := survivor1 + survivor2 + survivor3 + survivor4
	res := ""
	switch m.Type {
	case "":
		fallthrough
	case SURVIVOR_COMMAND:
		res += survivorStr
		res = fmt.Sprintf(res,
			m.SurvivorNickNames[0],
			m.SurvivorNickNames[1],
			m.SurvivorNickNames[2],
			m.SurvivorNickNames[3],
		)
	case CUSTOM_COMMAND, IS_RANK_COMMAND:
		res = mapStr + hunterStr + survivorStr
		res = fmt.Sprintf(res,
			m.MapName,
			m.HunterNickName,
			m.SurvivorNickNames[0],
			m.SurvivorNickNames[1],
			m.SurvivorNickNames[2],
			m.SurvivorNickNames[3],
		)
	case BAN_COMMAND:
		res = hunterStr + survivor1
		res = fmt.Sprintf(res,
			m.HunterNickName,
			m.SurvivorNickNames[0],
		)
	case FREQUENCY_COMMAND:
		res = frequencyStr
		res = fmt.Sprintf(res, m.Count)
	case MAP_COMMAND:
		res = mapStr
		res = fmt.Sprintf(res, m.MapName)
	}

	return res
}

func RandomPicker(s *discordgo.Session, m *discordgo.MessageCreate, cmd []string) error {
	if err := setData(); err != nil {
		return err
	}

	var (
		hunterRes         *model.Hunter
		survivorNickNames []string
		mapRes            *model.Map
		cnt               int
		cmdType           string
	)

	if len(cmd) >= 3 {
		switch cmd[2] {
		case SURVIVOR_COMMAND:
			survivorRes := model.Survivors(survivors).RandomPick()
			// get nicknames をあとでまとめる
			for _, v := range survivorRes {
				survivorNickNames = append(survivorNickNames, v.NickName)
			}

		case CUSTOM_COMMAND:
			hunterRes = model.Hunters(hunters).RandomPick()
			survivorRes := model.Survivors(survivors).RandomPick()
			// get nicknames をあとでまとめる
			for _, v := range survivorRes {
				survivorNickNames = append(survivorNickNames, v.NickName)
			}

			mapRes = model.Maps(maps).RandomPick(false)
		case IS_RANK_COMMAND:
			hunterRes = model.Hunters(hunters).RandomPick()
			survivorRes := model.Survivors(survivors).RandomPick()
			// get nicknames をあとでまとめる
			for _, v := range survivorRes {
				survivorNickNames = append(survivorNickNames, v.NickName)
			}
			mapRes = model.Maps(maps).RandomPick(true)

		case BAN_COMMAND:
			hunterRes = model.Hunters(hunters).RandomPick()
			survivorRes := model.Survivors(survivors).RandomPickOnce()
			survivorNickNames = append(survivorNickNames, survivorRes.NickName)
		case FREQUENCY_COMMAND:
			// hunterの方が数が少ない
			cnt = len(hunters)
			if cnt > len(survivors) {
				cnt = len(survivors)
			}

			rand.Seed(time.Now().UnixNano())
			cnt = rand.Intn(cnt)
		case MAP_COMMAND:
			// 殿堂のみ。殿堂制限なしはゲーム側に実装されている
			mapRes = model.Maps(maps).RandomPick(true)
		}

		cmdType = cmd[2]
	} else {
		survivorRes := model.Survivors(survivors).RandomPick()
		// get nicknames をあとでまとめる
		for _, v := range survivorRes {
			survivorNickNames = append(survivorNickNames, v.NickName)
		}
	}

	hunterNickName := ""
	if hunterRes != nil {
		hunterNickName = hunterRes.NickName
	}
	mapName := ""
	if mapRes != nil {
		mapName = mapRes.Name
	}

	res := &PickResult{
		HunterNickName:    hunterNickName,
		SurvivorNickNames: survivorNickNames,
		MapName:           mapName,
		Count:             cnt,
		Type:              cmdType,
	}

	discord.SendMessage(s, m, res.buildMessage())

	return nil
}
