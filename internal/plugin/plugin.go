package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/jvikstedt/awake/internal/domain"
)

// PluginPerformers loads all plugins in target folderPath and
// converts them to performer and calls callback fn with each of them
func PluginPerformers(folderPath string, fn func(domain.Performer, error)) error {
	files, err := findFiles(folderPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			fn(nil, err)
		}
		s, err := p.Lookup("Plugin")
		if err != nil {
			fn(nil, err)
		}

		performer, ok := s.(domain.Performer)
		if !ok {
			fn(nil, fmt.Errorf("Could not convert to Performer"))
		}

		fn(performer, nil)
	}

	return nil
}

func findFiles(folderPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".so" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
