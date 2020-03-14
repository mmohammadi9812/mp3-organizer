package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"path/filepath"
)

type check struct {
	label string
	checked bool
}

func (c* check) toggle (on bool) {
	c.checked = on
}

func CheckFiles(files []string) (checkedFiles []string) {
	//icon, err := fyne.LoadResourceFromPath("icon.png")
	//if err != nil {
	//	panic(err)
	//}
	a := app.New()
	//a.SetIcon(icon)
	w := a.NewWindow("MP3 Organizor")

	var checkBoxes []fyne.CanvasObject = []fyne.CanvasObject{widget.NewButton("Quit", func() {
		a.Quit()
	})}

	var storeChecks []*check

	for _, file := range files {
		var fileCheck = check{checked:false, label:filepath.Base(file)}
		storeChecks = append(storeChecks, &fileCheck)
		checkBoxes = append(checkBoxes, widget.NewCheck(fileCheck.label, fileCheck.toggle))
	}


	w.SetContent(widget.NewVBox(checkBoxes...))
	w.ShowAndRun()


	for _, checkBox := range storeChecks {
		if checkBox.checked {
			checkedFiles = append(checkedFiles, checkBox.label)
		}
	}


	return checkedFiles
}
