package main

import (
	"archive/tar"
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var fileFlag string
var extractFlag bool
var pathFlag string

func init() {
	flag.StringVar(&fileFlag, "f", "out.tar", "tar file name")
	flag.BoolVar(&extractFlag, "x", false, "extract tar")
	flag.StringVar(&pathFlag, "p", ".", "extract path")

	flag.Parse()
}

func main() {
	if extractFlag {
		if err := toUntar(); err != nil {
			println("failed to extract tar file")
			os.Exit(1)
		}
	} else {
		if err := toTar(); err != nil {
			println("failed to create tar")
			os.Exit(1)
		}
	}
	println("done")
}

func toTar() error {
	if len(flag.Args()) == 0 {
		println("files required")
		os.Exit(1)
	}

	tarAbs, err := filepath.Abs(fileFlag)
	if err != nil {
		return err
	}

	tarFile, err := os.Create(tarAbs)
	if err != nil {
		return err
	}

	defer func(tarFile *os.File) {
		if err := tarFile.Close(); err != nil {
			println("failed to close tar file")
		}
	}(tarFile)

	tw := tar.NewWriter(tarFile)
	defer func(tw *tar.Writer) {
		if err := tw.Close(); err != nil {
			println("failed to close tar writer")
		}
	}(tw)

	curAbs, err := filepath.Abs("")
	if err != nil {
		return err
	}

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if path == tarAbs {
			return nil
		}

		finfo, err := d.Info()
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(curAbs, path)
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(finfo, relPath)
		if err != nil {
			return err
		}

		hdr.Name = relPath

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func(src *os.File) {
			if err := src.Close(); err != nil {
				println("failed to close file " + src.Name())
			}
		}(src)

		if _, err := io.Copy(tw, src); err != nil {
			return err
		}

		return nil
	}

	for _, path := range flag.Args() {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if err := filepath.WalkDir(absPath, walkFunc); err != nil {
			return err
		}
	}

	return nil
}

func toUntar() error {
	absExPath, err := filepath.Abs(pathFlag)
	if err != nil {
		return err
	}

	absTarPath, err := filepath.Abs(fileFlag)
	if err != nil {
		return err
	}

	tarFile, err := os.Open(absTarPath)
	if err != nil {
		return err
	}

	defer func(tarFile *os.File) {
		if err := tarFile.Close(); err != nil {
			println("failed to close tar file")
		}
	}(tarFile)

	tr := tar.NewReader(tarFile)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		absTgtPath := filepath.Join(absExPath, hdr.Name)

		if hdr.FileInfo().IsDir() {
			if err := os.MkdirAll(absTgtPath, hdr.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		tgtFile, err := os.OpenFile(absTgtPath, os.O_CREATE|os.O_WRONLY, hdr.FileInfo().Mode())

		if _, err := io.Copy(tgtFile, tr); err != nil {
			return err
		}

		if err := tgtFile.Close(); err != nil {
			return err
		}
	}
	return nil
}
