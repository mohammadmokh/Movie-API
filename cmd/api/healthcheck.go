package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheck(rw http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"status": "available",
		"system_info": map[string]string{
			"version:": version,
		},
	}
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	rw.Header().Set("Content-Type:", "application/json")
	rw.Write(js)
}
