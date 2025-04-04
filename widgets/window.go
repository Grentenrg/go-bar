package widgets

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/grentenrg/go-bar/libs"
)

type Window struct {
	box           *gtk.Box
	label         *gtk.Label
	currentWindow string
	changed       bool
}

func NewWindow() *Window {
	return &Window{}
}

func (w *Window) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	w.box = box

	elem, err := gtk.LabelNew("Window")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	// Set label properties for text overflow handling
	elem.SetHExpand(false)
	elem.SetHAlign(gtk.ALIGN_START)
	elem.SetMaxWidthChars(30)
	elem.SetLineWrap(false)
	elem.SetEllipsize(pango.ELLIPSIZE_END)

	box.PackStart(elem, false, false, 0)

	w.label = elem

	libs.ListenForHyprlandEvents(w.Update)

	return nil
}

func (w *Window) Update(eventType, data string) {
	switch eventType {
	case "activewindow":
		w.currentWindow = data
		w.changed = true
		glib.IdleAdd(func() {
			if err := w.Render(); err != nil {
				fmt.Println("Unable to render window:", err)
			}
		})
	}
}

func (w *Window) Render() error {
	if !w.changed {
		return nil
	}

	w.label.SetLabel(w.currentWindow)
	w.changed = false

	return nil
}

func (w *Window) Name() string {
	return "window"
}

func (w *Window) Box() *gtk.Box {
	return w.box
}

// return label
func (w *Window) Label() *gtk.Label {
	return w.label
}
