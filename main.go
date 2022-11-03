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

func init() {
	flag.StringVar(&fileFlag, "f", "out.tar", "tar file name")
	flag.BoolVar(&extractFlag, "x", false, "extract tar")

	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		println("files required")
		os.Exit(1)
	}

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
		if d.IsDir() || path == tarAbs {
			return nil
		}

		relPath, err := filepath.Rel(curAbs, path)
		if err != nil {
			return err
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

		finfo, err := d.Info()
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
	//todo realuze
	return nil
}
