package storage

import (
	"encoding/json"
	"os"
)

type CacheRepository struct {
	filePath string
	Data     map[string]string `json:"data"`
	Order    []string          `json:"order"`
	limit    int
	dirty    bool
}

func NewCacheRepository(filePath string, limit int) *CacheRepository {
	return &CacheRepository{
		filePath: filePath,
		Data:     make(map[string]string),
		Order:    make([]string, 0),
		limit:    limit,
	}
}

func (r *CacheRepository) Load() error {
	f, err := os.Open(r.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(r); err != nil {
		return err
	}

	if r.Data == nil {
		r.Data = make(map[string]string)
	}
	if r.Order == nil {
		r.Order = make([]string, 0)
	}

	return nil
}

func (r *CacheRepository) Save() error {
	if !r.dirty {
		return nil
	}

	f, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(r); err != nil {
		return err
	}

	r.dirty = false
	return nil
}

func (r *CacheRepository) Get(term string) (string, bool) {
	code, ok := r.Data[term]
	if ok {
		// Move to front of LRU
		r.moveToFront(term)
	}
	return code, ok
}

func (r *CacheRepository) Set(term, code string) {
	if oldCode, ok := r.Data[term]; ok && oldCode == code {
		r.moveToFront(term)
		return
	}

	r.Data[term] = code
	r.moveToFront(term)
	r.dirty = true

	if len(r.Order) > r.limit {
		toRemove := r.Order[len(r.Order)-1]
		delete(r.Data, toRemove)
		r.Order = r.Order[:len(r.Order)-1]
	}
}

func (r *CacheRepository) Clear() {
	r.Data = make(map[string]string)
	r.Order = make([]string, 0)
	r.dirty = true
}

func (r *CacheRepository) moveToFront(term string) {
	// Find and remove if exists
	idx := -1
	for i, t := range r.Order {
		if t == term {
			idx = i
			break
		}
	}

	if idx != -1 {
		r.Order = append(r.Order[:idx], r.Order[idx+1:]...)
	}

	// Prepend
	r.Order = append([]string{term}, r.Order...)
}
