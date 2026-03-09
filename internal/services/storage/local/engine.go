package local

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/vayload/vayload/internal/services/storage"
)

type Store struct {
	root string
}

func NewManager(root string) *Store {
	return &Store{root: root}
}

func (s *Store) Put(ctx context.Context, bucket string, key string, reader io.Reader) (storage.PutResult, error) {
	path := objectPath(s.root, bucket, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return storage.PutResult{}, err
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return storage.PutResult{}, err
	}
	defer f.Close()

	h := sha256.New()
	w := io.MultiWriter(f, h)

	n, err := io.Copy(w, reader)
	if err != nil {
		return storage.PutResult{}, err
	}

	sum := h.Sum(nil)
	var sha [32]byte
	copy(sha[:], sum)

	if err := os.Rename(tmp, path); err != nil {
		return storage.PutResult{}, err
	}

	return storage.PutResult{
		Size:   n,
		SHA256: sha,
	}, nil
}

func (s *Store) Get(ctx context.Context, bucket string, key string) (io.ReadCloser, storage.GetInfo, error) {
	path := objectPath(s.root, bucket, key)

	info, err := os.Stat(path)
	if err != nil {
		return nil, storage.GetInfo{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, storage.GetInfo{}, err
	}

	return f, storage.GetInfo{
		Size: info.Size(),
	}, nil
}

func (s *Store) Delete(ctx context.Context, bucket string, key string) error {
	path := objectPath(s.root, bucket, key)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (s *Store) Exists(ctx context.Context, bucket string, key string) (bool, error) {
	path := objectPath(s.root, bucket, key)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func objectPath(root, bucket, key string) string {
	h := hashKey(key)
	return filepath.Join(root, "objects", bucket, h[:2], h+".data")
}
