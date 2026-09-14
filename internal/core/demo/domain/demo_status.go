package domain

type DemoStatus string

const (
	DemoStatusActive   DemoStatus = "active"
	DemoStatusInactive DemoStatus = "inactive"
)

//getters

func (s DemoStatus) String() string {
	return string(s)
}

// validators
func (s DemoStatus) IsActive() bool {
	return s == DemoStatusActive
}

func (s DemoStatus) IsInactive() bool {
	return s == DemoStatusInactive
}

func (s DemoStatus) IsValid() bool {
	switch s {
	case DemoStatusActive, DemoStatusInactive:
		return true
	}
	return false
}

func (s DemoStatus) CanTransitionTo(next DemoStatus) bool {
	if !next.IsValid() {
		return false
	}
	return s != next
}

// type ListDemoStatus struct {
// 	Label string `json:"name"`
// 	Value string `json:"value"`
// 	Color string `json:"color"`
// }

// func ListStatus() []ListDemoStatus {
// 	return []ListDemoStatus{
// 		ListDemoStatus{
// 			Label: DemoStatusActive.ToString(),
// 			Value: DemoStatusInactive.ToString(),
// 		},
// 		ListDemoStatus{
// 			Label: DemoStatusActive.ToString(),
// 			Value: DemoStatusInactive.ToString(),
// 		},
// 	}
// }
