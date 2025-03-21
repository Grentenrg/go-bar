package widgets

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gotk3/gotk3/gtk"
)

type hyprlandWorkspace struct {
	ID                   int
	Name                 string
	LastActiveWindowName string
	Monitor              string
}

type Workspace struct {
	box        *gtk.Box
	label      *gtk.Label
	workspaces []hyprlandWorkspace
}

func NewWorkspace() *Workspace {
	return &Workspace{
		workspaces: []hyprlandWorkspace{},
	}
}

func (w *Workspace) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("Unable to create box: %w", err)
	}

	w.box = box

	label, err := gtk.LabelNew("Workspaces...")
	if err != nil {
		return fmt.Errorf("Unable to create label: %w", err)
	}

	w.label = label

	box.PackStart(label, true, true, 0)

	w.workspaces = []hyprlandWorkspace{}
	w.workspaces = append(w.workspaces, hyprlandWorkspace{
		ID:                   1,
		Name:                 "1",
		LastActiveWindowName: "Terminal",
		Monitor:              "1",
	})

	w.workspaces = append(w.workspaces, hyprlandWorkspace{
		ID:                   2,
		Name:                 "Web",
		LastActiveWindowName: "Firefox",
		Monitor:              "1",
	})

	return nil
}

func (w *Workspace) Render() error {
	workspaceTemplate := `{{range .}}<span class="ws">{{.Name}}</span>{{end}}`

	tpl, err := template.New("workspace").Parse(workspaceTemplate)
	if err != nil {
		return fmt.Errorf("Unable to parse template: %w", err)
	}

	b := bytes.NewBuffer([]byte{})

	if err := tpl.Execute(b, w.workspaces); err != nil {
		return fmt.Errorf("Unable to execute template: %w", err)
	}

	content := b.String()
	content = strings.TrimSpace(content)

	// content = glib.MarkupEscapeText(content)

	w.label.SetMarkup(content)

	return nil
}

func (w *Workspace) Name() string {
	return "workspace"
}

func (w *Workspace) Box() *gtk.Box {
	return w.box
}
