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
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prometheus/alertmanager/template"
)

func TestService(t *testing.T) {
	data := newAlerts()
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := handler{logger: log.New(os.Stdout, "", log.LstdFlags)}
	http.Handle("/", &handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func newAlerts() template.Data {
	return template.Data{
		Alerts: template.Alerts{
			template.Alert{
				Status:       "critical",
				Annotations:  map[string]string{"a_key": "a_value"},
				Labels:       map[string]string{"l_key": "l_value"},
				StartsAt:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				EndsAt:       time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
				GeneratorURL: "file://generatorUrl",
				Fingerprint:  "3b15fd163d36582e",
			},
			template.Alert{
				Annotations: map[string]string{"a_key_warn": "a_value_warn"},
				Labels:      map[string]string{"l_key_warn": "l_value_warn"},
				Status:      "warning",
			},
		},
		CommonAnnotations: map[string]string{"ca_key": "ca_value"},
		CommonLabels:      map[string]string{"cl_key": "cl_value"},
		GroupLabels:       map[string]string{"gl_key": "gl_value"},
		ExternalURL:       "file://externalUrl",
		Receiver:          "test-receiver",
	}
}
