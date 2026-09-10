package productDemo_domain

type ProductStatus string

const (
	ProductStatusActive ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

//getters

func (s ProductStatus) ToString() string {
	return string(s)
}

//validators
func (s ProductStatus) IsActive() bool {
	return s == ProductStatusActive
}

func (s ProductStatus) IsInactive() bool {
	return s == ProductStatusInactive
}

func (s *ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusActive, ProductStatusInactive:
	return true
	}
	return false
}

type ListProductStatus struct {
	Label string `json:"name"`
	Value string `json:"value"`
	Color string `json:"color"`
}

func ListStatus() []ListProductStatus {
	return []ListProductStatus{
		ListProductStatus{
			Label: ProductStatusActive.ToString(),
			Value: ProductStatusActive.ToString(),
		},
		ListProductStatus{
			Label: ProductStatusInactive.ToString(),
			Value: ProductStatusInactive.ToString(),
		},
	}
}
