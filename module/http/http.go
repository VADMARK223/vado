package http

import (
	"fmt"
	"net/http"
)

func test(w http.ResponseWriter, _ *http.Request) {
	str := "Hello world"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success")
	}
}

func StartServer() {
	http.HandleFunc("/test", test)
	fmt.Println("Start server...")
	err := http.ListenAndServe(":9091", nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println("Stop server")
}
