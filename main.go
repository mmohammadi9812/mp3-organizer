package main

import (
	"fmt"
	"github.com/dhowden/tag"
	"mp3organize/gui"
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
	re := regexp.MustCompile("(?i)(\\(|\\[)?[^ \t]+\\.(com|ir|net|org)(\\)|])|(@[^ \t])|\"")
	m.title = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.title), ""))
	m.artist = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.artist), ""))
	m.album = strings.TrimSpace(re.ReplaceAllString(strings.TrimSpace(m.album), ""))
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

func moveFile(mp3Path, dest string, format [2]string, info metainfo) {
	if info.title == "" || info.artist == "" {
		fmt.Printf("Skipping %s as it doesn't have title or artist name in it's tags ...\n", filepath.Base(mp3Path))
		return
	}

	var (
		fileDest string
		destDir  string
		err      error
		fileName = info.title + ".mp3"
	)

	if format[0] == "artist" {
		if format[1] == "/" {
			destDir = path.Join(dest, info.artist)
			err = os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			fileDest = path.Join(destDir, fileName)
		} else {
			fileDest = path.Join(dest, fmt.Sprintf("%s - %s", info.artist, fileName))
		}
	} else {
		destDir = path.Join(dest, info.artist)
		if format[1] == "/" {
			destDir = path.Join(destDir, info.album)
			err = os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			fileDest = path.Join(destDir, fileName)
		} else {
			err = os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			fileDest = path.Join(destDir, fmt.Sprintf("%s - %s", info.album, fileName))
		}
	}

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

func singleAlbums(wantSingle bool, meta *metainfo) {
	if wantSingle {
		if meta.album == "" {
			meta.album = "Single"
		} else if matched, err := regexp.MatchString("(?i)single", meta.album); err == nil && matched {
			meta.album = "Single"
		} else if matched, err = regexp.MatchString("(?i)"+meta.title, meta.album); err == nil && matched {
			meta.album = "Single"
		} else if err != nil {
			panic(err)
		} else {
			return
		}
	} else {
		return
	}
}

func main() {
	files := gui.GetFiles()
	saveDir := gui.Dest
	for _, file := range files {
		meta, err := getMeta(file)
		if err != nil {
			panic(err)
		}
		if isEmpty(meta) {
			continue
		}
		singleAlbums(gui.Single, &meta)
		moveFile(file, saveDir, gui.Format, meta)
	}
}
