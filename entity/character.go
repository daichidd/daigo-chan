package entity

import (
	"encoding/csv"
)

type Character struct {
	NickName string
	Name     string
	Sex      bool
	Type     string
}

func SetCharacters(data *csv.Reader, characterType string) ([]*Character, error) {
	r, err := data.ReadAll()
	if err != nil {
		return nil, err
	}

	characters := []*Character{}
	for k, v := range r {
		if k == 0 {
			continue
		}

		characters = append(characters, &Character{
			NickName: v[0],
			Name:     v[1],
			Sex:      v[2] == "male",
			Type:     characterType,
		})
	}

	return characters, nil
}
