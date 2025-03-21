package widgets

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type Date struct {
	label      *gtk.Label
	box        *gtk.Box
	lastChange time.Time
}

func NewDate() *Date {
	return &Date{}
}

func (c *Date) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	c.box = box

	elem, err := gtk.LabelNew(time.Now().Format("Monday 02 January 2006"))
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	box.PackStart(elem, true, true, 0)
	c.label = elem
	return nil
}

func (c *Date) Render() error {
	if time.Now().Format("02") == c.lastChange.Format("02") {
		return nil
	}

	c.label.SetText(time.Now().Format("Monday 02 January 2006"))
	c.lastChange = time.Now()
	return nil
}

func (c *Date) Name() string {
	return "date"
}

func (c *Date) Box() *gtk.Box {
	return c.box
}
