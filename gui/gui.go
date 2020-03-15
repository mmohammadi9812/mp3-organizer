package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"path/filepath"
)

type check struct {
	label   string
	checked bool
}

func (c *check) toggle(on bool) {
	c.checked = on
}

func CheckFiles(files []string) (checkedFiles []string) {
	icon, err := fyne.LoadResourceFromPath("C:\\Users\\Mohammad\\Documents\\Code\\GO\\mp3Organizer\\gui\\icon.png")
	if err != nil {
		panic(err)
	}
	a := app.New()
	a.SetIcon(icon)
	w := a.NewWindow("MP3 Organizor")

	var (
		//checkBoxes []fyne.CanvasObject
		storeChecks []*check
		confirmed   = false
	)

	var (
		quitButton = widget.NewButton("Quit", func() {
			a.Quit()
		})

		okButton = widget.NewButton("Ok", func() {
			confirmed = true
			a.Quit()
		})
	)

	var (
		filesBox   = widget.NewGroupWithScroller("Musics")
		buttonsBox = fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2), quitButton, okButton)
	)

	for _, file := range files {
		var fileCheck = check{
			checked: false,
			label:   file,
		}
		storeChecks = append(storeChecks, &fileCheck)
		filesBox.Append(widget.NewCheck(filepath.Base(fileCheck.label), fileCheck.toggle))
	}

	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, buttonsBox, nil, nil),
			filesBox,
			buttonsBox,
		),
	)

	w.Resize(fyne.Size{
		Width:  320,
		Height: 480,
	})
	w.ShowAndRun()

	for _, checkBox := range storeChecks {
		if checkBox.checked {
			checkedFiles = append(checkedFiles, checkBox.label)
		}
	}

	return checkedFiles
}
