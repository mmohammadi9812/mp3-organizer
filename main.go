package main

import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type metainfo struct {
	title  string
	artist string
	album  string
}

func (m *metainfo) clean() {
	re := regexp.MustCompile("(?i)(\\(|\\[)?[^ \t]+\\.(com|ir|net)(\\)|])|(@[^ \t])|\"")
	m.title = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.title), ""))
	m.artist = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.artist), ""))
	m.album = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.album), ""))
}

func showHelp() {
	fmt.Printf("Usage:\n%s [source music directory] (optional)[destination]", os.Args[0])
}

func validateRaise(src, dest string) {
	var (
		windowsAbsRegex = regexp.MustCompile("[a-zA-Z]:\\\\.+")
		cwd             = os.Getenv("PWD")
	)
	if windowsAbsRegex.MatchString(src) || path.IsAbs(src) {
		if fs, err := os.Stat(src); os.IsNotExist(err) || !fs.IsDir() {
			panic(err)
		}
	} else {
		absSrc := path.Join(cwd, src)
		if fs, err := os.Stat(absSrc); os.IsNotExist(err) || !fs.IsDir() {
			panic(err)
		}
	}

	if windowsAbsRegex.MatchString(dest) || path.IsAbs(dest) {
		if fs, err := os.Stat(dest); os.IsExist(err) && !fs.IsDir() {
			panic("Destination is not a folder.")
		} else if os.IsNotExist(err) {
			err = os.MkdirAll(dest, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	} else {
		absDest := path.Join(cwd, dest)
		if fs, err := os.Stat(dest); os.IsExist(err) && !fs.IsDir() {
			panic("Destination is not a folder.")
		} else if os.IsNotExist(err) {
			err = os.MkdirAll(absDest, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getMP3Files(src string) (mp3Files []string, err error) {
	mp3Files, err = filepath.Glob(path.Join(src, "*.[mM][pP]3"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d mp3 files in %s ...\n", len(mp3Files), src)
	return mp3Files, nil
}

func getMeta(filePath string) (metainfo, error) {
	file, err2 := os.Open(filePath)
	if err2 != nil {
		panic(err2)
	}
	defer func() {
		_ = file.Close()
	}()
	meta, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("No tags found for %s\n", filepath.Base(filePath))
		return metainfo{}, err
	}
	title := meta.Title()
	artist := meta.Artist()
	album := meta.Album()
	if album == "" {
		album = title + " Single"
	}
	result := &metainfo{title, artist, album}
	result.clean()
	return *result, nil
}

func moveFile(mp3Path, dest string, info metainfo) {
	if info.title == "" || info.artist == "" {
		fmt.Printf("Skipping %s as it doesn't have title or artist name in it's tags ...\n", filepath.Base(mp3Path))
		return
	}
	destDir := path.Join(dest, info.artist)
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fileName := info.title + ".mp3"
	fileDest := path.Join(destDir, fileName)
	err = os.Rename(mp3Path, fileDest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s was moved to %s\n", filepath.Base(mp3Path), filepath.Base(destDir))
}

func isEmpty(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	//Then see if it's a struct
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// and create an empty copy of the struct object to compare against
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}
func moveFiles(mp3Dir, saveDir string) {
	files, err := getMP3Files(mp3Dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		meta, err := getMeta(file)
		if isEmpty(meta) {
			continue
		}
		if err != nil {
			panic(err)
		}
		moveFile(file, saveDir, meta)
	}
}

func main() {
	var (
		src  string
		dest string
		home string
		err  error
	)
	home, err = homedir.Dir()
	if err != nil {
		panic(err)
	}
	switch len(os.Args) {
	case 2:
		src = os.Args[1]
		dest = path.Join(home, "Music")
	case 3:
		src = os.Args[1]
		dest = os.Args[2]
	default:
		showHelp()
		os.Exit(0)
	}
	validateRaise(src, dest)
	fmt.Printf("Moving files from %s to %s\n", filepath.Base(src), filepath.Base(dest))
	moveFiles(src, dest)
}
