package devcert

import "testing"

func TestSignCSR_InvalidCSR(t *testing.T) {
	_, _, err := SignCSR([]byte("not pem"), []byte("not"), []byte("not"))
	if err == nil {
		t.Fatal("expected error")
	}
}
