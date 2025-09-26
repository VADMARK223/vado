package component

import (
	"fmt"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, _ *http.Request) {
	msg := "All tasks list."
	_, err := w.Write([]byte(msg))
	if err != nil {
		return
	}
}

func slowHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Started slow request...")
	time.Sleep(time.Second * slowRequestDelaySecond)
	str := "Hello from slow handler!"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Finished slow request")
	}
}
