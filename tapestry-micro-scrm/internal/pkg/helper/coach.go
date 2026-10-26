package helper

import "strings"

func GetSalesPosition() []string {
	return []string{
		strings.ToLower("Sr Sales Associate"),
		strings.ToLower("Sales Associate"),
		strings.ToLower("Senior Sales"),
		strings.ToLower("Seasonal Sales Associate"),
		strings.ToLower("eSA"),
		strings.ToLower("Barista"),
		strings.ToLower("Caf Supervisor"),
		strings.ToLower("Kitchen Helper"),
		strings.ToLower("Retail Management Trainee SA"),
		strings.ToLower("Retail Management Trainee SSA"),
		strings.ToLower("SALES"),
		strings.ToLower("Sales Assistant"),
		strings.ToLower("Service Crew"),
	}
}
