package main

import ( 

	"github.com/fsnotify/fsnotify"

	"log"

	"regexp"

	"io/ioutil"

	"os"
)

func NewWatcher(dir string, cache *Cache) *CacheWatcher {
	w := CacheWatcher{}
	
	w.done = make(chan bool, 1)
	w.cache = cache
	
	w.watch(dir)

	w.scanCache(dir)

	return &w
}

type CacheWatcher struct {
	fileWatcher *fsnotify.Watcher
	cache *Cache
	done chan bool
}

// Close the watcher
func (watcher *CacheWatcher) Close() {
	watcher.done <- true			// signal watch worker
	watcher.fileWatcher.Close()		// close file watcher
}

func (w *CacheWatcher) scanCache(dir string) {
	// get directories
	d, err := os.Open(dir)

	checkAndPanic(err, "could not open cache directory %s ", dir)

	w.fileWatcher.Add(dir)

	// add subdirectories 

	dirs, err := d.Readdir(0)

	// foreach
	for _, folder := range dirs {
		log.Printf("Watching directory %s/%s \n", dir, folder.Name())
		w.fileWatcher.Add(dir + "/" + folder.Name())
	}
}

// Watch cache 
func (watcher *CacheWatcher) watch(dir string) {

	fileWatcher, err := fsnotify.NewWatcher()

	if (err != nil) {
		log.Fatal(err)
	}

	watcher.fileWatcher = fileWatcher

	/* worker */
	go func() {
		for {
			select {
				case <- watcher.done: 
					return

				case event := <-fileWatcher.Events:
					log.Println("event:", event)
					
					if (event.Op&fsnotify.Create == fsnotify.Create ) {
						watcher.onCreate(event)
					}

					// send update
					if (event.Op & fsnotify.Write == fsnotify.Write) {
						watcher.onWrite(event)
					}

					// send remove 
					if (event.Op&fsnotify.Remove == fsnotify.Remove) {
						watcher.onRemove(event)
					}

				case err := <-fileWatcher.Errors: 
					log.Println("error:", err)
			}
		}
	}()

}

func (w *CacheWatcher) onCreate(event fsnotify.Event) {
	log.Println("created file:", event.Name)
						
	f, err := os.Open(event.Name)

	checkAndPanic(err, "couldn't open file %s", event.Name)

	finfo, err := f.Stat()
	
	// if it's a directory add it and return
	if finfo.IsDir() {
		log.Printf("adding directory to watcher: %s ", event.Name)

		w.fileWatcher.Add(event.Name)
	}
}

// look for key and fire cache event
func (w *CacheWatcher) onWrite(event fsnotify.Event) {
	log.Println("modified file:", event.Name)
						
	f, err := os.Open(event.Name)

	checkAndPanic(err, "couldn't open file %s", event.Name)

	finfo, err := f.Stat()

	if finfo.IsDir() {
		return
	}

	key := getKeyFromFile(f)
	
	if key != "" {
		w.cache.Events <- CacheEvent{Key: key, File: event.Name, Op: UPDATE}
	}
}

func (w *CacheWatcher) onRemove(event fsnotify.Event) {
	w.cache.Events <- CacheEvent{Key: "", File: event.Name, Op: DELETE}
}

// Get cache key from file
func getKeyFromFile(file *os.File) string {

	contents, err := ioutil.ReadFile(file.Name())

	if check(err, "error reading file: ", file.Name()) {
		return ""
	}

	r, _ := regexp.Compile("KEY: (.+)")
	
	matches := r.FindStringSubmatch(string(contents))

	var key string 
	
	if len(matches) > 1 {
		key = matches[1]
	}
	
	if key == "" {
		log.Println("Could not find key in file")
	}

	return key
}