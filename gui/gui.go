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
	//icon, err := fyne.LoadResourceFromPath("gui\\icon.png")
	//if err != nil {
	//	panic(err)
	//}
	a := app.New()
	//a.SetIcon(icon)
	w := a.NewWindow("MP3 Organizor")

	var (
		confirmed  = false
		quitButton = widget.NewButton("Quit", func() {
			a.Quit()
		})

		okButton = widget.NewButton("Ok", func() {
			confirmed = true
			a.Quit()
		})
		filesBox      = widget.NewGroupWithScroller("Musics")
		buttonsBox    = fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2), quitButton, okButton)
		labelsArr     []*check
		checkBoxesArr []*widget.Check
		selectAll     = widget.NewCheck("Select All", func(on bool) {
			for _, checkBox := range checkBoxesArr {
				checkBox.SetChecked(on)
			}
		})
	)

	selectAll.SetChecked(true)

	filesBox.Append(selectAll)

	for _, file := range files {
		var newLabel = check{
			checked: true,
			label:   file,
		}
		labelsArr = append(labelsArr, &newLabel)
		checkBox := widget.NewCheck(filepath.Base(newLabel.label), newLabel.toggle)
		checkBox.SetChecked(true)
		checkBoxesArr = append(checkBoxesArr, checkBox)
		filesBox.Append(checkBox)
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
