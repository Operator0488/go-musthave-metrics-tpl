package handler

import (
	"net/http"
)

func errorBadRequest(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func errorStatusNotFound(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusNotFound)
}

func errorInternalServer(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
