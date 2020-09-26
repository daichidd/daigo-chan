package model

import (
	"log"
	"math/rand"
	"time"

	"github.com/daichidd/daigo-chan/entity"
)

const (
	SURVIVORS_NUM = 4
)

type Survivor struct {
	*entity.Character
}

type Survivors []*Survivor

func NewSurvivor(character *entity.Character) *Survivor {
	return &Survivor{
		character,
	}
}

func NewSurvivors(characters []*entity.Character) []*Survivor {
	survivors := []*Survivor{}
	for _, v := range characters {
		survivors = append(survivors, NewSurvivor(v))
	}

	return survivors
}

// サバは四人選出する
func (ms Survivors) RandomPick() []*Survivor {
	survivors := []*Survivor{}
	i := 0

	m := make(map[*Survivor]struct{})
	survivorsRes := make([]*Survivor, 0)

	for i < SURVIVORS_NUM {
		survivor := ms.RandomPickOnce()
		survivors = append(survivors, survivor)
		// remove duplicate
		for _, v := range survivors {
			if _, ok := m[v]; !ok {
				m[v] = struct{}{}
				survivorsRes = append(survivorsRes, v)
				i++
			}
		}
	}

	return survivorsRes
}

func (ms Survivors) RandomPickOnce() *Survivor {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(ms))

	// for debug
	log.Println("r: %d, lenms: %d", r, len(ms))
	return ms[r-1]
}
