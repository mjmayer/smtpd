package smtpd_test

import (
	"strings"
	"testing"

	"net/mail"

	"github.com/mailproto/smtpd"
)

const (
	multipartEmail = `From: Sender <sender@example.com>
Date: Mon, 16 Jan 2017 16:59:33 -0500
Subject: Multipart Message
MIME-Version: 1.0
Content-Type: multipart/mixed;
 	 boundary="_=test=_bbd1e98aa6c34ef59d8d102a0e795027"
Feedback-ID: 2575071:54973::postmark
To: recipient1@example.com, "Recipient 2" <recipient2@example.com>
Message-ID: <examplemessage@example.com>

--_=test=_bbd1e98aa6c34ef59d8d102a0e795027
Content-Type: multipart/alternative; boundary="_=ALT_=test=_bbd1e98aa6c34ef59d8d102a0e795027"

--_=ALT_=test=_bbd1e98aa6c34ef59d8d102a0e795027
Content-Type: text/plain; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

Sending bees

=F0=9F=90=9D

--_=ALT_=test=_bbd1e98aa6c34ef59d8d102a0e795027
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

<!DOCTYPE html>
<html>
  <body>
    Sending bees<br><br>=F0=9F=90=9D
  </body>
</html>

--_=ALT_=test=_bbd1e98aa6c34ef59d8d102a0e795027--
--_=test=_bbd1e98aa6c34ef59d8d102a0e795027
Content-Type: text/calendar; name="invite.ics"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="invite.ics"

QkVHSU46VkNBTEVOREFSClZFUlNJT046Mi4wClBST0RJRDotLy9tYWlscHJvdG8vL01haWxQcm90bwpDQUxTQ0FMRTpHUkVHT1JJQU4KQkVHSU46VkVWRU5UCkRUU1RBTVA6MjAxNzAxMTZUMTU0MDAwClVJRDpteWNvb2xldmVudEBtYWlscHJvdG8KCkRUU1RBUlQ7VFpJRD0iQW1lcmljYS9OZXdfWW9yayI6MjAxNzAxMThUMTEwMDAwCkRURU5EO1RaSUQ9IkFtZXJpY2EvTmV3X1lvcmsiOjIwMTcwMTE4VDEyMDAwMApTVU1NQVJZOlNlbmQgYW4gZW1haWwKTE9DQVRJT046VGVzdApFTkQ6VkVWRU5UCkVORDpWQ0FMRU5EQVI=
--_=test=_bbd1e98aa6c34ef59d8d102a0e795027--`
)

func TestMessageParsing(t *testing.T) {

	msg, err := smtpd.NewMessage([]byte(multipartEmail), nil)

	if err != nil {
		t.Error("error creating message", err)
		return
	}

	expectTo := []mail.Address{
		{
			Name:    "",
			Address: "recipient1@example.com",
		},
		{
			Name:    "Recipient 2",
			Address: "recipient2@example.com",
		},
	}

	if len(msg.To) < len(expectTo) {
		t.Errorf("Not enough recipients, want: %v, got: %v", len(expectTo), len(msg.To))

	}

	for i, expect := range expectTo {
		if i >= len(msg.To) {
			break
		}
		if msg.To[i].Address != expect.Address || msg.To[i].Name != expect.Name {
			t.Errorf("Wrong recipient %v want: %v, got: %v", i, expect, msg.To[i])
		}
	}

	expectHTML := `<!DOCTYPE html>
<html>
  <body>
    Sending bees<br><br>🐝
  </body>
</html>`

	if html, err := msg.HTML(); err != nil {
		t.Error(err)
	} else if strings.TrimSpace(string(html)) != expectHTML {
		t.Errorf("Wrong HTML content, want: %v, got: %v", expectHTML, strings.TrimSpace(string(html)))
	}

	expectPlain := `Sending bees

🐝`

	if plain, err := msg.Plain(); err != nil {
		t.Error(err)
	} else if strings.TrimSpace(string(plain)) != expectPlain {
		t.Errorf("Wrong Plaintext content, want: %v, got: %v", expectPlain, strings.TrimSpace(string(plain)))
	}

	// TODO: check rest of parse proceeded as expected

}
