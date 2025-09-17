package http

import (
	"fmt"
	"net/http"
	"time"
	"vado/constant"
)

func slowHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Started slow request...")
	time.Sleep(time.Second * constant.SlowRequestDelaySecond)
	str := "Hello from slow handler!"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Finished slow request")
	}
}

func queryParamsHandler(w http.ResponseWriter, r *http.Request) {
	nameParam := r.URL.Query().Get("name")
	surnameParam := r.URL.Query().Get("surname")

	msg := fmt.Sprintf("Name: %s, Surname: %s", nameParam, surnameParam)
	_, err := w.Write([]byte(msg))
	if err != nil {
		return
	}
}
