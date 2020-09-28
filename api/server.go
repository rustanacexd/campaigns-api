package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Campaign struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

type campaignHandlers struct {
	sync.Mutex
	store map[string]Campaign
}

func (h *campaignHandlers) campaigns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *campaignHandlers) get(w http.ResponseWriter, r *http.Request) {
	campaigns := make([]Campaign, len(h.store))

	h.Lock()
	i := 0
	for _, campaign := range h.store {
		campaigns[i] = campaign
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(campaigns)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *campaignHandlers) getCampaign(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	campaign, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(campaign)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *campaignHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var campaign Campaign
	err = json.Unmarshal(bodyBytes, &campaign)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	campaign.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[campaign.ID] = campaign
	defer h.Unlock()
}

func newcampaignHandlers() *campaignHandlers {
	return &campaignHandlers{
		store: map[string]Campaign{},
	}

}

func main() {
	campaignHandlers := newcampaignHandlers()
	http.HandleFunc("/campaigns", campaignHandlers.campaigns)
	http.HandleFunc("/campaigns/", campaignHandlers.getCampaign)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
