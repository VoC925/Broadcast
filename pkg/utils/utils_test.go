package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testingDataString struct {
	incomingData []any
	expected     []string
}

func getTestingDataString() []testingDataString {
	return []testingDataString{
		{
			incomingData: []any{"hello", "this", "is", "a", "line"},
			expected:     []string{"hello", "this", "is", "a", "line"},
		},
		{
			incomingData: []any{"34", "3d3", "4vf"},
			expected:     []string{"34", "3d3", "4vf"},
		},
	}
}

type testingDataFloat64 struct {
	incomingData []any
	expected     []float64
}

func getTestingDataFloat64() []testingDataFloat64 {
	return []testingDataFloat64{
		{
			incomingData: []any{7., 5., 3., 1.},
			expected:     []float64{7., 5., 3., 1.},
		},
		{
			incomingData: []any{0., .1},
			expected:     []float64{0., .1},
		},
	}
}

func TestConvertingAnyToString(t *testing.T) {
	data := getTestingDataString()
	for _, dataLine := range data {
		got, err := ConvertSliceAnyToSliceString(dataLine.incomingData)
		require.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, dataLine.expected))
	}
}

func TestConvertingAnyToFloat(t *testing.T) {
	data := getTestingDataFloat64()
	for _, dataLine := range data {
		got, err := ConvertSliceAnyToSliceFloat(dataLine.incomingData)
		require.NoError(t, err)
		assert.True(t, reflect.DeepEqual(got, dataLine.expected))
	}
}

type testingSlice struct {
	incomingData  []string
	target        string
	expectedIndex int
}

func getTestingSlice() []testingSlice {
	return []testingSlice{
		{
			incomingData: []string{
				"2024-04-22T00:00",
				"2024-04-22T01:00",
				"2024-04-22T02:00",
				"2024-04-22T03:00",
				"2024-04-22T04:00",
				"2024-04-22T05:00",
				"2024-04-22T06:00",
				"2024-04-22T07:00",
				"2024-04-22T08:00",
				"2024-04-22T09:00",
				"2024-04-22T10:00",
				"2024-04-22T11:00",
			},
			target:        "2024-04-22T01:00",
			expectedIndex: 1,
		},
		{
			incomingData: []string{
				"2024-04-22T00:00",
				"2024-04-22T01:00",
				"2024-04-22T02:00",
				"2024-04-22T03:00",
				"2024-04-22T04:00",
				"2024-04-22T05:00",
				"2024-04-22T06:00",
				"2024-04-22T07:00",
				"2024-04-22T08:00",
				"2024-04-22T09:00",
				"2024-04-22T10:00",
				"2024-04-22T11:00",
			},
			target:        "2024-04-22T07:00",
			expectedIndex: 7,
		},
		{
			incomingData: []string{
				"2024-04-22T00:00",
				"2024-04-22T01:00",
				"2024-04-22T02:00",
				"2024-04-22T03:00",
				"2024-04-22T04:00",
				"2024-04-22T05:00",
				"2024-04-22T06:00",
				"2024-04-22T07:00",
				"2024-04-22T08:00",
				"2024-04-22T09:00",
				"2024-04-22T10:00",
				"2024-04-22T11:00",
			},
			target:        "2024-04-22T012:00",
			expectedIndex: -1,
		},
	}
}

func TestBinarySearch(t *testing.T) {
	arrs := getTestingSlice()
	for _, dataLine := range arrs {
		index := BinarySearch(dataLine.incomingData, dataLine.target)
		assert.Equal(t, dataLine.expectedIndex, index)
	}
}
