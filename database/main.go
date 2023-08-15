package database

// Stupid simple JSON database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Database struct {
	// The path to the database file
	Path string

	// The data in the database
	Data map[string]interface{}

	// The mutex to lock the database
	Mutex sync.Mutex

	// The file handle to the database
	File *os.File
}

// Create a new database
func NewDatabase(path string) (*Database, error) {
	// Create the database
	db := &Database{
		Path: path,
		Data: make(map[string]interface{}),
	}

	// Load the database
	err := db.Load()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Load the database
func (db *Database) Load() error {
	// Open the database file
	file, err := os.OpenFile(db.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// Set the file handle
	db.File = file

	// Read the database file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Check if the database is empty
	if len(data) == 0 {
		return nil
	}

	// Unmarshal the database
	err = json.Unmarshal(data, &db.Data)
	if err != nil {
		return err
	}

	return nil
}

// Save the database
func (db *Database) Save() error {
	// Marshal the database
	data, err := json.MarshalIndent(db.Data, "", "  ")
	if err != nil {
		return err
	}

	// Write the database
	_, err = db.File.WriteAt(data, 0)
	if err != nil {
		return err
	}

	return nil
}

// Close the database
func (db *Database) Close() error {
	// Close the file handle
	err := db.File.Close()
	if err != nil {
		return err
	}

	return nil
}

// Get a value from the database
func (db *Database) Get(key string, out interface{}) error {
	// Load the database
	err := db.Load()

	if err != nil {
		return err
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Get the value
	value, ok := db.Data[key]
	if !ok {
		return fmt.Errorf("key not found")
	}

	// Marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Unmarshal the value
	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

// Get all values from the database
func (db *Database) GetAll(out interface{}) error {
	// Load the database
	err := db.Load()

	if err != nil {
		return err
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Marshal the values
	data, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	// Unmarshal the value
	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

// Set a value in the database
func (db *Database) Set(key string, value interface{}) error {
	// Load the database
	err := db.Load()

	if err != nil {
		return err
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Set the value
	db.Data[key] = value

	// Save the database
	err = db.Save()

	return err
}

// Delete a value from the database
func (db *Database) Delete(key string) error {
	// Load the database
	err := db.Load()

	if err != nil {
		return err
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Delete the value
	delete(db.Data, key)

	// Save the database
	err = db.Save()

	return err
}

// Check if a key exists in the database
func (db *Database) Exists(key string) bool {
	// Load the database
	err := db.Load()

	if err != nil {
		return false
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Check if the key exists
	_, ok := db.Data[key]

	return ok
}

// Get the keys in the database
func (db *Database) Keys() []string {
	// Load the database
	err := db.Load()

	if err != nil {
		return nil
	}

	// Lock the database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	// Get the keys
	keys := make([]string, 0, len(db.Data))
	for key := range db.Data {
		keys = append(keys, key)
	}

	return keys
}
