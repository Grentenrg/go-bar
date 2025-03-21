package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Network struct {
	box            *gtk.Box
	label          *gtk.Label
	interface_     string
	prevRx, prevTx uint64
	ticker         *time.Ticker
}

func NewNetwork(interface_ string) *Network {
	return &Network{
		interface_: interface_,
	}
}

func (n *Network) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	n.box = box

	label, err := gtk.LabelNew("NET: ↓0KB/s ↑0KB/s")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	box.PackStart(label, false, false, 0)
	n.label = label

	// Start monitoring network usage
	n.ticker = time.NewTicker(1 * time.Second)
	go n.monitor()

	return nil
}

func (n *Network) monitor() {
	for range n.ticker.C {
		if err := n.updateUsage(); err != nil {
			fmt.Printf("Error updating network usage: %v\n", err)
			continue
		}
	}
}

func (n *Network) findActiveInterface() (string, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return "", fmt.Errorf("unable to read /proc/net/dev: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip header lines
		if strings.Contains(line, "Inter-|") || strings.Contains(line, "face |") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Remove the colon from the interface name
		iface := strings.TrimSuffix(fields[0], ":")

		// Skip loopback and virtual interfaces
		if iface == "lo" || strings.HasPrefix(iface, "tun") || strings.HasPrefix(iface, "docker") {
			continue
		}

		// Check if interface is up
		upFile := fmt.Sprintf("/sys/class/net/%s/operstate", iface)
		upData, err := os.ReadFile(upFile)
		if err != nil {
			continue
		}

		if strings.TrimSpace(string(upData)) == "up" {
			return iface, nil
		}
	}

	return "", fmt.Errorf("no active network interface found")
}

func (n *Network) updateUsage() error {
	// If no interface specified or interface is down, try to find an active one
	if n.interface_ == "" {
		iface, err := n.findActiveInterface()
		if err != nil {
			glib.IdleAdd(func() {
				n.label.SetLabel("NET: No active interface")
			})
			return err
		}
		n.interface_ = iface
	}

	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return fmt.Errorf("unable to read /proc/net/dev: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, n.interface_) {
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			rx, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse rx bytes: %w", err)
			}

			tx, err := strconv.ParseUint(fields[9], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse tx bytes: %w", err)
			}

			if n.prevRx > 0 && n.prevTx > 0 {
				rxSpeed := float64(rx-n.prevRx) / 1024 // KB/s
				txSpeed := float64(tx-n.prevTx) / 1024 // KB/s

				glib.IdleAdd(func() {
					n.label.SetLabel(fmt.Sprintf("↓%.1fKB/s ↑%.1fKB/s", rxSpeed, txSpeed))
				})
			}

			n.prevRx = rx
			n.prevTx = tx
			break
		}
	}

	return nil
}

func (n *Network) Name() string {
	return "network"
}

func (n *Network) Box() *gtk.Box {
	return n.box
}

func (n *Network) Render() error {
	return nil // Updates handled by monitor goroutine
}
