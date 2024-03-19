package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
		var files = FilesList{
			filesMap:       map[string]*os.File{},
			workDir:        flags.outDir,
			truncateOnOpen: flags.clearOpen,
		}

		var compRegEx = regexp.MustCompile(flags.regex)
		err := runCommand(files, compRegEx)
		closeAllFiles(files)
		return err
	},
}

type FilesList struct {
	filesMap       map[string]*os.File
	workDir        string
	truncateOnOpen bool
}

type FileHandler interface {
	getFile(s string) *os.File
}

func closeAllFiles(r FilesList) {
	for s := range r.filesMap {
		err := r.filesMap[s].Close()
		if err != nil {
			println(err)
		}
	}
}
func (r FilesList) getFile(s string) *os.File {
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

	flag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if r.truncateOnOpen {
		flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(path, flag, 0600)
	if err != nil {
		panic(err)
	}
	r.filesMap[s] = f
	return f

}

var flags struct {
	filepath  string
	verbose   bool
	regex     string
	outDir    string
	clearOpen bool
}

var flagsName = struct {
	file, fileShort               string
	verbose, verboseShort         string
	regex, regexShort             string
	outDir, outDirShort           string
	clearOnOpen, clearOnOpenShort string
}{
	"file", "f",
	"verbose", "v",
	"regex", "r",
	"outdir", "o",
	"clear", "c",
}

var printLog func(s string)

func main() {
	rootCmd.Flags().StringVarP(
		&flags.filepath,
		flagsName.file,
		flagsName.fileShort,
		"", "path to the file")
	rootCmd.Flags().StringVarP(
		&flags.regex,
		flagsName.regex,
		flagsName.regexShort,
		"[^ ]+ (?P<fileName>[^ ]+) (?P<log>.*)", "line regex should have two named capture group 'fileName' and 'log'")
	rootCmd.Flags().StringVarP(
		&flags.outDir,
		flagsName.outDir,
		flagsName.outDirShort,
		"./tmp", "out directory")
	rootCmd.PersistentFlags().BoolVarP(
		&flags.verbose,
		flagsName.verbose,
		flagsName.verboseShort,
		false, "log verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&flags.clearOpen,
		flagsName.clearOnOpen,
		flagsName.clearOnOpenShort,
		false, "clear each log file on open")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func logNoop(s string) {}

func logOut(s string) {
	log.Println(s)
}
