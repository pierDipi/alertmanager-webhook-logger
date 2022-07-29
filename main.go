/*
 * Copyright (C) 2020 TomTom N.V. (www.tomtom.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/alertmanager/template"
)

type handler struct {
	logger   *log.Logger
	alertsMu sync.Mutex
	alerts   []template.Data
}

func main() {
	address := flag.String("address", ":6725", "address and port of service")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix)
	http.Handle("/", &handler{logger: logger})

	if err := http.ListenAndServe(*address, nil); err != nil {
		logger.Fatalf("failed to start http server: %v", err)
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var alerts template.Data
	err := json.NewDecoder(r.Body).Decode(&alerts)
	if err != nil {
		h.logger.Printf("cannot parse content because of %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.storeAlerts(alerts)

	err = logAlerts(alerts, h.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) storeAlerts(alerts template.Data) {
	h.alertsMu.Lock()
	defer h.alertsMu.Unlock()

	h.alerts = append(h.alerts, alerts)
}

func logAlerts(alerts template.Data, logger *log.Logger) error {
	s, _ := json.Marshal(alerts)
	logger.Println(string(s))
	return nil
}
