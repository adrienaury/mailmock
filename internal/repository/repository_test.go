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
	id := repository.Len()
	repository.Store("1")
	repository.Store("2")
	repository.Store("3")
	repository.Store("4")
	repository.Store("5")

	slice := repository.All(id, 2)
	assert.Equal(t, []interface{}{"1", "2"}, slice, "")

	slice = repository.All(id, 5)
	assert.Equal(t, []interface{}{"1", "2", "3", "4", "5"}, slice, "")

	slice = repository.All(id, 10)
	assert.Equal(t, []interface{}{"1", "2", "3", "4", "5"}, slice, "")

	slice = repository.All(id+2, 2)
	assert.Equal(t, []interface{}{"3", "4"}, slice, "")

	slice = repository.All(id+2, 5)
	assert.Equal(t, []interface{}{"3", "4", "5"}, slice, "")

	slice = repository.All(id+5, 5)
	assert.Equal(t, []interface{}{}, slice, "")

	slice = repository.All(id+10, 5)
	assert.Equal(t, []interface{}{}, slice, "")
}
