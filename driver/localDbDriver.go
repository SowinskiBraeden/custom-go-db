package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/jcelliott/lumber"
)

const Version = "1.0.1"
const parentDir = "C:/local-go-db/"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Info(string, ...interface{})
		Warn(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

type Data struct {
	ID     string
	Object interface{}
}

func NewConnection(dbName string, options *Options) (*Driver, error) {
	if dbName == "" {
		return nil, fmt.Errorf("Cannot create Nameless database - No name given")
	}
	dir := parentDir + dbName
	dir = filepath.Clean(dir)

	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) InsertOne(collection string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	// Create new uuid for object and append object to database

	var data Data
	data.ID = uuid.New().String()
	data.Object = v

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, data.ID+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) InsertMany(collection, resource string) {

}

func (d *Driver) FindOne(collection, id string) (interface{}, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to read!")
	}

	record := filepath.Join(d.dir, collection, id)

	if _, err := stat(record); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("The object you're searching for cannot be found")
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return nil, err
	}
	var v Data
	json.Unmarshal(b, &v)

	return v, nil
}

func (d *Driver) FindAll(collection string) ([]string, error) {

	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to read!")
	}
	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection, id string) error {

	path := filepath.Join(collection, id)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to find file or directory name %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

type User struct {
	Name string
	Age  json.Number
}
