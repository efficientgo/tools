package extkingpin

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// logger is an interface compatible with go-kit/logger.
type logger interface {
	Log(keyvals ...interface{}) error
}

// pathOrContent is an interface compatible with PathOrContent.
type pathOrContent interface {
	Content() ([]byte, error)
	Path() string
}

// PathContentReloader starts a file watcher that monitors the file indicated by pathOrContent.Path() and runs
// reloadFunc whenever a change is detected.
// A debounce timer can be configured via function args to handle situations where many events that would trigger
// a reload are receive in a short period of time. Files will be effectively reloaded at the latest after 2 times
// the debounce timer. By default the debouncer timer is 1 second.
// To ensure renames and deletes are properly handled, the file watcher is put at the file's parent folder. See
// https://github.com/fsnotify/fsnotify/issues/214 for more details.
func PathContentReloader(ctx context.Context, fileContent pathOrContent, debugLogger logger, errorLogger logger, reloadFunc func(), debounceTime time.Duration) error {
	filePath, err := filepath.Abs(fileContent.Path())
	if err != nil {
		return errors.Wrap(err, "getting absolute file path")
	}

	watcher, err := fsnotify.NewWatcher()
	if filePath == "" {
		_ = debugLogger.Log("msg", "no path detected for config reload")
	}
	if err != nil {
		return errors.Wrap(err, "creating file watcher")
	}
	go func() {
		var reloadTimer *time.Timer
		if debounceTime != 0 {
			reloadTimer = time.AfterFunc(debounceTime, func() {
				reloadFunc()
				_ = debugLogger.Log("msg", "configuration reloaded after debouncing")
			})
		}
		defer watcher.Close()
		for {
			select {
			case <-ctx.Done():
				if reloadTimer != nil {
					reloadTimer.Stop()
				}
				return
			case event := <-watcher.Events:
				// fsnotify sometimes sends a bunch of events without name or operation.
				// It's unclear what they are and why they are sent - filter them out.
				if event.Name == "" {
					break
				}
				// We are watching the file's parent folder (more details on why this is done can be found below), but
				// we are only interested in changes to the target file. Discard every other file as quickly as possible.
				if event.Name != filePath {
					break
				}
				// We only react to files being written or created.
				// On "chmod" or "remove" we have nothing to do.
				// On "rename" we have the old file name (not useful). A "create" event for the new file will come later.
				if !event.Op.Has(fsnotify.Write) || !event.Op.Has(fsnotify.Create) {
					break
				}
				_ = debugLogger.Log("msg", fmt.Sprintf("change detected for %s", filePath), "eventName", event.Name, "eventOp", event.Op)
				if reloadTimer != nil {
					reloadTimer.Reset(debounceTime)
				}
			case err := <-watcher.Errors:
				_ = errorLogger.Log("msg", "watcher error", "error", err)
			}
		}
	}()
	// We watch the file's parent folder and not the file itself to better handle DELETE and RENAME events. Check
	// https://github.com/fsnotify/fsnotify/issues/214 for more details.
	if err := watcher.Add(path.Dir(filePath)); err != nil {
		return errors.Wrapf(err, "adding path %s to file watcher", filePath)
	}
	return nil
}

// StaticPathContent serves the contents of a given file through the pathOrContent interface. It's useful for tests
// that rely on such interface.
type StaticPathContent struct {
	content []byte
	path    string
}

var _ pathOrContent = (*StaticPathContent)(nil)

// Content returns the static content.
func (t *StaticPathContent) Content() ([]byte, error) {
	return t.content, nil
}

// Path returns the path to the file that contains the content.
func (t *StaticPathContent) Path() string {
	return t.path
}

// NewStaticPathContent creates a new content that can be used to serve a static configuration.
func NewStaticPathContent(fromPath string) (*StaticPathContent, error) {
	content, err := ioutil.ReadFile(fromPath)

	if err != nil {
		return nil, errors.Wrapf(err, "could not load test content: %s", fromPath)
	}
	return &StaticPathContent{content, fromPath}, nil
}

// Rewrite rewrites the file backing this StaticPathContent and swaps the local content cache. The file writing
// is needed to trigger the file system monitor.
func (t *StaticPathContent) Rewrite(newContent []byte) error {
	t.content = newContent
	// Write the file to ensure possible file watcher reloaders get triggered.
	return ioutil.WriteFile(t.path, newContent, 0666)
}
