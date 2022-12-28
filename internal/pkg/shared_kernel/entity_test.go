package sharedkernel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	shared "github.com/thangchung/go-coffeeshop/internal/pkg/shared_kernel"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	id := shared.NewID()
	assert.NotNil(t, id)
}

func TestStringToID(t *testing.T) {
	t.Parallel()

	_, err := shared.StringToID("fd14c028-5f56-488a-8c29-3186fd62395c")
	assert.Nil(t, err)
}
