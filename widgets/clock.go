package widgets

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type Clock struct {
	box   *gtk.Box
	label *gtk.Label
}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("Unable to create box: %w", err)
	}

	c.box = box

	elem, err := gtk.LabelNew(time.Now().Format("15:04:05"))
	if err != nil {
		return fmt.Errorf("Unable to create label: %w", err)
	}

	box.PackStart(elem, true, true, 0)

	c.label = elem
	return nil
}

func (c *Clock) Render() error {
	c.label.SetText(time.Now().Format("15:04:05Z07"))
	return nil
}

func (c *Clock) Name() string {
	return "clock"
}

func (c *Clock) Box() *gtk.Box {
	return c.box
}
