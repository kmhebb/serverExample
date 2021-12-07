package os_test

import (
	goos "os"
	"testing"

	"github.com/pborman/uuid"

	"github.com/kmhebb/serverExample/internal/assert"
	"github.com/kmhebb/serverExample/lib/os"
)

func TestGetStringEnv(t *testing.T) {
	assert := assert.New(t)
	key := uuid.New()

	assert.Equals(os.GetStringEnv(key), "")

	goos.Setenv(key, "test")
	assert.Equals(os.GetStringEnv(key), "test")
}

func TestGetIntEnv(t *testing.T) {
	assert := assert.New(t)
	key := uuid.New()

	// Not set -> 0
	assert.Equals(os.GetIntEnv(key), 0)

	// Set to intval -> intval
	goos.Setenv(key, "3")
	assert.Equals(os.GetIntEnv(key), 3)

	// Set to otherval -> panic
	goos.Setenv(key, "4.5")
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		os.GetIntEnv(key)
	}()
	assert.True(panicked)
}
