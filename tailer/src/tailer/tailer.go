package tailer

import (
	"os"

	"golang.org/x/exp/inotify"
)

// Processor is line processing handler in Tailer
type Processor interface {
	Process([]byte) error
}

type Tailer struct {
	watcher   *inotify.Watcher
	files     map[string]*os.File
	processor Processor
}

func NewTailer(files []string, p Processor) (res *Tailer, err error) {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if watcher != nil {
				watcher.Close()
			}
		}
	}()

	filesMap := make(map[string]*os.File)
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		filesMap[fileName] = file
		err = watcher.AddWatch(
			fileName,
			inotify.IN_DELETE_SELF|
				inotify.IN_CLOSE_WRITE|
				inotify.IN_MODIFY|
				inotify.IN_MOVE_SELF,
		)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		if err != nil {
			for _, v := range filesMap {
				v.Close()
			}
		}
	}()

	return &Tailer{
		watcher:   watcher,
		files:     filesMap,
		processor: p,
	}, nil
}
