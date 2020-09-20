package entity

import (
	"encoding/csv"
	"strconv"
)

type Map struct {
	Name   string
	IsRank bool
}

func SetMaps(data *csv.Reader) ([]*Map, error) {
	r, err := data.ReadAll()
	if err != nil {
		return nil, err
	}

	maps := []*Map{}
	for k, v := range r {
		if k == 0 {
			continue
		}

		isRank, err := strconv.ParseBool(v[1])
		if err != nil {
			return nil, err
		}
		maps = append(maps, &Map{
			Name:   v[0],
			IsRank: isRank,
		})
	}

	return maps, nil
}
