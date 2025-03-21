package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Volume struct {
	box    *gtk.Box
	label  *gtk.Label
	sink   string
	ticker *time.Ticker
}

func NewVolume(sink string) *Volume {
	if sink == "" {
		sink = "@DEFAULT_SINK@"
	}
	return &Volume{
		sink: sink,
	}
}

func (v *Volume) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	v.box = box

	label, err := gtk.LabelNew("VOL: 0%")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	// Create event box for click handling
	eventBox, err := gtk.EventBoxNew()
	if err != nil {
		return fmt.Errorf("unable to create event box: %w", err)
	}

	// Enable scroll events
	eventBox.AddEvents(int(gdk.SCROLL_MASK))

	eventBox.Add(label)
	eventBox.Connect("button-press-event", v.handleClick)
	eventBox.Connect("scroll-event", v.handleScroll)

	box.PackStart(eventBox, false, false, 0)
	v.label = label

	// Start monitoring volume
	v.ticker = time.NewTicker(1 * time.Second)
	go v.monitor()

	return nil
}

func (v *Volume) monitor() {
	for range v.ticker.C {
		if err := v.updateUsage(); err != nil {
			fmt.Printf("Error updating volume: %v\n", err)
			continue
		}
	}
}

func (v *Volume) updateUsage() error {
	cmd := exec.Command("pactl", "get-sink-volume", v.sink)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to get volume: %w", err)
	}

	// Parse volume percentage
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Volume:") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}

			// Extract percentage value
			percent := strings.TrimSuffix(fields[4], "%")
			volume, err := strconv.Atoi(percent)
			if err != nil {
				return fmt.Errorf("unable to parse volume: %w", err)
			}

			// Check if muted
			cmd = exec.Command("pactl", "get-sink-mute", v.sink)
			output, err = cmd.Output()
			if err != nil {
				return fmt.Errorf("unable to get mute status: %w", err)
			}

			muted := strings.Contains(string(output), "yes")

			glib.IdleAdd(func() {
				icon := "🔊"
				if muted {
					icon = "🔇"
				} else if volume == 0 {
					icon = "🔈"
				} else if volume < 50 {
					icon = "🔉"
				}
				v.label.SetLabel(fmt.Sprintf("%s %d%%", icon, volume))
			})
			break
		}
	}

	return nil
}

func (v *Volume) handleClick(event *gtk.EventBox, eventBtn *gdk.Event) bool {
	buttonEvent := gdk.EventButtonNewFromEvent(eventBtn)
	if buttonEvent.Button() == 1 { // Left click
		// Toggle mute
		cmd := exec.Command("pactl", "set-sink-mute", v.sink, "toggle")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error toggling mute: %v\n", err)
		}
	}
	return true
}

func (v *Volume) handleScroll(event *gtk.EventBox, scrollEvent *gdk.Event) bool {
	scroll := gdk.EventScrollNewFromEvent(scrollEvent)
	direction := scroll.Direction()

	if direction == gdk.SCROLL_UP || direction == gdk.SCROLL_SMOOTH && scroll.DeltaY() < 0 {
		// Increase volume by 5%
		cmd := exec.Command("pactl", "set-sink-volume", v.sink, "+5%")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error increasing volume: %v\n", err)
		}
	} else if direction == gdk.SCROLL_DOWN || direction == gdk.SCROLL_SMOOTH && scroll.DeltaY() > 0 {
		// Decrease volume by 5%
		cmd := exec.Command("pactl", "set-sink-volume", v.sink, "-5%")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error decreasing volume: %v\n", err)
		}
	}
	return true
}

func (v *Volume) Name() string {
	return "volume"
}

func (v *Volume) Box() *gtk.Box {
	return v.box
}

func (v *Volume) Render() error {
	return nil // Updates handled by monitor goroutine
}
