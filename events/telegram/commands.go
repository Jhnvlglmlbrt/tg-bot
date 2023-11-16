package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/Jhnvlglmlbrt/tg-bot/clients/telegram"
	"github.com/Jhnvlglmlbrt/tg-bot/lib/e"
	"github.com/Jhnvlglmlbrt/tg-bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	ListCmd  = "/list"
)

func (d *Dispatcher) docmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		return d.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		return d.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return d.sendHelp(ctx, chatID)
	case StartCmd:
		return d.sendHello(ctx, chatID)
	case ListCmd:
		return d.sendList(ctx, chatID)
	default:
		return d.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (d *Dispatcher) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(ctx, chatID, d.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := d.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := d.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := d.tg.SendMessage(ctx, chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (d *Dispatcher) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	sendMsg := NewMessageSender(ctx, chatID, d.tg)

	page, err := d.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := d.tg.SendMessage(ctx, chatID, page.URL); err != nil {
		return err
	}

	return d.storage.Remove(ctx, page)
}

func (d *Dispatcher) sendHelp(ctx context.Context, chatID int) error {
	return d.tg.SendMessage(ctx, chatID, msgHelp)
}

func (d *Dispatcher) sendList(ctx context.Context, chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send list", err) }()

	sendMsg := NewMessageSender(ctx, chatID, d.tg)

	urls, err := d.storage.List(ctx)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := d.tg.SendMessage(ctx, chatID, urls); err != nil {
		return err
	}

	return nil
}

func (d *Dispatcher) sendHello(ctx context.Context, chatID int) error {
	return d.tg.SendMessage(ctx, chatID, msgHello)
}

func NewMessageSender(ctx context.Context, chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(ctx, chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
