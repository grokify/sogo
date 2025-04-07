package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/emersion/go-imap"
	"github.com/grokify/mogo/config"
	"github.com/grokify/sogo/net/imaputil"
)

func main() {
	_, err := config.LoadDotEnv([]string{".env"}, 1)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	cm, err := imaputil.NewClientMoreEnv(imaputil.DefaultEnvPrefix)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	err = cm.ConnectAndLogin()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(3)
	}

	// defer cm.Logout() // lint: triggers: `Error return value of `cm.Logout` is not checked (errcheck)`

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- cm.Client.List("", "*", mailboxes)
	}()

	slog.Info("Mailboxes:")
	for m := range mailboxes {
		slog.Info("* " + m.Name)
	}

	err = cm.Logout()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(4)
	}

	fmt.Println("DONE")
	os.Exit(0)
}
