package lib

// Extracted from https://github.com/mauri870/gcsfs

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

type File struct {
	reader io.ReadCloser
	writer io.WriteCloser
	attrs  *storage.ObjectAttrs
}

type fileInfo struct {
	dirModTime time.Time
	attrs      *storage.ObjectAttrs
}

func (f *fileInfo) Name() string {
	name := f.attrs.Name
	if f.IsDir() {
		name = f.attrs.Prefix
	}
	return filepath.Base(name)
}

func (f *fileInfo) Info() (fs.FileInfo, error) {
	return f, nil
}

func (f *fileInfo) Size() int64 {
	return f.attrs.Size
}

func (f *fileInfo) Mode() fs.FileMode {
	if f.IsDir() {
		return fs.ModeDir
	}

	return 0
}

func (f *fileInfo) ModTime() time.Time {
	if f.IsDir() {
		return f.dirModTime
	}
	return f.attrs.Updated
}

func (f *fileInfo) IsDir() bool {
	return f.attrs.Prefix != ""
}

func (f *fileInfo) Sys() interface{} {
	return nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return &fileInfo{attrs: f.attrs}, nil
}

func (f *File) Read(p []byte) (int, error) {
	return f.reader.Read(p)
}

func (f *File) Close() error {
	return f.reader.Close()
}

type GCSFS struct {
	bucket string
}

func (gcsfs GCSFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return gcsfs.getFile(name)
}

func (gcsfs *GCSFS) getFile(name string) (*File, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(gcsfs.bucket)
	obj := bucket.Object(name)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	w := obj.NewWriter(ctx)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return &File{reader: r, writer: w, attrs: attrs}, nil
}

func GetProviderByEnv() (fs.FS, error) {
	switch os.Getenv("env") {
	case "local":
		return os.DirFS("../function-data/dev/"), nil
	case "dev":
		return GCSFS{bucket: "function-data"}, nil
	case "prod":
		return GCSFS{bucket: "core-350507-function-data"}, nil
	default:
		return nil, fmt.Errorf("unhandled environment")
	}
}
