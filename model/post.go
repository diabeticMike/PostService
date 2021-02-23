package model

import (
	"encoding/json"
	"time"
)

// Post entity
type Post struct {
	Name   string    `json:"post_name"`
	Date   time.Time `json:"date"`
	Author string    `json:"author"`
}

// MarshalJSON needed for formatting date parameter
func (p *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(p),
		Date:  time.Time(p.Date).Format("02.01.06"),
	})
}

type Posts []Post

func (d Posts) Len() int           { return len(d) }
func (d Posts) Less(i, j int) bool { return d[i].Date.After(d[j].Date) }
func (d Posts) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
