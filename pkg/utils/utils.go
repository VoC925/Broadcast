package utils

import (
	"fmt"
	"net/http"
)

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

func GetRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	// запрос к API
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
