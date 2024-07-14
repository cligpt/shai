package drive

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrive(t *testing.T) {
	m := New(context.Background(), DefaultConfig())
	assert.NotEqual(t, nil, m)
}
