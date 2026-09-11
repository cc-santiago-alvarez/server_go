package domain

import "fmt"

type CreateParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (p CreateParams) Validation() error {
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price < 0 {
		return fmt.Errorf("product price is required")
	}

	return nil
}

type UpdateParams struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Status      DemoStatus `json:"status"`
}
