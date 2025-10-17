package sdp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateForAudio(t *testing.T) {
	current := GenerateForAudio(
		net.ParseIP("1.1.1.1"),
		net.ParseIP("2.2.2.2"),
		50000,
		ModeSendrecv,
		[]string{"0", "101", "8", "3"},
	)
	assert.NotNil(t, current)

	currentSd, err := FromString(current)
	assert.Nil(t, err)
	assert.NotNil(t, currentSd)

	withPion := GenAudio(
		net.ParseIP("1.1.1.1"),
		net.ParseIP("2.2.2.2"),
		50000,
		ModeSendrecv,
		[]string{"0", "101", "8", "3"},
	)

	assert.NotNil(t, current)
	pionSd, err := FromString(withPion)
	assert.Nil(t, err)
	assert.NotNil(t, pionSd)

}
