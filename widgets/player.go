package widgets

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Player struct {
	box       *gtk.Box
	label     *gtk.Label
	status    string
	title     string
	artist    string
	position  time.Duration
	duration  time.Duration
	isPlaying bool
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Create() error {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return fmt.Errorf("unable to create box: %w", err)
	}

	p.box = box

	label, err := gtk.LabelNew("")
	if err != nil {
		return fmt.Errorf("unable to create label: %w", err)
	}

	p.label = label
	p.box.PackStart(label, true, true, 0)

	// Start listening for player events
	go p.listenForPlayerEvents()

	return nil
}

func (p *Player) listenForPlayerEvents() {
	for {
		cmd := exec.Command("playerctl", "status", "--format", "{{status}}")
		status, err := cmd.Output()
		if err != nil {
			p.status = "Stopped"
			p.isPlaying = false
			glib.IdleAdd(p.updateLabel)
			time.Sleep(1 * time.Second)
			continue
		}

		p.status = strings.TrimSpace(string(status))
		p.isPlaying = p.status == "Playing"

		// Get metadata
		cmd = exec.Command("playerctl", "metadata", "--format", "{{title}}|{{artist}}")
		metadata, err := cmd.Output()
		if err == nil {
			parts := strings.Split(strings.TrimSpace(string(metadata)), "|")
			if len(parts) >= 2 {
				p.title = parts[0]
				p.artist = parts[1]
			}
		}

		// Get position and duration
		cmd = exec.Command("playerctl", "position")
		pos, err := cmd.Output()
		if err == nil {
			if pos, err := time.ParseDuration(strings.TrimSpace(string(pos)) + "s"); err == nil {
				p.position = pos
			}
		}

		cmd = exec.Command("playerctl", "metadata", "mpris:length")
		dur, err := cmd.Output()
		if err == nil {
			if dur, err := time.ParseDuration(strings.TrimSpace(string(dur)) + "ns"); err == nil {
				p.duration = dur
			}
		}

		glib.IdleAdd(p.updateLabel)
		time.Sleep(1 * time.Second)
	}
}

func (p *Player) updateLabel() {
	if !p.isPlaying {
		p.label.SetText("No media playing")
		return
	}

	// Format time
	pos := formatDuration(p.position)
	dur := formatDuration(p.duration)

	// Truncate long titles
	title := p.title
	if len(title) > 40 {
		title = title[:37] + "..."
	}

	// Truncate long artist names
	artist := p.artist
	if len(artist) > 30 {
		artist = artist[:27] + "..."
	}

	text := fmt.Sprintf("%s - %s [%s/%s]", artist, title, pos, dur)
	p.label.SetText(text)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "00:00"
	}

	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (p *Player) Render() error {
	return nil
}

func (p *Player) Name() string {
	return "player"
}

func (p *Player) Box() *gtk.Box {
	return p.box
}
