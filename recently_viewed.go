package main

import (
	"encoding/gob"
	"os"
	"sync"
)

var recentlyViewed struct {
	Packages [10]string // Index 0 is the top (most recently viewed Go package).
	lock     sync.RWMutex

	Production bool
}

// sendToTop sends importPath to top of recentlyViewed.Packages if it's not already present.
func sendToTop(importPath string) {
	recentlyViewed.lock.Lock()
	defer recentlyViewed.lock.Unlock()
	// Check if package is already present, then do nothing.
	for _, p := range recentlyViewed.Packages {
		if p == importPath {
			return
		}
	}
	// Shift all packages down by one.
	for i := len(recentlyViewed.Packages) - 1; i > 0; i-- {
		recentlyViewed.Packages[i] = recentlyViewed.Packages[i-1]
	}
	recentlyViewed.Packages[0] = importPath
}

func saveState(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	recentlyViewed.lock.RLock()
	err = enc.Encode(recentlyViewed.Packages)
	recentlyViewed.lock.RUnlock()

	return err
}

func loadState(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)

	err = dec.Decode(&recentlyViewed.Packages)

	return err
}
