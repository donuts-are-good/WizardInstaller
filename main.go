package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	wizard := NewWizard(a)
	wizard.Show()
}

type Wizard struct {
	Window      fyne.Window
	CurrentView int
	Views       []*fyne.Container
}

func NewWizard(a fyne.App) *Wizard {
	w := a.NewWindow(app_name)
	w.Resize(fyne.NewSize(500, 400))
	w.SetFixedSize(true)

	wizard := &Wizard{
		Window:      w,
		CurrentView: 0,
		Views:       make([]*fyne.Container, 0),
	}

	wizard.AddView(widget.NewLabel(msg_welcome))
	wizard.AddScrollableView(msg_terms_of_service)
	wizard.AddScrollableView(msg_privacy_statement)
	wizard.AddScrollableView(license_text)
	wizard.AddView(widget.NewLabel(msg_post_install))

	return wizard
}

func (w *Wizard) generateLayout(content fyne.CanvasObject, buttons ...fyne.CanvasObject) *fyne.Container {
	scroll := container.NewVScroll(content)
	scrollContainer := container.NewMax(scroll)

	buttonsContainer := container.NewHBox(layout.NewSpacer())
	for _, button := range buttons {
		buttonsContainer.Add(button)
	}

	return container.NewBorder(nil, buttonsContainer, nil, nil, scrollContainer)
}

func (w *Wizard) AddView(content fyne.CanvasObject) {
	btnOK := widget.NewButton("OK", w.Next)
	btnCancel := widget.NewButton("Cancel", w.Cancel)
	layout := w.generateLayout(content, btnCancel, btnOK)
	w.Views = append(w.Views, layout)
}

func (w *Wizard) SetWindowContent() {
	w.Window.SetContent(w.Views[w.CurrentView])
}

func (w *Wizard) AddScrollableView(text string) {
	entry := widget.NewEntry()
	entry.MultiLine = true
	entry.SetText(text)
	w.AddView(entry)
}

func (w *Wizard) Next() {
	if w.CurrentView < len(w.Views)-1 {
		w.CurrentView++
	} else {
		w.Finish()
		return
	}
	w.SetWindowContent()
}

func (w *Wizard) Cancel() {
	w.Window.Close()
}

func (w *Wizard) Finish() {
	w.Window.Close()
}

func (w *Wizard) Show() {
	w.Window.SetContent(w.Views[w.CurrentView])
	w.Window.ShowAndRun()
}

var (
	app_name    = "wizard installer"
	app_email   = "help@wizardinstaller.com"
	app_website = "https://wizardinstaller.com"
	app_credits = "https://github.com/donuts-are-good"
	app_version = "v1.0.0"

	msg_welcome           = "Welcome! You're about to install WizardInstaller. If you want to continue, click NEXT"
	msg_post_install      = "Thank you for installing WizardInstaller. You can close this window by pressing FINISH"
	msg_terms_of_service  = "<the full terms of service text goes here>"
	msg_privacy_statement = "<the full privacy statement text goes here>"

	license_name = "MIT License"
	license_text = "<the full mit license text goes here>"

	dev_name    = "Donut Wizard"
	dev_email   = "donut@example.com"
	dev_website = "donuts-are-good.com"

	installFiles = []InstallFile{}
	addToPath    = true
)

type InstallFile struct {
	FromFilePath string
	ToFilePath   string
	AddToDesktop bool
	Success      bool
}
