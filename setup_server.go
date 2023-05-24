package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

const requiredVersion = "1.12.6"

func setupServer() error {
	if _, err := os.Stat("httptoolkit-server/package.json"); err != nil {
		if err = downloadServer(); err != nil {
			return err
		}
	}
	return nil
}

func downloadServer() error {
	fmt.Println("Downloading httptoolkit-server " + requiredVersion + " for " + platform)
	fileName := "httptoolkit-server.tar.gz"
	err := dl(
		fmt.Sprintf(
			"https://github.com/juby-httptoolkit/httptoolkit-server/releases/download/v%s/httptoolkit-server-v%s-%s-x64.tar.gz",
			requiredVersion, requiredVersion, platform,
		),
		fileName,
	)
	if err != nil {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	err = extract(file)
	file.Close()
	if err != nil {
		return err
	}
	os.Remove(fileName)
	fmt.Println("Server download completed")
	return nil
}

func dl(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return file.Close()
}

// https://codereview.stackexchange.com/questions/272457/decompress-tar-gz-file-in-go
func extract(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	var header *tar.Header
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("extract: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(header.Name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return fmt.Errorf("extract: Create() failed: %w", err)
			}

			if _, err = io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("extract: Copy() failed: %w", err)
			}
			if err = outFile.Close(); err != nil {
				return fmt.Errorf("extract: Close() failed: %w", err)
			}
		}
	}
	if err != io.EOF {
		return fmt.Errorf("extract: Next() failed: %w", err)
	}
	return nil
}
