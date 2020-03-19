package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/mitchellh/go-homedir"
	"github.com/sqweek/dialog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return home
}

var (
	confirmed     = false
	myapp         fyne.App
	mywindow      fyne.Window
	filesCheckBox []*widget.Check
	files         []string
	filesBox      *widget.Group
	labelsArr     []*check
	home          = getHome()
	src           string
	singleAlbums  *widget.Check
	opts          = make(map[string]bool, 0)
	recursive     = true
	Dest          = path.Join(home, "Music")
	Format        [2]string
	Single        = false
)

type check struct {
	label   string
	checked bool
}

func (c *check) toggle(on bool) {
	c.checked = on
}

func getQOKButtonsContainer() *fyne.Container {
	var (
		quitButton = widget.NewButton("Quit", func() {
			myapp.Quit()
		})

		okButton = widget.NewButton("Ok", func() {
			confirmed = true
			myapp.Quit()
		})
		buttonsContainer = fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2), quitButton, okButton)
	)
	return buttonsContainer
}

func getMP3Files(src string) (mp3Files []string, err error) {
	mp3Files, err = filepath.Glob(path.Join(src, "*.[mM][pP]3"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d mp3 files in %s ...\n", len(mp3Files), src)
	return mp3Files, nil
}

func getMP3FilesRecursive(src string) (mp3Files []string, err error) {
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if matched, reErr := regexp.MatchString("(?i)mp3", filepath.Ext(path)); reErr == nil && matched {
			mp3Files = append(mp3Files, path)
		}
		return nil
	})
	return mp3Files, err
}

func addFiles() {
	for _, file := range files {
		var newLabel = check{
			checked: true,
			label:   file,
		}

		labelsArr = append(labelsArr, &newLabel)
		checkBox := widget.NewCheck(filepath.Base(newLabel.label), newLabel.toggle)
		checkBox.SetChecked(true)
		filesCheckBox = append(filesCheckBox, checkBox)
		filesBox.Append(checkBox)
	}

	filesBox.Refresh()
}

func getSrcButton(label string) *widget.Button {
	var srcButton *widget.Button
	srcButton = widget.NewButton(label, func() {
		if len(filesCheckBox) != 0 {
			for _, fileCheckBox := range filesCheckBox {
				fileCheckBox.Disable()
			}
			filesBox.Refresh()
		}

		filesCheckBox = []*widget.Check{}
		labelsArr = []*check{}

		directoryStruct := &dialog.DirectoryBuilder{
			Dlg: dialog.Dlg{
				Title: "Select source directory",
			},
			StartDir: os.Getenv("PWD"),
		}
		directory, err := directoryStruct.Browse()
		if err != nil {
			//panic(err)
			fmt.Print(err)
		}

		if recursive {
			files, err = getMP3FilesRecursive(directory)
		} else {
			files, err = getMP3Files(directory)
		}

		if err != nil {
			panic(err)
		}

		srcButton.SetText(fmt.Sprintf("Source: %s", filepath.Base(directory)))

		src = directory
		addFiles()
	})
	return srcButton
}

func getDestButton(labelFormatStr string) *widget.Button {
	var (
		destButton *widget.Button
		label      = fmt.Sprintf("%s: %s", labelFormatStr, filepath.Base(Dest))
	)
	destButton = widget.NewButton(label, func() {
		directory, err := dialog.Directory().Title("Select destination folder").Browse()
		if err != nil {
			//panic(err)
			fmt.Print(err)
		}
		destButton.SetText(fmt.Sprintf("Destination: %s", filepath.Base(directory)))
		Dest = directory
	})
	return destButton
}

func getFormat() *fyne.Container {
	var albumOrNotOpts = []string{
		"artist",
		"artist/album",
	}
	var albumOrNot = widget.NewSelect(albumOrNotOpts, func(selected string) {
		Format[0] = selected
		if singleAlbums == nil {
			return
		} else if selected == "artist" {
			singleAlbums.Hide()
		} else {
			singleAlbums.Show()
		}
	})
	var fileOrDirOpts = []string{
		"/ (make directory)",
		"- (just rename files)",
	}
	var fileOrDir = widget.NewSelect(fileOrDirOpts, func(selected string) {
		Format[1] = strings.Fields(selected)[0]
	})

	albumOrNot.SetSelected("artist")
	fileOrDir.SetSelected("/ (make directory)")

	var formatLayout = fyne.NewContainerWithLayout(layout.NewVBoxLayout(), albumOrNot, fileOrDir)
	return formatLayout
}

func getOptsSlice() []*widget.Check {
	var (
		recursiveButton = widget.NewCheck("Recursive file select", func(on bool) {
			opts["recursive"] = on
			handleOpt("recursive", on)
		})
		compressSingleButton = widget.NewCheck("Move single tracks to one folder", func(on bool) {
			opts["single"] = on
			handleOpt("single", on)
		})
		selectAllButton = widget.NewCheck("Select All", func(on bool) {
			opts["selectall"] = on
			handleOpt("selectall", on)
		})
	)
	recursiveButton.SetChecked(true)
	compressSingleButton.SetChecked(false)
	compressSingleButton.Hide()

	singleAlbums = compressSingleButton

	selectAllButton.SetChecked(true)
	return []*widget.Check{
		selectAllButton,
		recursiveButton,
		compressSingleButton,
	}
}

func handleOpt(option string, value bool) {
	var err error
	switch option {
	case "selectall":
		for _, checkBox := range filesCheckBox {
			checkBox.SetChecked(value)
		}
		break
	case "single":
		Single = value
		break
	case "recursive":
		recursive = value
		if len(filesCheckBox) != 0 {
			for _, fileCheckBox := range filesCheckBox {
				fileCheckBox.Disable()
			}
			filesBox.Refresh()
		}
		filesCheckBox = []*widget.Check{}
		labelsArr = []*check{}
		if value {
			files, err = getMP3FilesRecursive(src)
		} else {
			files, err = getMP3Files(src)
		}
		if err != nil {
			panic(err)
		}
		addFiles()
		break
	}
}

func GetFiles() (checkedFiles []string) {
	//icon, err := fyne.LoadResourceFromPath("gui\\icon.png")
	//if err != nil {
	//	panic(err)
	//}
	myapp = app.New()
	//a.SetIcon(icon)
	mywindow = myapp.NewWindow("MP3 Organizor")
	filesBox = widget.NewGroupWithScroller("Musics")
	var (
		buttonsContainer = getQOKButtonsContainer()
		srcButton        = getSrcButton("Source Directory")
		destButton       = getDestButton("Default Destination")
		formatSelector   = getFormat()
		pathButtonLayout = fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(3), srcButton, destButton, formatSelector)
		otherOptions     = getOptsSlice()
	)

	for _, option := range otherOptions {
		filesBox.Prepend(option)
	}

	mywindow.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(pathButtonLayout, buttonsContainer, nil, nil),
		pathButtonLayout,
		filesBox,
		buttonsContainer,
	),
	)

	mywindow.Resize(fyne.Size{
		Width:  350,
		Height: 480,
	})
	mywindow.ShowAndRun()

	for _, label := range labelsArr {
		if label.checked {
			checkedFiles = append(checkedFiles, label.label)
		}
	}

	if confirmed {
		return checkedFiles
	} else {
		return []string{}
	}

}
