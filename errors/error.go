package errors

import (
	"fmt"
	"log"
	"net/http"
)

func HandleErrors(w http.ResponseWriter, err error) {
	log.Printf("%s\n", err.Error())
	emsg := fmt.Sprintf("%s", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(emsg))
}
