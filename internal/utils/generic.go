package utils

import (
	"sort"
	"time"
)

// Comparable constraint for types that can be compared
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// SortBy provides generic sorting functionality for any map with integer keys
func SortBy[T any, K Comparable](items map[int]T, keyExtractor func(T) K, reverse bool) []T {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value T
	}, 0, len(items))

	for k, v := range items {
		pairs = append(pairs, struct {
			Key   int
			Value T
		}{k, v})
	}

	// Sort the slice by the extracted key
	sort.SliceStable(pairs, func(i, j int) bool {
		key1, key2 := keyExtractor(pairs[i].Value), keyExtractor(pairs[j].Value)
		if reverse {
			return key1 > key2
		}
		return key1 < key2
	})

	// Extract sorted values
	result := make([]T, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

// Transform applies a transformation function to each element in a slice
func Transform[TSource, TTarget any](source []TSource, transformer func(TSource) TTarget) []TTarget {
	result := make([]TTarget, len(source))
	for i, item := range source {
		result[i] = transformer(item)
	}
	return result
}

// ConvertSlice converts a slice of pointers from one type to another using a converter function
func ConvertSlice[TSource, TTarget any](source []*TSource, converter func(*TSource) TTarget) []TTarget {
	result := make([]TTarget, len(source))
	for i, item := range source {
		result[i] = converter(item)
	}
	return result
}

// FilterSlice filters a slice based on a predicate function
func FilterSlice[T any](source []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range source {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// MapSliceToMap converts a slice to a map using key and value extractors
func MapSliceToMap[T any, K comparable, V any](
	source []T,
	keyExtractor func(T) K,
	valueExtractor func(T) V,
) map[K]V {
	result := make(map[K]V, len(source))
	for _, item := range source {
		key := keyExtractor(item)
		value := valueExtractor(item)
		result[key] = value
	}
	return result
}

// ParseTimeString provides generic time parsing with error handling
func ParseTimeString[T any](timeStr string, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// MaxLen returns the maximum length that should be processed for a given slice
func MaxLen[T any](items []T, maxItems int) int {
	if len(items) > maxItems {
		return maxItems
	}
	return len(items)
}

// StringExtractor defines a function type for extracting strings from any type
type StringExtractor[T any] func(T) string

// NumberExtractor defines a function type for extracting comparable numbers from any type
type NumberExtractor[T any, N Comparable] func(T) N

// DateExtractor defines a function type for extracting time values from any type
type DateExtractor[T any] func(T) time.Time