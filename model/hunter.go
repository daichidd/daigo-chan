package model

import (
	"math/rand"
	"time"

	"github.com/daichidd/daigo-chan/entity"
)

type Hunter struct {
	*entity.Character
}

type Hunters []*Hunter

func NewHunter(character *entity.Character) *Hunter {
	return &Hunter{
		character,
	}
}

func NewHunters(characters []*entity.Character) []*Hunter {
	hunters := []*Hunter{}
	for _, v := range characters {
		hunters = append(hunters, NewHunter(v))
	}

	return hunters
}

func (ms Hunters) RandomPick() *Hunter {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(ms))

	return ms[r-1]
}
