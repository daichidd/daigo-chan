package model

import (
	"math/rand"
	"time"

	"github.com/daichidd/daigo-chan/entity"
)

type Map struct {
	*entity.Map
}

type Maps []*Map

func NewMap(m *entity.Map) *Map {
	return &Map{
		m,
	}
}

func NewMaps(ms []*entity.Map) []*Map {
	maps := []*Map{}
	for _, v := range ms {
		maps = append(maps, NewMap(v))
	}

	return maps
}

func (ms Maps) RandomPick(isRank bool) *Map {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(ms))

	res := ms[r-1]

	if isRank {
		if !res.IsRank {
			res = ms.RandomPick(true)
		}
	}

	return res
}
