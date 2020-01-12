package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/urfave/cli"
)

const Name string = "flac2mp3"
const Version string = "1"

var ffmpegBinaryPath string
var cmdq chan *exec.Cmd
var dest string
var src string

var isScan bool = false

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "source-directory",
		Value: GetCwd(),
		Usage: "<dir>",
	},
	cli.StringFlag{
		Name:  "destination-directory",
		Value: GetCwd() + "/destination",
		Usage: "<dir>",
	},
}

var Commands = []cli.Command{
	{
		Name:        "scan",
		ShortName:   "s",
		Usage:       "scan",
		Description: "Outputs a list of canidate files for processing.",
		Flags:       []cli.Flag{},
		Action: func(c *cli.Context) {
			src := c.GlobalString("source-directory")
			isScan = true
			filepath.Walk(src, traverse)
		},
	},
	{
		Name:        "convert",
		ShortName:   "c",
		Usage:       "convert",
		Description: ".",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "workers",
				Value: runtime.NumCPU(),
				Usage: "<n>",
			},
		},
		Action: func(c *cli.Context) error {

			worker_count := c.Int("workers")

			src = c.GlobalString("source-directory")
			dest = c.GlobalString("destination-directory")

			cmdq = make(chan *exec.Cmd, 1024)

			go func() {
				filepath.Walk(src, traverse)
				close(cmdq)
			}()

			var wg sync.WaitGroup
			for i := 0; i < worker_count; i++ {
				wg.Add(1)
				go func() {
					for cmd := range cmdq {
						_, fierr := os.Stat(cmd.Args[7])

						if fierr != nil {
							fmt.Println(cmd.Args[7])
							cerr := cmd.Run()

							if cerr != nil {
								log.Println(cmd, cerr)
							}
						}
					}
					wg.Done()
				}()
			}

			wg.Wait()
			return nil
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "Adam Flott"
	app.Email = "adam@adamflott.com"
	app.Usage = "parallel conversion of flacs to V0 mp3s with ffmpeg"

	app.Compiled = time.Now()
	app.Copyright = "2020"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
		os.Exit(2)
	}

	ffmpegpath, fferr := exec.LookPath("ffmpeg")
	if fferr != nil {
		log.Fatalln(fferr)
	}

	ffmpegBinaryPath = ffmpegpath

	app.Run(os.Args)
}

func traverse(fpathi string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// Skip all non-regular files (dirs, etc)
	if info.IsDir() == true || info.Mode().IsRegular() == false {
		return nil
	}

	if filepath.Base(fpathi) == ".flac2mp3ignore" {
		fmt.Println("Skipping", filepath.Dir(fpathi))
		return filepath.SkipDir
	}

	if info.Size() == 0 {
		return nil
	}

	if strings.ToLower(path.Ext(fpathi)) != ".flac" {
		return nil
	}

	if isScan {
		fmt.Println(fpathi)
	} else {
		filename := filepath.Base(fpathi)
		ext := filepath.Ext(fpathi)
		filename = strings.TrimSuffix(filename, ext) + ".mp3"

		dir := filepath.Join(dest, strings.TrimPrefix(filepath.Dir(fpathi), src))

		direrr := os.MkdirAll(dir, 0755)

		if direrr != nil {
			log.Fatalln(direrr)
		}

		fpatho := filepath.Clean(filepath.Join(dir, filename))

		ffmpegcmd := exec.Command(ffmpegBinaryPath, "-i", fpathi, "-codec:a", "libmp3lame", "-qscale:a", "0", fpatho)
		cmdq <- ffmpegcmd
	}

	return nil
}

func GetCwd() string {
	dir, derr := os.Getwd()

	if derr != nil {
		log.Fatalln(derr)
	}
	return dir
}
