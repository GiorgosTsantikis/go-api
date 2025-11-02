package model

import "fmt"

type MenuItem struct {
	ID            int             `json:"id"`
	StoreId       int32           `json:"store_id"`
	Name          string          `json:"name"`
	Price         float64         `json:"price"`
	ConfigOptions []*ConfigOption `json:"configOptions"`
}

func (m *MenuItem) String() string {
	return fmt.Sprintf("User[ID=%d, Name=%s, Price=%s, store_id=%s]", m.ID, m.Name, m.Price, m.StoreId)
}

type ConfigOption struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	MaxSelectable int       `json:"maxSelectable"`
	Options       []*Option `json:"options"`
}

type Option struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceDelta int    `json:"priceDelta"`
}
