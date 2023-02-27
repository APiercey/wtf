package simple_github

import (
	"fmt"

	ghb "github.com/google/go-github/v32/github"
)

func (widget *Widget) display() {
	widget.TextWidget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	// Choses the correct place to scroll to when changing sources
	if len(widget.View.GetHighlights()) > 0 {
		widget.View.ScrollToHighlight()
	} else {
		widget.View.ScrollToBeginning()
	}

	title := widget.CommonSettings().Title

	str := ""

	if widget.settings.showOpenReviewRequests {
		str += fmt.Sprintf("\n [%s]Open Review Requests[white]\n", widget.settings.Colors.Subheading)
		str += widget.displayMyReviewRequests()
	}

	if widget.settings.showMyPullRequests {
		str += fmt.Sprintf("\n [%s]My Pull Requests[white]\n", widget.settings.Colors.Subheading)
		str += widget.displayMyPullRequests()
	}

	return title, str, false
}

func (widget *Widget) displayMyPullRequests() string {
	if !anyAuthoriedRequests(widget.PullRequests) {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range widget.PullRequests {
		if pr.ItemType == AuthoredPullRequest {
			str += fmt.Sprintf(` %s[green]["%d"]%4d[""][white] %s`, widget.mergeString(pr.PullRequest), pr.ID, *pr.PullRequest.Number, *pr.PullRequest.Title)
			str += "\n"
		}
	}


	return str
}

func (widget *Widget) displayMyReviewRequests() string {
	if !anyReviewRequests(widget.PullRequests) {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range widget.PullRequests {
		if pr.ItemType == ReviewRequested {
			str += fmt.Sprintf(` %s[green]["%d"]%4d[""][white] %s`, widget.mergeString(pr.PullRequest), pr.ID, *pr.PullRequest.Number, *pr.PullRequest.Title)
			str += "\n"
		}
	}

	return str
}

var mergeIcons = map[string]string{
	"dirty":    "[red]\u0021[white] ",
	"clean":    "[green]\u2713[white] ",
	"unstable": "[red]\u2717[white] ",
	"blocked":  "[red]\u2717[white] ",
}

func (widget *Widget) mergeString(pr *ghb.PullRequest) string {
	if !widget.settings.enableStatus {
		return ""
	}
	if str, ok := mergeIcons[pr.GetMergeableState()]; ok {
		return str
	}
	return "? "
}
