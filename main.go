package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	log.Print("Start working", "\n")
	file_to_pars, out_dir := parseArgs()
	log.Print("Find file to parse and out dirrectory")
	folders_to_copy, _ := getFolders(file_to_pars)
	log.Print("Got folders to copy")
	for _, folder := range folders_to_copy {
		cpDir(folder, path.Join(out_dir, path.Base(folder)))
	}
	log.Print("End working")

}

func getFolders(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		panic("Not enough args. Need file to parse (.txt), folder to dump.")
	}
	file_to_pars := os.Args[1]
	out_dir = path.Join("..", os.Args[2])
	return file_to_pars, out_dir
}

func cpDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = cpDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = cpFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func cpFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
