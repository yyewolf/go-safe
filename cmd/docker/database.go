package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

var databaseFile string

type File struct {
	Sum string `json:"s"`
}

var database map[string]*File
var databaseDigest string

func loadDatabase(f string) {
	databaseFile = f
	database = make(map[string]*File)

	// Read the database file
	data, err := os.ReadFile(databaseFile)
	if err != nil {
		return
	}

	sum := sha256.Sum256(data)
	databaseDigest = hex.EncodeToString(sum[:])

	json.Unmarshal(data, &database)
}

func saveDatabase() error {
	data, err := json.Marshal(database)
	if err != nil {
		return err
	}

	return os.WriteFile(databaseFile, data, 0644)
}
