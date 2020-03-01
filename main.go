package main


import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)


type metainfo struct {
	title string
	artist string
	album string
}

func (m *metainfo) clean() {
	re := regexp.MustCompile("(?i)(\\(|\\[)\\w+\\.(com|ir|net)(\\)|\\])")
	re2 := regexp.MustCompile("(?i)\\w+\\.(net)\\s+\\|")
	m.title = re.ReplaceAllString(m.title, "")
	m.title = re2.ReplaceAllString(m.title, "")
	m.artist = re.ReplaceAllString(m.artist, "")
	m.artist = re2.ReplaceAllString(m.artist, "")
	m.album = re.ReplaceAllString(m.album, "")
	m.album = re2.ReplaceAllString(m.album, "")
	m.title = strings.TrimSpace(m.title)
	m.artist = strings.TrimSpace(m.artist)
	m.album = strings.TrimSpace(m.album)
}

func getMP3Files(startPath string) (res []string, err error) {
	res, err = filepath.Glob("(?i)" + path.Join(startPath, "*.mp3"))
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func getPath() string {
	var dir string
	if path.IsAbs(os.Args[1]){
		dir = os.Args[1]
	}else{
		var currentDir string
		currentDir = filepath.Dir(os.Args[0])
		dir, _ = filepath.Abs(path.Join(currentDir, os.Args[1]))
	}
	return dir
}

func getMeta(filePath string) (metainfo, error) {
	f, ferr := os.Open(filePath)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer func(){
		_ = f.Close()
	}()
	meta, err := tag.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
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

func moveFile(mp3Path, saveDir string, info metainfo) {
	newDir := path.Join(saveDir, info.artist)
	err := os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fileName := info.title + ".mp3"
	newPath := path.Join(newDir, fileName)
	err = os.Rename(mp3Path, newPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s was moved to %s", mp3Path, newPath)
}

func moveFiles(mp3Dir, saveDir string) {
	files, ferr := getMP3Files(mp3Dir)
	if ferr != nil {
		log.Fatal(ferr)
	}
	for _, file := range files {
		meta, err := getMeta(file)
		if err != nil {
			log.Fatal(err)
		}
		moveFile(file, saveDir, meta)
	}
}

func main(){
	Home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	dir := getPath()
	var saveDir string
	if len(os.Args) == 2 {
		saveDir = os.Args[2]
	} else {
		saveDir = path.Join(Home, "Music")
	}
	moveFiles(dir, saveDir)
}
