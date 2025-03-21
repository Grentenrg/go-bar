package widgets

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/grentenrg/go-bar/libs"
)

type hyprlandWorkspace struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	LastActiveWindowName string `json:"lastwindowtitle"`
	Monitor              string `json:"monitor"`
	IsActive             bool   `json:"active"`
}

type monitors struct {
	Name              string `json:"name"`
	Active            bool   `json:"active"`
	ID                int    `json:"id"`
	ActiveWorkspaceID int    `json:"activeWorkspace.id"`
}

type Workspace struct {
	box             *gtk.Box
	workspaces      []hyprlandWorkspace
	buttons         map[string]*gtk.Button
	activeWorkspace string
	activeMonitor   string // Current monitor name
	monitors        []monitors
}

func NewWorkspace() *Workspace {
	return &Workspace{
		workspaces: []hyprlandWorkspace{},
		buttons:    make(map[string]*gtk.Button),
	}
}

func (w *Workspace) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	w.box = box

	// Get current monitor
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting monitors: %v\n", err)
		return err
	}

	var monitors []struct {
		Name   string `json:"name"`
		Active bool   `json:"active"`
		ID     int    `json:"id"`
	}
	if err := json.Unmarshal(output, &monitors); err != nil {
		fmt.Printf("Error parsing monitors JSON: %v\n", err)
		return err
	}

	// Find the active monitor
	for _, m := range monitors {
		if m.Active {
			w.activeMonitor = strconv.Itoa(m.ID)
			break
		}
	}

	w.updateMonitors()

	// Subscribe to Hyprland workspace events
	go w.subscribeToHyprland()

	go libs.ListenForHyprlandEvents(func(eventType, data string) {
		switch eventType {
		case "workspace":
			w.activeWorkspace = data
		case "focusedmon":
			d := strings.Split(data, ",")
			id, err := strconv.Atoi(d[1])
			if err != nil {
				fmt.Printf("Error converting id to int: %v\n", err)
				return
			}
			w.activeMonitor = d[0]
			for _, m := range w.monitors {
				if m.ID == id {
					w.activeMonitor = d[0]
					w.activeWorkspace = (string)(m.ActiveWorkspaceID)
					break
				}
			}
		}
		glib.IdleAdd(func() {
			w.Render()
		})
	})

	return nil
}

func (w *Workspace) subscribeToHyprland() {
	cmd := exec.Command("hyprctl", "workspaces", "-j")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting workspaces: %v\n", err)
		return
	}

	// Parse workspaces and update buttons
	workspaces := parseWorkspaces(string(output))
	w.updateWorkspaces(workspaces)
}

func (w *Workspace) updateMonitors() {

	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting monitors: %v\n", err)
		return
	}

	var monitors []monitors
	if err := json.Unmarshal(output, &monitors); err != nil {
		fmt.Printf("Error parsing monitors JSON: %v\n", err)
		return
	}

	w.monitors = monitors

}

func (w *Workspace) updateWorkspaces(workspaces []hyprlandWorkspace) {
	// Create a map of existing workspace names
	existingWorkspaces := make(map[string]bool)
	for _, ws := range workspaces {
		existingWorkspaces[ws.Name] = true
	}

	// Remove buttons for workspaces that no longer exist
	for name, button := range w.buttons {
		if !existingWorkspaces[name] {
			w.box.Remove(button)
			delete(w.buttons, name)
		}
	}

	// Update or create buttons for current workspaces
	for _, ws := range workspaces {
		button, exists := w.buttons[ws.Name]
		if !exists {
			var err error
			// Create new button if it doesn't exist
			button, err = gtk.ButtonNewWithLabel(ws.Name)
			if err != nil {
				fmt.Printf("Error creating button: %v\n", err)
				continue
			}

			// Connect click handler
			button.Connect("clicked", func() {
				exec.Command("hyprctl", "dispatch", "workspace", ws.Name).Run()
			})

			w.buttons[ws.Name] = button
			w.box.PackStart(button, false, false, 0)
		}

		// Update button style based on state
		styleContext, _ := button.GetStyleContext()
		styleContext.RemoveClass("workspace-active")
		styleContext.RemoveClass("workspace-inactive")
		styleContext.RemoveClass("workspace-other-display")

		if ws.IsActive {
			styleContext.AddClass("workspace-active")
		} else if ws.Monitor == w.activeMonitor {
			styleContext.AddClass("workspace-inactive")
		} else {
			styleContext.AddClass("workspace-other-display")
		}
	}

	w.box.ShowAll()
}

func (w *Workspace) Render() error {
	// Refresh workspace list
	w.subscribeToHyprland()
	return nil
}

func (w *Workspace) Name() string {
	return "workspace"
}

func (w *Workspace) Box() *gtk.Box {
	return w.box
}

func parseWorkspaces(output string) []hyprlandWorkspace {
	var workspaces []hyprlandWorkspace
	err := json.Unmarshal([]byte(output), &workspaces)
	if err != nil {
		fmt.Printf("Error parsing workspaces JSON: %v\n", err)
		return nil
	}
	return workspaces
}
