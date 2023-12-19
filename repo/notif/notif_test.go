package notif

import (
	"testing"
	"os"
	"strings"
)

const (
	TEST_VAR_APP_NAME = "APP_NAME"
	TEST_VAR_APP_PWD = "APP_PWD"
	TEST_VAR_NOTIF_HOST = "NOTIF_HOST"
	TEST_VAR_TM_CHAT_ID = "TM_CHAT_ID"
	TEST_VAR_MSG_TEXT = "MSG_TEXT"
	TEST_VAR_MSG_TEL = "MSG_TEL"
	
	TEST_VAR_MAIL_SENDER_ADDR	= "MAIL_SENDER_ADDR"
	TEST_VAR_MAIL_SENDER_USER_NAME	= "MAIL_SENDER_USER_NAME"
	TEST_VAR_MAIL_HOST		= "MAIL_HOST"
	TEST_VAR_MAIL_PWD		= "MAIL_PWD"
	
	TEST_VAR_MAIL_TO_ADDR		="MAIL_TO_ADDR"
	TEST_VAR_MAIL_TO_NAME		="MAIL_TO_NAME"
	TEST_VAR_MAIL_SUBJECT		="MAIL_SUBJECT"
	TEST_VAR_MAIL_BODY		="MAIL_BODY"
	TEST_VAR_MSG_TYPE		="MSG_TYPE"
	TEST_VAR_MAIL_ATTACHMENTS	="MAIL_ATTACHMENTS"
	TEST_VAR_MAIL_ATTACHMENT_ALIAS	="MAIL_ATTACHMENT_ALIAS"
	
)

func getTestVar(t *testing.T, n string) string {
	v := os.Getenv(n)
	if v == "" {
		t.Fatalf("getTestVar() failed: %s environment variable is not set", n)
	}
	return v
}

type Messager interface{
	NewNotif(string) *NotifMessage
}
func ProviderSend(t *testing.T, msg Messager) {
	n := NewNotifier(getTestVar(t, TEST_VAR_NOTIF_HOST), getTestVar(t, TEST_VAR_APP_NAME), getTestVar(t, TEST_VAR_APP_PWD))
	batch := []*NotifMessage{msg.NewNotif(getTestVar(t, TEST_VAR_MSG_TYPE))}
	resp_ar, err := n.Send(batch)
	if err != nil {
		t.Fatalf("n.Send() failed: %v", err)
	}
	for i,resp := range resp_ar {
		if resp.Error != "" {
			t.Fatalf("n.Send() message with index %d failed: %v", i, err)
		}
		t.Logf("n.Send() index %d, ID: %d", i, resp.ID)
	}
}

func TestSendTM(t *testing.T) {	
	msg := &TMMessage{ChatID: getTestVar(t, TEST_VAR_TM_CHAT_ID),
		Text: getTestVar(t, TEST_VAR_MSG_TEXT),
	}
	ProviderSend(t, msg)
}

func TestSendWA(t *testing.T) {
	msg := &WAMessage{Tel: getTestVar(t, TEST_VAR_MSG_TEL),
		Text: getTestVar(t, TEST_VAR_MSG_TEXT),
	}
	ProviderSend(t, msg)
}
func TestSendSMS(t *testing.T) {
	msg := &SMSMessage{Tel: getTestVar(t, TEST_VAR_MSG_TEL),
		Text: getTestVar(t, TEST_VAR_MSG_TEXT),
	}
	ProviderSend(t, msg)
}

func TestSendEmail(t *testing.T) {
	var attachments []string
	msg := &EmailMessage{FromAddr: getTestVar(t, TEST_VAR_MAIL_SENDER_ADDR),
		FromName: getTestVar(t, TEST_VAR_MAIL_SENDER_USER_NAME),
		ToAddr: getTestVar(t, TEST_VAR_MAIL_TO_ADDR),
		ToName: getTestVar(t, TEST_VAR_MAIL_TO_NAME),
		ReplyName: getTestVar(t, TEST_VAR_MAIL_SENDER_USER_NAME),
		Body: getTestVar(t, TEST_VAR_MAIL_BODY),
		SenderAddr: getTestVar(t, TEST_VAR_MAIL_SENDER_ADDR),
		Subject: getTestVar(t, TEST_VAR_MAIL_SUBJECT),
		Attachments: attachments,
	}	
	ProviderSend(t, msg)
}
func TestSendEmailAtt(t *testing.T) {
	var attachments []string
	var attachment_alias []string
	att_s := getTestVar(t, TEST_VAR_MAIL_ATTACHMENTS)
	if att_s != "" {
		attachments = strings.Split(att_s, " ")
	}
	att_al_s := getTestVar(t, TEST_VAR_MAIL_ATTACHMENT_ALIAS)
	if att_al_s != "" {
		attachment_alias = strings.Split(att_al_s, " ")
	}
	
	msg := &EmailMessage{FromAddr: getTestVar(t, TEST_VAR_MAIL_SENDER_ADDR),
		FromName: getTestVar(t, TEST_VAR_MAIL_SENDER_USER_NAME),
		ToAddr: getTestVar(t, TEST_VAR_MAIL_TO_ADDR),
		ToName: getTestVar(t, TEST_VAR_MAIL_TO_NAME),
		ReplyName: getTestVar(t, TEST_VAR_MAIL_SENDER_USER_NAME),
		Body: getTestVar(t, TEST_VAR_MAIL_BODY),
		SenderAddr: getTestVar(t, TEST_VAR_MAIL_SENDER_ADDR),
		Subject: getTestVar(t, TEST_VAR_MAIL_SUBJECT),
		Attachments: attachments,
		AttachmentAlias: attachment_alias,
	}	
	ProviderSend(t, msg)
}
