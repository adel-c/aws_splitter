package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "aws_splitter",
	Short: "split cloudwatch log to separated files",
	Long:  `split cloudwatch log to separated files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printLog = logNoop
		if flags.verbose {
			printLog = logOut
		}
		var files = FilesList{filesMap: map[string]*os.File{},
			workDir: "/home/adel/projects/aws_splitter/work/",
		}
		err := runCommand(files)
		files.Close()
		return err
	},
}

type FilesList struct {
	filesMap map[string]*os.File
	workDir  string
}

type FileHandler interface {
	GetFile(s string) *os.File
	Close()
}

func (r FilesList) Close() {
	for s := range r.filesMap {
		err := r.filesMap[s].Close()
		if err != nil {
			println(err)
		}
	}
}
func (r FilesList) GetFile(s string) *os.File {
	var existingFile, ok = r.filesMap[s]
	// If the key exists
	if ok {
		return existingFile
	}

	var path = filepath.Join(r.workDir, s)
	var dir = filepath.Dir(path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	r.filesMap[s] = f
	return f

}

var flags struct {
	filepath string
	verbose  bool
}

var flagsName = struct {
	file, fileShort       string
	verbose, verboseShort string
}{
	"file", "f",
	"verbose", "v",
}

var printLog func(s string)

func main() {
	rootCmd.Flags().StringVarP(
		&flags.filepath,
		flagsName.file,
		flagsName.fileShort,
		"", "path to the file")
	rootCmd.PersistentFlags().BoolVarP(
		&flags.verbose,
		flagsName.verbose,
		flagsName.verboseShort,
		false, "log verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func logNoop(s string) {}

func logOut(s string) {
	log.Println(s)
}
