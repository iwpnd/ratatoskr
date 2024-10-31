package compress

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
)

// GzipCompressor ...
type GzipCompressor struct{}

func createGZipArchive(ctx context.Context, files []string, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for i := range files {
		err := appendToArchive(ctx, tw, files[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func appendToArchive(ctx context.Context, tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

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
	out, err := os.Create(archive + ".tar.gz")
	if err != nil {
		return err
	}
	defer out.Close()

	err = createGZipArchive(ctx, files, out)
	if err != nil {
		return err
	}

	return nil
}
