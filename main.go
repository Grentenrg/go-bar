package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/grentenrg/go-bar/widgets"
	"github.com/grentenrg/go-bar/widgets/system"

	layershell "github.com/dlasky/gotk3-layershell/layershell"
)

func main() {
	// Initialize GTK
	gtk.Init(nil)

	bar := NewBar()

	bar.setupStyle()
	bar.setPosition()

	// Create container for status items

	windowNameChan = make(chan string)

	bar.createMainBox()

	var enabledWidgets []Widget

	// Create three sections: left, center, and right
	leftBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create left box:", err)
	}
	leftBox.SetHExpand(true)
	leftBox.SetHAlign(gtk.ALIGN_START)

	centerBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal("Unable to create center box:", err)
	}
	centerBox.SetHExpand(true)
	centerBox.SetHAlign(gtk.ALIGN_CENTER)

	rightBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create right box:", err)
	}
	rightBox.SetHExpand(true)
	rightBox.SetHAlign(gtk.ALIGN_END)

	// Create widgets
	window := widgets.NewWindow()
	if err := window.Create(); err != nil {
		log.Fatal("Unable to create window:", err)
	}
	enabledWidgets = append(enabledWidgets, window)

	workspaces := widgets.NewWorkspace()
	if err := workspaces.Create(); err != nil {
		log.Fatal("Unable to create workspaces:", err)
	}
	enabledWidgets = append(enabledWidgets, workspaces)

	player := widgets.NewPlayer()
	if err := player.Create(); err != nil {
		log.Fatal("Unable to create player:", err)
	}
	enabledWidgets = append(enabledWidgets, player)

	// System widgets
	cpu := system.NewCPU()
	if err := cpu.Create(); err != nil {
		log.Fatal("Unable to create CPU widget:", err)
	}
	enabledWidgets = append(enabledWidgets, cpu)

	memory := system.NewMemory()
	if err := memory.Create(); err != nil {
		log.Fatal("Unable to create memory widget:", err)
	}
	enabledWidgets = append(enabledWidgets, memory)

	disk := system.NewDisk("/")
	if err := disk.Create(); err != nil {
		log.Fatal("Unable to create disk widget:", err)
	}
	enabledWidgets = append(enabledWidgets, disk)

	network := system.NewNetwork("") // Empty string for automatic interface detection
	if err := network.Create(); err != nil {
		log.Fatal("Unable to create network widget:", err)
	}
	enabledWidgets = append(enabledWidgets, network)

	volume := system.NewVolume("")
	if err := volume.Create(); err != nil {
		log.Fatal("Unable to create volume widget:", err)
	}
	enabledWidgets = append(enabledWidgets, volume)

	notification := system.NewNotification()
	if err := notification.Create(); err != nil {
		log.Fatal("Unable to create notification widget:", err)
	}
	enabledWidgets = append(enabledWidgets, notification)

	clock := widgets.NewClock()
	if err := clock.Create(); err != nil {
		log.Fatal("Unable to create clock:", err)
	}
	enabledWidgets = append(enabledWidgets, clock)

	date := widgets.NewDate()
	if err := date.Create(); err != nil {
		log.Fatal("Unable to create date:", err)
	}
	enabledWidgets = append(enabledWidgets, date)

	// Pack widgets in their respective boxes
	leftBox.PackStart(workspaces.Box(), false, false, 5)
	leftBox.PackStart(window.Box(), false, false, 5)

	centerBox.PackStart(player.Box(), false, false, 0)

	// Pack system widgets in the right box
	rightBox.PackStart(cpu.Box(), false, false, 5)
	rightBox.PackStart(memory.Box(), false, false, 5)
	rightBox.PackStart(disk.Box(), false, false, 5)
	rightBox.PackStart(network.Box(), false, false, 5)
	rightBox.PackStart(volume.Box(), false, false, 5)
	rightBox.PackStart(notification.Box(), false, false, 5)
	rightBox.PackEnd(clock.Box(), false, false, 5)
	rightBox.PackEnd(date.Box(), false, false, 5)

	// Add all sections to the main box
	bar.box.PackStart(leftBox, true, true, 0)
	bar.box.PackStart(centerBox, true, true, 0)
	bar.box.PackStart(rightBox, true, true, 0)

	for _, enableWidget := range enabledWidgets {
		styleContext, err := enableWidget.Box().GetStyleContext()
		if err != nil {
			log.Fatal("Unable to get style context:", err)
		}

		styleContext.AddClass(enableWidget.Name())
	}

	// Show all widgets and the window
	bar.window.ShowAll()

	// Connect signals
	bar.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	timer := time.NewTimer(500 * time.Millisecond)

	go func(ws []Widget, timer *time.Timer) {
		for {
			<-timer.C

			for _, widget := range ws {
				if err := widget.Render(); err != nil {
					log.Fatal("Unable to render widget:", err)
				}
			}

			timer.Reset(500 * time.Millisecond)
		}
	}(enabledWidgets, timer)

	// Start the GTK main loop
	gtk.Main()
}

func SetupStyle() {
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Unable to create CSS provider:", err)
	}

	cssFile, err := os.Open("style.css")
	if err != nil {
		log.Fatal("Unable to open CSS file:", err)
	}
	css, err := io.ReadAll(cssFile)
	if err != nil {
		log.Fatal("Unable to read CSS file:", err)
	}

	// Load CSS with font size
	err = cssProvider.LoadFromData(string(css))
	if err != nil {
		log.Fatal("Unable to load CSS:", err)
	}

	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func DockTop(win *gtk.Window) {
	// Initialize gtk-layer-shell for the window
	layershell.InitForWindow(win)

	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)

	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_BOTTOM)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, 0)

	// Set exclusive zone to prevent windows from overlapping with the bar
	layershell.SetExclusiveZone(win, 50) // Match your bar height
}

func getScreenDimensions() (int, int) {
	// Get the default display
	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Fatal("Could not get default display:", err)
	}

	m, err := display.GetMonitor(1)
	if err != nil {
		log.Fatal("Could not get monitor:", err)
	}

	return m.GetGeometry().GetWidth(), m.GetGeometry().GetHeight()
}

type Bar struct {
	window *gtk.Window
	box    *gtk.Box
}

func NewBar() *Bar {
	bar := &Bar{}
	bar.window = bar.createWindow()
	return bar
}

func (b *Bar) createWindow() *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	// Set window properties
	win.SetTitle("Hyprland Status Bar")

	width, height := getScreenDimensions()
	win.SetDefaultSize(width, height)

	return win
}

func (b *Bar) setupStyle() {
	SetupStyle()
}

func (b *Bar) setPosition() {
	DockTop(b.window)
}

func (b *Bar) createMainBox() {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}

	styleContext, err := box.GetStyleContext()
	if err != nil {
		log.Fatal("Unable to get style context:", err)
	}

	styleContext.AddClass("bar")

	b.window.Add(box)
	b.box = box
}

var (
	activeWorkspace int
	windowNameChan  chan string
)
