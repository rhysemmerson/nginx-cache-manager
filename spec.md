# Implementation

## Files

#### watcher.go
    + Watch(dir) chan

* Encapsulates boiler plate for file watching
* writes to file watcher

#### cache.go
    + Cache() chan

- listens to file watcher
- updates index when files change
- listens to cache channel
- matches and deletes files in index when message received 

#### server.go
    + Server(port)

* Listens for requests
* writes to cache channel

#### main.go
    watcher := Watch(dir)
    cache := Cache(watcher)
    server := Server(cache)
    
## channels 

chan watcherChan = Watch()
chan cacheChan = Cache()

fileWatcher 
* watcher writes to watcherChan 
* cache listens to watcherChan

cacheChan
* cache listens to cacheChan
* server writes to cacheChan

Watch and index
- watch files
- when file changes, read contents
- read key from file
- add/update key in index

Match and delete files
- listen
- recieve request from client
- look for files in index by key
- unlink files