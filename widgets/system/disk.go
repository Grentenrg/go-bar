package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Disk struct {
	box    *gtk.Box
	label  *gtk.Label
	path   string
	ticker *time.Ticker
}

func NewDisk(path string) *Disk {
	return &Disk{
		path: path,
	}
}

func (d *Disk) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	d.box = box

	label, err := gtk.LabelNew("💾: ---%")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	box.PackStart(label, false, false, 0)
	d.label = label

	// Start monitoring disk usage
	d.ticker = time.NewTicker(30 * time.Second) // Less frequent updates for disk
	go d.monitor()

	return nil
}

func (d *Disk) monitor() {
	for range d.ticker.C {
		if err := d.updateUsage(); err != nil {
			fmt.Printf("Error updating disk usage: %v\n", err)
			continue
		}
	}
}

func (d *Disk) updateUsage() error {
	cmd := exec.Command("sh", "-c", "df -h "+d.path+" | tail -1 | awk '{print $5}' | sed 's/%//'")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to get disk usage: %w", err)
	}

	usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return fmt.Errorf("unable to parse disk usage: %w", err)
	}

	glib.IdleAdd(func() {
		d.label.SetLabel(fmt.Sprintf("💾 %.1f%%", usage))
	})

	return nil
}

func (d *Disk) Name() string {
	return "disk"
}

func (d *Disk) Box() *gtk.Box {
	return d.box
}

func (d *Disk) Render() error {
	return nil // Updates handled by monitor goroutine
}
