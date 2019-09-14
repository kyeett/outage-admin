package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func markAsHandled(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updateHandled(ids, true)
}

func markAsUnhandled(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updateHandled(ids, false)
}

func updateHandled(ids []int, newStatus bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	for i, order := range orderFailures {
		for _, id := range ids {
			if order.ID == id {
				orderFailures[i].Handled = newStatus
				fmt.Println("order updated", id)
				break
			}
		}
	}
}

var dbMutex sync.Mutex

func listOrderFailures(w http.ResponseWriter, r *http.Request) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	res, err := json.MarshalIndent(&orderFailures, "", "  ")
	if err != nil {
		http.Error(w, "something went wrong string", http.StatusInternalServerError)
		return
	}

	w.Write(res)
}
