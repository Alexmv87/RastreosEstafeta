package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type TrackingResponse struct {
	DatoBusqueda string `json:"dato_busqueda"`
	Guia         string `json:"guia"`
	Fecha        string `json:"fecha"`
	Hora         string `json:"hora"`
	Status       string `json:"status"`
}

type MyRequest struct {
	Items []string `body:"items"`
}

func main() {
	http.HandleFunc("/buscarGuia/{rastreo}", buscarRastreo)
	http.HandleFunc("/buscarGuia", buscarGuia)

	http.ListenAndServe(":9000", nil)
}

func buscarGuia(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "error leyendo body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var guias []string

	if err := json.Unmarshal(body, &guias); err == nil {
		// concurrencia limitada
		//maxConcurrent := 5
		//sem := make(chan struct{}, maxConcurrent)
		var wg sync.WaitGroup
		resCh := make(chan TrackingResponse, len(guias))

		for _, guia := range guias {
			wg.Add(1)
			go func(g string) {
				defer wg.Done()
				//		sem <- struct{}{}
				resultado := busquedaxRastreo(g)
				resCh <- resultado
				//		<-sem
			}(guia)
		}

		wg.Wait()
		close(resCh)

		var resultados []TrackingResponse
		for r := range resCh {
			resultados = append(resultados, r)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resultados)
		return
	}

}

func busquedaxGuia(guia string) TrackingResponse {

	c := colly.NewCollector()

	var tr TrackingResponse
	var found bool

	c.OnHTML("div.historyEventRow", func(h *colly.HTMLElement) {
		if found {
			return
		}

		found = true
		tr.Guia = guia
		tr.Fecha = strings.TrimSpace(h.DOM.Find("div.col-xs-2").First().Text())
		tr.Hora = strings.TrimSpace(h.DOM.Find("div.col-sm-2").First().Text())
		tr.Status = strings.TrimSpace(h.DOM.Find("div.col-sm-7").First().Text())

	})

	c.OnHTML("div.HistoryNoInfo", func(h *colly.HTMLElement) {
		tr.Guia = guia
		tr.Status = "No se encontró información para la guía proporcionada."

	})

	c.Post("https://cs.estafeta.com/es/Tracking/GetTrackingItemHistory", map[string]string{
		"waybill": guia,
	})

	return tr

}

func buscarRastreo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	codigo := strings.Split(r.URL.String(), "/")[len(strings.Split(r.URL.String(), "/"))-1]

	result := busquedaxRastreo(codigo)
	json.NewEncoder(w).Encode(result)
}

func busquedaxRastreo(rastreo string) TrackingResponse {
	c := colly.NewCollector()
	var tr TrackingResponse
	tr.Guia = rastreo

	c.OnHTML("div.shipmentOperationsDiv div.col-sm-5 input", func(h *colly.HTMLElement) {
		guia := h.Attr("data-guia")
		if guia == "" {
			tr.Guia = rastreo
			tr.Status = "No se encontró información para la guía proporcionada."
			return
		}
		tr = busquedaxGuia(guia)
	})

	c.OnHTML("div.titleError", func(h *colly.HTMLElement) {
		tr.Guia = rastreo
		tr.Status = "No se encontró información para la guía proporcionada."
	})

	url := fmt.Sprintf("https://cs.estafeta.com/es/Tracking/searchByGet?wayBill=%s&wayBillType=0&isShipmentDetail=False", rastreo)
	c.Visit(url)
	tr.DatoBusqueda = rastreo
	return tr
}
