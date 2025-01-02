package utils

func CalculateOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return (page - 1) * limit
}
