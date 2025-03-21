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

type CPU struct {
	box    *gtk.Box
	label  *gtk.Label
	prev   cpuTimes
	usage  float64
	ticker *time.Ticker
}

type cpuTimes struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
}

func NewCPU() *CPU {
	return &CPU{}
}

func (c *CPU) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	c.box = box

	label, err := gtk.LabelNew("CPU: ---%")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	box.PackStart(label, false, false, 0)
	c.label = label

	// Start monitoring CPU usage
	c.ticker = time.NewTicker(2 * time.Second)
	go c.monitor()

	return nil
}

func (c *CPU) monitor() {
	for range c.ticker.C {
		if err := c.updateUsage(); err != nil {
			fmt.Printf("Error updating CPU usage: %v\n", err)
			continue
		}

		glib.IdleAdd(func() {
			c.label.SetLabel(fmt.Sprintf("💻 %.1f%%", c.usage))
		})
	}
}

func (c *CPU) updateUsage() error {
	cmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2}'")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to get CPU usage: %w", err)
	}

	usage, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return fmt.Errorf("unable to parse CPU usage: %w", err)
	}

	glib.IdleAdd(func() {
		c.label.SetLabel(fmt.Sprintf("💻 %.1f%%", usage))
	})

	return nil
}

func (c *CPU) Name() string {
	return "cpu"
}

func (c *CPU) Box() *gtk.Box {
	return c.box
}

func (c *CPU) Render() error {
	return nil // Updates handled by monitor goroutine
}
