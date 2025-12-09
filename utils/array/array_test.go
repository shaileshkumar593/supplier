package array_test

import (
	"testing"

	"swallow-supplier/utils/array"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
	// Test case 1: Existing value in slice
	exists, index := array.InArray(3, []int{1, 2, 3, 4})
	assert.True(t, exists)
	assert.Equal(t, 2, index)

	// Test case 2: Non-existing value in slice
	exists, index = array.InArray(5, []int{1, 2, 3, 4})
	assert.False(t, exists)
	assert.Equal(t, -1, index)

	// Test case 3: Existing value in empty slice
	exists, index = array.InArray("hello", []string{})
	assert.False(t, exists)
	assert.Equal(t, -1, index)

	// Test case 4: Existing value in slice with mixed types
	exists, index = array.InArray(2, []interface{}{"hello", 2, true})
	assert.True(t, exists)
	assert.Equal(t, 1, index)
}

func TestKeyExist(t *testing.T) {
	// Test case 1: Existing key in map
	exists := array.KeyExist("key1", map[string]string{"key1": "value1", "key2": "value2"})
	assert.True(t, exists)

	// Test case 2: Non-existing key in map
	exists = array.KeyExist("key3", map[string]string{"key1": "value1", "key2": "value2"})
	assert.False(t, exists)

	// Test case 3: Existing key in empty map
	exists = array.KeyExist("key1", map[string]string{})
	assert.False(t, exists)
}

func TestFlip(t *testing.T) {
	// Test case 1: Regular map flip
	input := map[string]string{"key1": "value1", "key2": "value2"}
	expectedOutput := map[string]string{"value1": "key1", "value2": "key2"}
	output := array.Flip(input)
	assert.Equal(t, expectedOutput, output)

	// Test case 2: Empty map flip
	input = map[string]string{}
	expectedOutput = map[string]string{}
	output = array.Flip(input)
	assert.Equal(t, expectedOutput, output)
}

func TestUnique(t *testing.T) {
	// Test case 1: Remove duplicates from a slice
	input := []string{"a", "b", "c", "a", "d", "b"}
	expectedOutput := []string{"a", "b", "c", "d"}
	output := array.Unique(input)
	assert.Equal(t, expectedOutput, output)

	// Test case 2: No duplicates in the slice
	input = []string{"a", "b", "c", "d"}
	expectedOutput = []string{"a", "b", "c", "d"}
	output = array.Unique(input)
	assert.Equal(t, expectedOutput, output)

	// Test case 3: Empty slice
	input = []string{}
	expectedOutput = nil
	output = array.Unique(input)
	assert.Equal(t, expectedOutput, output)
}
