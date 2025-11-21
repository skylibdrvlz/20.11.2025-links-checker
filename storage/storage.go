package storage

import (
	"encoding/json"
	"fmt"
	"github.com/skylibdrvlz/20.11.2025-links-checker/models"
	"os"
	"sync"
)

type Storage struct {
	mu       sync.RWMutex
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

func (s *Storage) GetLinkSets(ids []int) ([]*models.LinkSet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.LinkSet
	var missingIDs []int

	for _, id := range ids {
		if linkSet, exists := s.linkSets[id]; exists {
			result = append(result, linkSet)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		return nil, fmt.Errorf("IDs not found: %v", missingIDs)
	}
	return result, nil
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
