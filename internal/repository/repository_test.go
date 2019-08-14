package repository_test

import (
	"testing"

	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryNominal(t *testing.T) {
	in := "test"
	id := repository.Store("test")
	out := repository.Use(id)
	assert.Equal(t, in, out, "")
}

func TestRepositoryNil(t *testing.T) {
	out := repository.Use(9999)
	assert.Nil(t, out, "")
}

func TestRepositoryAll(t *testing.T) {
	repository.Reset()
	repository.Store("1")
	repository.Store("2")
	repository.Store("3")
	repository.Store("4")
	repository.Store("5")

	len := repository.Len()
	assert.Equal(t, 5, len, "")

	slice, full := repository.All(0, 2)
	assert.Equal(t, map[int]interface{}{0: "1", 1: "2"}, slice, "")
	assert.Equal(t, false, full, "")

	slice, full = repository.All(0, 5)
	assert.Equal(t, map[int]interface{}{0: "1", 1: "2", 2: "3", 3: "4", 4: "5"}, slice, "")
	assert.Equal(t, true, full, "")

	slice, full = repository.All(0, 10)
	assert.Equal(t, map[int]interface{}{0: "1", 1: "2", 2: "3", 3: "4", 4: "5"}, slice, "")
	assert.Equal(t, true, full, "")

	slice, full = repository.All(2, 2)
	assert.Equal(t, map[int]interface{}{2: "3", 3: "4"}, slice, "")
	assert.Equal(t, false, full, "")

	slice, full = repository.All(2, 5)
	assert.Equal(t, map[int]interface{}{2: "3", 3: "4", 4: "5"}, slice, "")
	assert.Equal(t, false, full, "")

	slice, full = repository.All(5, 5)
	assert.Equal(t, map[int]interface{}{}, slice, "")
	assert.Equal(t, false, full, "")

	slice, full = repository.All(10, 5)
	assert.Nil(t, slice, "")
	assert.Equal(t, false, full, "")
}
