package storage

import (
	"encoding/json"
	"github.com/skylibdrvlz/20.11.2025-links-checker/models"
	"os"
	"sync"
)

type Storage struct {
	mu       sync.Mutex
	linkSets map[int]*models.LinkSet
	nextID   int
	dataFile string
}

func NewStorage(dataFile string) *Storage {
	s := &Storage{
		linkSets: make(map[int]*models.LinkSet),
		nextID:   1,
		dataFile: dataFile,
	}

	s.loadFromFile()
	s.updateNextID()
	return s
}

func (s *Storage) SaveLinkSet(links map[string]string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	s.linkSets[id] = &models.LinkSet{
		ID:    id,
		Links: links,
	}

	s.saveToFile()
	return id
}

func (s *Storage) saveToFile() {
	file, _ := os.Create(s.dataFile)
	defer file.Close()
	json.NewEncoder(file).Encode(s.linkSets)
}

func (s *Storage) loadFromFile() {
	file, err := os.Open(s.dataFile)
	if err != nil {
		return
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&s.linkSets)
}

func (s *Storage) updateNextID() {
	maxID := 0
	for id := range s.linkSets {
		if id > maxID {
			maxID = id
		}
	}
	if maxID > 0 {
		s.nextID = maxID + 1
	}
}
