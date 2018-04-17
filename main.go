package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/urfave/cli"
)

var humanReadable bool

func main() {

	app := cli.NewApp()

	app.Name = "gs"
	app.Usage = "List files in a directory (An ls clone in go.)"
	app.HideVersion = true

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "display this help and exit",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "human-readable, h",
			Usage:       "with -l and -s, print sizes like 1K 234M 2G etc.",
			Destination: &humanReadable,
		},
	}

	app.Action = func(c *cli.Context) error {
		directory := "."

		if c.NArg() > 0 {
			directory = c.Args()[0]
		}

		return ShowOutput(directory)
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// ShowOutput runs the main function of the program to list files
func ShowOutput(directory string) error {

	path := filepath.Dir(directory)

	files, err := OSReadDir(path)

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	for _, fileInfo := range files {

		size, err := determineSize(fileInfo)

		if err != nil {
			log.Fatal(err)
		}

		time := determineDate(fileInfo.ModTime())

		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", fileInfo.Mode().String(), size, time, fileInfo.Name())
	}

	w.Flush()

	return nil
}

func determineSize(fi os.FileInfo) (string, error) {
	var size = fi.Size()

	if fi.IsDir() {
		return "-", nil
	}

	if humanReadable {
		var v datasize.ByteSize
		err := v.UnmarshalText([]byte(strconv.FormatInt(size, 10) + "b"))
		return v.HumanReadable(), err
	}

	return strconv.FormatInt(size, 10), nil
}

func determineDate(t time.Time) string {
	if t.Year() != time.Now().Year() {
		return t.Format("01 Jan 2006")
	}

	return t.Format("01 Jan 15:04")

}

// OSReadDir reads a directory and returns all info objects
func OSReadDir(root string) ([]os.FileInfo, error) {
	var files []os.FileInfo

	f, err := os.Open(root)
	defer f.Close()

	if err != nil {
		return files, err
	}

	// the -1 means it will read all fileinfo
	fileInfo, err := f.Readdir(-1)

	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file)
	}

	return files, nil
}
