package backup

import (
	"path/filepath"
	"time"
)

type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := m.Archiver.DestFmt()(time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
