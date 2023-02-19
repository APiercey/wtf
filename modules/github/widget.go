package github

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	ghb "github.com/google/go-github/v32/github"
)

// Widget define wtf widget to register widget later
type Widget struct {
	view.MultiSourceWidget
	view.TextWidget

	Client *ghb.Client
	PullRequests []*ghb.PullRequest
	ReviewRequests []*ghb.PullRequest

	settings *Settings
	Selected int
	maxItems int
	Items    []int
}

// NewWidget creates a new instance of the widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		MultiSourceWidget: view.NewMultiSourceWidget(settings.Common, "repository", "repositories"),
		TextWidget:        view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		settings: settings,
	}

	widget.initializeKeyboardControls()

	widget.View.SetRegions(true)
	widget.SetDisplayFunction(widget.display)

	widget.Unselect()

	widget.Sources = widget.settings.repositories

	client, err := NewGithubClient(
									widget.settings.apiKey,
									widget.settings.baseURL,
									widget.settings.uploadURL,
								).githubClient()

	if err != nil {
		panic("Could not init GH Client")
	}

	widget.Client = client

	widget.reloadData()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// SetItemCount sets the amount of PRs RRs and other PRs throughout the widgets display creation
func (widget *Widget) SetItemCount(items int) {
	widget.maxItems = items
}

// GetItemCount returns the amount of PRs RRs and other PRs calculated so far as an int
func (widget *Widget) GetItemCount() int {
	return widget.maxItems
}

// GetSelected returns the index of the currently highlighted item as an int
func (widget *Widget) GetSelected() int {
	if widget.Selected < 0 {
		return 0
	}
	return widget.Selected
}

// Next cycles the currently highlighted text down
func (widget *Widget) Next() {
	widget.Selected++
	if widget.Selected >= widget.maxItems {
		widget.Selected = 0
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Prev cycles the currently highlighted text up
func (widget *Widget) Prev() {
	widget.Selected--
	if widget.Selected < 0 {
		widget.Selected = widget.maxItems - 1
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Unselect stops highlighting the text and jumps the scroll position to the top
func (widget *Widget) Unselect() {
	widget.Selected = -1
	widget.View.Highlight()
	widget.View.ScrollToBeginning()
}

// Refresh reloads the github data via the Github API and reruns the display
func (widget *Widget) Refresh() {
	widget.reloadData()

	widget.display()
}

func (widget *Widget) reloadData() {
	prs ,_ := loadPullRequests(widget.Client, widget.settings.username)
	reviews ,_ := loadReviewRequests(widget.Client, widget.settings.username)

	widget.SetItemCount(len(prs) + len(reviews))
	widget.PullRequests = prs
	widget.ReviewRequests = reviews
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) openPr() {
	if widget.Selected >= 0 && len(widget.Items) > 0 {
		utils.OpenFile(*widget.PullRequests[widget.Selected].HTMLURL)
	}
}
