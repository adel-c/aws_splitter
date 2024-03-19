package main

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"os"
	"regexp"
)

func runCommand(files FileHandler, regexp *regexp.Regexp) error {

	if isInputFromPipe() {
		printLog("data is from pipe")
		return splitLine(files, os.Stdin, regexp)
	} else {
		file, e := getFile()
		if e != nil {
			return e
		}
		defer file.Close()
		return splitLine(files, file, regexp)
	}
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func getFile() (*os.File, error) {
	if flags.filepath == "" {
		return nil, errors.New("please input a file")
	}
	if !fileExists(flags.filepath) {
		return nil, errors.New("the file provided does not exist")
	}
	file, e := os.Open(flags.filepath)
	if e != nil {
		return nil, errors.Wrapf(e,
			"unable to read the file %s", flags.filepath)
	}
	return file, nil
}

func splitLine(files FileHandler, r io.Reader, regexp *regexp.Regexp) error {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		line := scanner.Text()
		var log = parseLine(regexp, line)
		printLog(log.file + "->" + log.line)
		var outFile = files.getFile(log.file)
		_, err := outFile.WriteString(log.line + "\n")

		if err != nil {
			return err
		}
		errS := outFile.Sync()
		if errS != nil {
			return errS
		}

	}
	return nil
}

func fileExists(filepath string) bool {
	info, e := os.Stat(filepath)
	if os.IsNotExist(e) {
		return false
	}
	return !info.IsDir()
}

type LogLine struct {
	file string
	line string
}

func parseLine(compRegEx *regexp.Regexp, url string) LogLine {

	match := compRegEx.FindStringSubmatch(url)
	var log = LogLine{}

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			if name == "fileName" {
				log.file = match[i]
			}
			if name == "log" {
				log.line = match[i]
			}

		}
	}
	return log
}
