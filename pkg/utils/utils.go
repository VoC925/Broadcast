package utils

import (
	"fmt"
	"net/http"
)

// функция на вход принимает слайс интерфейсов, а на выходе отдает слайс string
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

// функция на вход принимает слайс интерфейсов, а на выходе отдает слайс float64
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

// функция для отправки http запросов по адресу uri
// функция отдает структуру http.Response и возможную ошибку
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

// алгоритм бинарного поиска в слайсе arr значения target, в результате функция
// возвращает значение индекса в слайсе
func BinarySearch(arr []string, target string) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := low + (high-low)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}
