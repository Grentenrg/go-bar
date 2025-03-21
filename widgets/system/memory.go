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

type Memory struct {
	box    *gtk.Box
	label  *gtk.Label
	ticker *time.Ticker
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	m.box = box

	label, err := gtk.LabelNew("MEM: ---%")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	box.PackStart(label, false, false, 0)
	m.label = label

	// Start monitoring memory usage
	m.ticker = time.NewTicker(2 * time.Second)
	go m.monitor()

	return nil
}

func (m *Memory) monitor() {
	for range m.ticker.C {
		if err := m.updateUsage(); err != nil {
			fmt.Printf("Error updating memory usage: %v\n", err)
			continue
		}
	}
}

func (m *Memory) updateUsage() error {
	cmd := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to get memory usage: %w", err)
	}

	usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return fmt.Errorf("unable to parse memory usage: %w", err)
	}

	glib.IdleAdd(func() {
		m.label.SetLabel(fmt.Sprintf("🧠 %.1f%%", usage))
	})

	return nil
}

func (m *Memory) Name() string {
	return "memory"
}

func (m *Memory) Box() *gtk.Box {
	return m.box
}

func (m *Memory) Render() error {
	return nil // Updates handled by monitor goroutine
}
