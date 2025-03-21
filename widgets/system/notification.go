package system

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Notification struct {
	box   *gtk.Box
	label *gtk.Label
}

func NewNotification() *Notification {
	return &Notification{}
}

func (n *Notification) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	n.box = box

	label, err := gtk.LabelNew("🔔")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	// Create event box for click handling
	eventBox, err := gtk.EventBoxNew()
	if err != nil {
		return fmt.Errorf("unable to create event box: %w", err)
	}

	eventBox.Add(label)
	eventBox.Connect("button-press-event", n.handleClick)

	box.PackStart(eventBox, false, false, 0)
	n.label = label

	// Start monitoring notifications
	go n.monitor()

	return nil
}

func (n *Notification) monitor() {
	cmd := exec.Command("swaync-client", "-c")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting notification count: %v\n", err)
		return
	}

	count := strings.TrimSpace(string(output))
	if count != "0" {
		glib.IdleAdd(func() {
			n.label.SetLabel(fmt.Sprintf("🔔 %s", count))
		})
	} else {
		glib.IdleAdd(func() {
			n.label.SetLabel("🔔")
		})
	}
}

func (n *Notification) handleClick(event *gtk.EventBox, eventBtn *gdk.Event) bool {
	buttonEvent := gdk.EventButtonNewFromEvent(eventBtn)
	if buttonEvent.Button() == 1 { // Left click
		// Toggle notification center
		cmd := exec.Command("swaync-client", "-t")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error toggling notification center: %v\n", err)
		}
	}
	return true
}

func (n *Notification) Name() string {
	return "notification"
}

func (n *Notification) Box() *gtk.Box {
	return n.box
}

func (n *Notification) Render() error {
	return nil // Updates handled by monitor goroutine
}
