package mail

import "testing"

func TestSendSMTP_SkipWhenNoHost(t *testing.T) {
	if err := SendSMTP("", "587", "", "", "f", []string{"t@x"}, "s", "b"); err != nil {
		t.Fatal(err)
	}
}

func TestSendSMTP_SkipWhenNoRecipients(t *testing.T) {
	if err := SendSMTP("h", "587", "", "", "f", nil, "s", "b"); err != nil {
		t.Fatal(err)
	}
}
