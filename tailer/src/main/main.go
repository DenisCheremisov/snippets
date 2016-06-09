package main

import (
	"log"

	"golang.org/x/exp/inotify"
)

func main() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.AddWatch("/home/emacs/Desktop/test", inotify.IN_ALL_EVENTS)
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.AddWatch("/home/emacs/Desktop", inotify.IN_ALL_EVENTS)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <-watcher.Event:
			log.Println(ev, ev.Name)
		case err := <-watcher.Error:
			log.Fatal(err)
		}
	}
}
