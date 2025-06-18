package compress

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GzipCompressor implementing Compressor interface.
type GzipCompressor struct{}

func createGZipArchive(ctx context.Context, files []string, buf io.Writer) error {
	var err error

	gw := gzip.NewWriter(buf)
	defer func() {
		if cerr := gw.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("closing gzip writer: %w", err)
			}
		}
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		if terr := tw.Close(); terr != nil {
			if err == nil {
				err = fmt.Errorf("closing tar writer: %w", err)
			}
		}
	}()

	for i := range files {
		err := appendToArchive(ctx, tw, files[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func appendToArchive(_ context.Context, tw *tar.Writer, filename string) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			if err == nil {
				err = ferr
			}
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// https://golang.org/src/archive/tar/common.go?#L626
	// If we want to preserve the directory structure
	// header.Name = filename
	// else files only then
	header.Name = info.Name()

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

// Compress tages an archive and an array of file paths to compress into archive
func (gz *GzipCompressor) Compress(ctx context.Context, archive string, files ...string) error {
	out, err := os.Create(filepath.Clean(archive + ".tar.gz"))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	err = createGZipArchive(ctx, files, out)
	if err != nil {
		return err
	}

	return nil
}
