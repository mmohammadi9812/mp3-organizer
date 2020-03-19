package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"path/filepath"
)

var (
	confirmed     = false
	myapp         fyne.App
	mywindow      fyne.Window
	filesCheckBox []*widget.Check
)

type check struct {
	label   string
	checked bool
}

func (c *check) toggle(on bool) {
	c.checked = on
}

func getButtonsContainer() *fyne.Container {
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

func getSelectAllButton() *widget.Check {
	selectAllButton := widget.NewCheck("Select All", func(on bool) {
		for _, checkBox := range filesCheckBox {
			checkBox.SetChecked(on)
		}
	})
	return selectAllButton
}

func getFilesBox(files []string, title string) (filesBox *widget.Group, labelsArr []*check) {
	filesBox = widget.NewGroupWithScroller(title)

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

	return filesBox, labelsArr
}

func CheckFiles(files []string) (checkedFiles []string) {
	//icon, err := fyne.LoadResourceFromPath("gui\\icon.png")
	//if err != nil {
	//	panic(err)
	//}
	myapp = app.New()
	//a.SetIcon(icon)
	mywindow = myapp.NewWindow("MP3 Organizor")

	var (
		filesBox, labelsArr = getFilesBox(files, "Musics")
		buttonsContainer    = getButtonsContainer()
		selectAll           = getSelectAllButton()
	)

	selectAll.SetChecked(true)

	filesBox.Append(selectAll)

	mywindow.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, buttonsContainer, nil, nil),
			filesBox,
			buttonsContainer,
		),
	)

	mywindow.Resize(fyne.Size{
		Width:  320,
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
