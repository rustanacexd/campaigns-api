package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Campaign struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

type campaignHandlers struct {
	sync.Mutex
	store map[int]Campaign
}

func (h *campaignHandlers) campaigns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
	case "POST":
		h.post(w, r)
	case "DELETE":
		h.remove(w, r)
	case "PUT":
		h.put(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *campaignHandlers) get(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) == 3 {
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
		return
	}

	part, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	h.Lock()
	campaign, ok := h.store[part]
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
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var campaign Campaign
	err = json.Unmarshal(bodyBytes, &campaign)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.Lock()
	lastCampaign := h.store[len(h.store)]
	campaign.ID = lastCampaign.ID + 1
	h.store[campaign.ID] = campaign
	defer h.Unlock()
}

func (h *campaignHandlers) remove(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	parts := strings.Split(r.URL.String(), "/")
	part, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	h.Lock()
	campaign, ok := h.store[part]
	delete(h.store, campaign.ID)
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *campaignHandlers) put(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	parts := strings.Split(r.URL.String(), "/")
	part, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}
	h.Lock()
	campaign, ok := h.store[part]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.Unmarshal(bodyBytes, &campaign)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.store[campaign.ID] = campaign
	defer h.Unlock()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func newcampaignHandlers() *campaignHandlers {
	return &campaignHandlers{
		store: map[int]Campaign{
			1: Campaign{
				1,
				"mam,e",
				"status",
				time.Now(),
			},
			2: Campaign{
				2,
				"mam,e",
				"status",
				time.Now(),
			},
			3: Campaign{
				3,
				"mam,e",
				"status",
				time.Now(),
			},
		},
	}

}

func main() {
	campaignHandlers := newcampaignHandlers()
	http.HandleFunc("/campaigns/", campaignHandlers.campaigns)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
