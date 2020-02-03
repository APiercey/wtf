package cdsqueue

import (
	"fmt"
	"strings"
	"time"

	"github.com/ovh/cds/sdk"
)

func (widget *Widget) display() {
	widget.TextWidget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	if len(widget.View.GetHighlights()) > 0 {
		widget.View.ScrollToHighlight()
	} else {
		widget.View.ScrollToBeginning()
	}

	widget.Items = make([]sdk.WorkflowNodeJobRun, 0)
	filter := widget.currentFilter()
	_, _, width, _ := widget.View.GetRect()

	str := widget.settings.common.SigilStr(len(widget.filters), widget.Idx, width) + "\n"
	str += widget.displayQueue(filter)

	title := fmt.Sprintf("%s - %s", widget.CommonSettings().Title, widget.title(filter))

	return title, str, false
}

func (widget *Widget) title(filter string) string {
	return fmt.Sprintf(
		"[%s]%d - %s[white]",
		widget.settings.common.Colors.TextTheme.Title,
		widget.maxItems,
		filter,
	)
}

func (widget *Widget) displayQueue(filter string) string {
	runs, _ := widget.client.QueueWorkflowNodeJobRun(filter)

	widget.SetItemCount(len(runs))

	if len(runs) == 0 {
		return " [grey]none[white]\n"
	}

	var content string
	for idx, job := range runs {
		content += fmt.Sprintf(`[grey]["%d"]%s`,
			idx, widget.generateQueueJobLine(job.ID, job.Parameters, job.Job, time.Since(job.Queued), job.BookedBy, job.Status))

		widget.Items = append(widget.Items, job)
	}

	return content
}

func (widget *Widget) generateQueueJobLine(id int64, parameters []sdk.Parameter, executedJob sdk.ExecutedJob,
	duration time.Duration, bookedBy sdk.Service, status string) string {
	prj := getVarsInPbj("cds.project", parameters)
	workflow := getVarsInPbj("cds.workflow", parameters)
	node := getVarsInPbj("cds.node", parameters)
	run := getVarsInPbj("cds.run", parameters)
	triggeredBy := getVarsInPbj("cds.triggered_by.username", parameters)

	row := make([]string, 6)
	row[0] = pad(fmt.Sprintf(sdk.Round(duration, time.Second).String()), 9)
	row[2] = pad(run, 7)
	row[3] = fmt.Sprintf("%s", pad(prj+"/"+workflow+"/"+node, 40))

	if status == sdk.StatusBuilding {
		row[1] = pad(fmt.Sprintf(" %s.%s ", executedJob.WorkerName, executedJob.WorkerID), 27)
	} else if bookedBy.ID != 0 {
		row[1] = pad(fmt.Sprintf(" %s.%d ", bookedBy.Name, bookedBy.ID), 27)
	} else {
		row[1] = pad("", 27)
	}

	row[4] = fmt.Sprintf("➤ %s", pad(triggeredBy, 17))

	c := "grey"
	if status == sdk.StatusWaiting {
		if duration > 120*time.Second {
			c = "red"
		} else if duration > 50*time.Second {
			c = "yellow"
		}
	}

	return fmt.Sprintf("[%s]%s [grey]%s %s %s %s\n", c, row[0], row[1], row[2], row[3], row[4])
}

func getStatusColor(status string) string {
	switch status {
	case sdk.StatusSuccess:
		return "green"
	case sdk.StatusBuilding, sdk.StatusWaiting:
		return "blue"
	case sdk.StatusFail:
		return "red"
	case sdk.StatusStopped:
		return "red"
	case sdk.StatusSkipped:
		return "grey"
	case sdk.StatusDisabled:
		return "grey"
	}
	return "red"
}

func pad(t string, size int) string {
	if len(t) > size {
		return t[0:size-3] + "..."
	}
	return t + strings.Repeat(" ", size-len(t))
}

func getVarsInPbj(key string, ps []sdk.Parameter) string {
	for _, p := range ps {
		if p.Name == key {
			return p.Value
		}
	}
	return ""
}
