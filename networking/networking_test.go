package networking_test

import (
	"testing"

	"cs425_mp2/config"
	"cs425_mp2/networking"
)

func TestSend(t *testing.T) {
	err := networking.Send("127.0.0.1:"+config.PORT, []byte("test"))
	if err != nil {
		t.Errorf("Error sending UDP message: \"%s\"", err)
	}
}
