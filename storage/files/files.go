package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Jhnvlglmlbrt/tg-bot/lib/e"
	"github.com/Jhnvlglmlbrt/tg-bot/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	userPath, err := s.createUserDir(page.UserName)
	if err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath := filepath.Join(userPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	userPath, err := s.createUserDir(userName)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(userPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(userPath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err = os.Remove(path); err != nil {
		msg := fmt.Sprintf("cant remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil

}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)

	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) createUserDir(userName string) (string, error) {
	userPath := filepath.Join(s.basePath, userName)

	_, err := os.Stat(userPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(userPath, defaultPerm); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return userPath, nil
}
