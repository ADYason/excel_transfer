package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) != 4 {
		log.Print("Not enough args")
		log.Print("Need folder to search, field to search, output folder")
		return
	}
	log.Print("Start working", "\n")
	search_dir := os.Args[1]
	search_field := os.Args[2]
	output_dir := os.Args[3]
	sheet_to_search := parseSheetArg()
	files, err := os.ReadDir(search_dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		src := path.Join(search_dir, file.Name())
		dest := path.Join(output_dir, file.Name())
		f, err := excelize.OpenFile(src)
		if err != nil {
			log.Printf("%s not an excel file", file.Name())
			continue
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()
		rows, err := f.GetRows(sheet_to_search)
		if err != nil {
			log.Println(err)
			return
		}
		var found_flag bool
		found_flag = false
		for _, row := range rows {
			for _, colCell := range row {
				if colCell == search_field {
					log.Printf("%s found in %s", search_field, file.Name())
					log.Printf("Transfering %s to %s", file.Name(), output_dir)
					copyFileContents(src, dest)
					found_flag = true
					break
				}
			}
		}
		if !found_flag {
			log.Printf("'%s' not found in %s file\n", search_field, file.Name())
		}
	}
	log.Print("End working")

}

func parseSheetArg() string {
	if len(os.Args) == 5 {
		return os.Args[4]
	} else {
		return "Sheet1"
	}
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
