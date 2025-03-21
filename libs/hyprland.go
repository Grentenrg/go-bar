package libs

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gotk3/gotk3/glib"
)

func ListenForHyprlandEvents(updateFunc func(eventType, data string)) {
	path := os.Getenv("XDG_RUNTIME_DIR")

	socketPath := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if socketPath == "" {
		log.Fatal("HYPRLAND_INSTANCE_SIGNATURE environment variable not set")
	}

	socketPath = fmt.Sprintf("%s/hypr/%s/.socket2.sock", path, socketPath)

	fmt.Println("Connecting to Hyprland socket:", socketPath)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal("Failed to connect to Hyprland socket:", err)
	}

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ">>", 2)
			if len(parts) != 2 {
				continue
			}

			eventType := parts[0]
			data := parts[1]

			glib.IdleAdd(func() {
				updateFunc(eventType, data)
			})
		}
		conn.Close()
	}()
}
