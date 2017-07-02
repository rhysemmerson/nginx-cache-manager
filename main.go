package main

import ( 
	"os"
	"bufio"
	"fmt"
	"strings"
	"syscall"
	"os/signal"
	"log"
)

func main() {

	const dir string = "/home/rhys/watchthis"
	const port string = "8080"

	cache := NewCache()
	watcher := NewWatcher(dir, cache)
	server := NewServer(cache, port)

	/* cleanup before exiting */
	defer watcher.Close()
	defer cache.Close()
	defer server.Close()
	
	// exit when prompted by 
	exit := make(chan bool, 1)
	
	listenForSignal(exit)
	scanner(exit)
	
	<- exit
}

// listen for signal
func listenForSignal(exit chan bool) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <- sigs
		fmt.Println(sig)
		exit <- true
	}()
}

// exit when prompted by user
func scanner(exit chan bool) {
	reader := bufio.NewReader(os.Stdin)
	
	/* listen for input */
	go func() {
		for {
			fmt.Print("Press q to quit...")
			
			text, _ := reader.ReadString('\n')
			
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			if text == "q" {
				/* return main */
				exit <- true
			}
		}
	}()
}

func check(err error, message string, params ...interface{}) bool {
	if err == nil {
		return false
	}

	message = "[error] " + message

	log.Printf(message, params...)

	return true
}

func checkAndPanic(err error, message string, params ...interface{}) bool {
	if err == nil {
		return false
	}

	message = "[error:panic] " + message
	
	log.Panicf(message, params...)

	return true
}

func checkAndExit(err error, message string, params ...interface{}) bool {
	if err == nil {
		return false
	}

	message = "[error:fatal] " + message

	log.Fatalf(message, params...)

	return true
}