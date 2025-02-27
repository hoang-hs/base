package helpers

func ConvertTypesToInterfaces[T any](tt []T) []interface{} {
	var ii []interface{}
	for _, t := range tt {
		ii = append(ii, t)
	}
	return ii
}
