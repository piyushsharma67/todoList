package structs

import (
	"errors"
	"time"
)

type TodoItem struct {
	Id         string `json:"id"`
	Item       string `json:"item"`
	Priority   int `json:"priority"`
	Updated_at time.Time `json:"created_at"`
	Created_at time.Time `json:"updated_at"`
}

type TodoItemList struct {
	Items []TodoItem
	Count int
}

func (t *TodoItem) Validate() error {
	

	if t.Item == "" {
		return errors.New("item is required")
	}

	if t.Priority <= 0 {
		return errors.New("priority cannot be less than 0")
	}

	return nil
}