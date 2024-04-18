package utils

import "fmt"

func ConvertSliceAnyToSliceString(slice []any) ([]string, error) {
	sliceOfString := make([]string, 0, len(slice))
	for _, elem := range slice {
		str, ok := elem.(string)
		if !ok {
			return nil, fmt.Errorf(`type assertion failed: %v to string`, elem)
		}
		sliceOfString = append(sliceOfString, str)
	}
	return sliceOfString, nil
}

func ConvertSliceAnyToSliceFloat(slice []any) ([]float64, error) {
	sliceOfFloat := make([]float64, 0, len(slice))
	for _, elem := range slice {
		valFloat, ok := elem.(float64)
		if !ok {
			return nil, fmt.Errorf(`type assertion failed: %v to float64`, elem)
		}
		sliceOfFloat = append(sliceOfFloat, valFloat)
	}
	return sliceOfFloat, nil
}
