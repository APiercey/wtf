package tmuxinator

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
	"github.com/wtfutil/wtf/wtf"
	tc "github.com/wtfutil/wtf/modules/tmuxinator/client"
)

const (
	modalHeight = 7
	modalWidth  = 80
	offscreen   = -1000
)

type Widget struct {
	pages         *tview.Pages

	settings *Settings
	Selected int
	maxItems int
	Items    []string

	tviewApp      *tview.Application
	view.ScrollableWidget
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.common),

		tviewApp:      tviewApp,
		settings: settings,
		pages:         pages,
	}

	widget.initializeKeyboardControls()

	widget.Items = tc.ProjectList()

	widget.Unselect()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) GetSelected() int {
	if widget.Selected < 0 {
		return 0
	}

	return widget.Selected
}

func (widget *Widget) MaxItems() int {
	return len(widget.Items)
}

func (widget *Widget) Refresh() {
	widget.Items = tc.ProjectList()
	widget.Unselect()
	widget.display()
}

func (widget *Widget) RowColor(idx int) string {
	if widget.View.HasFocus() && (idx == widget.GetSelected()) {
		foreground := widget.CommonSettings().Colors.RowTheme.HighlightedForeground

		return fmt.Sprintf(
			"%s:%s",
			foreground,
			widget.CommonSettings().Colors.RowTheme.HighlightedBackground,
		)
	}

	if idx%2 == 0 {
		return fmt.Sprintf(
			"%s:%s",
			widget.settings.common.Colors.RowTheme.EvenForeground,
			widget.settings.common.Colors.RowTheme.EvenBackground,
		)
	}

	return fmt.Sprintf(
		"%s:%s",
		widget.settings.common.Colors.RowTheme.OddForeground,
		widget.settings.common.Colors.RowTheme.OddBackground,
	)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}

func (widget *Widget) content() string {
	cnt := ""

	if len(widget.Items) <= 0 {
		cnt += " [grey]No projects found[white]\n"
	} 

	for idx, projectName := range widget.Items {
		cnt += fmt.Sprintf("[%s] %s \n",
																	widget.RowColor(idx),
																	projectName)
	}

	return cnt
}

func (widget *Widget) processFormInput(prompt string, initValue string, onSave func(string)) {
	form := widget.modalForm(prompt, initValue)

	saveFctn := func() {
		onSave(form.GetFormItem(0).(*tview.InputField).GetText())

		widget.pages.RemovePage("modal")
		widget.tviewApp.SetFocus(widget.View)
		widget.display()
	}

	widget.addButtons(form, saveFctn)
	widget.modalFocus(form)

	// Tell the app to force redraw the screen
	widget.Base.RedrawChan <- true
}

/* -------------------- Modal Form -------------------- */

func (widget *Widget) addButtons(form *tview.Form, saveFctn func()) {
	widget.addSaveButton(form, saveFctn)
	widget.addCancelButton(form)
}

func (widget *Widget) addCancelButton(form *tview.Form) {
	cancelFn := func() {
		widget.pages.RemovePage("modal")
		widget.tviewApp.SetFocus(widget.View)
		widget.display()
	}

	form.AddButton("Cancel", cancelFn)
	form.SetCancelFunc(cancelFn)
}

func (widget *Widget) addSaveButton(form *tview.Form, fctn func()) {
	form.AddButton("Save", fctn)
}

func (widget *Widget) modalFocus(form *tview.Form) {
	frame := widget.modalFrame(form)
	widget.pages.AddPage("modal", frame, false, true)
	widget.tviewApp.SetFocus(frame)

	// Tell the app to force redraw the screen
	widget.Base.RedrawChan <- true
}

func (widget *Widget) modalForm(lbl, text string) *tview.Form {
	form := tview.NewForm()
	form.SetFieldBackgroundColor(wtf.ColorFor(widget.settings.common.Colors.Background))
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonTextColor(wtf.ColorFor(widget.settings.common.Colors.Text))

	form.AddInputField(lbl, text, 60, nil, nil)

	return form
}

func (widget *Widget) modalFrame(form *tview.Form) *tview.Frame {
	frame := tview.NewFrame(form)
	frame.SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetRect(offscreen, offscreen, modalWidth, modalHeight)
	frame.SetBorder(true)
	frame.SetBorders(1, 1, 0, 0, 1, 1)

	drawFunc := func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w, h := screen.Size()
		frame.SetRect((w/2)-(width/2), (h/2)-(height/2), width, height)
		return x, y, width, height
	}

	frame.SetDrawFunc(drawFunc)

	return frame
}
