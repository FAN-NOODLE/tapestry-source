package ba_kanban

type baPerformanceIndicator struct {
	Name string
}

func NewBaPerformanceIndicator() *baPerformanceIndicator {
	return &baPerformanceIndicator{
		Name: "Ba业绩指标",
	}
}
