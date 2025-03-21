package main

import "github.com/gotk3/gotk3/gtk"

type Widget interface {
	Create() error
	Render() error
	Name() string
	Box() *gtk.Box
}
