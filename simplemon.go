package simplemon

import (
	"encoding/json"
	"net/http"
)

const (
	good = "VERYHAPPY"
	bad  = "MUCHSAD"
)

var AllChecks = map[string]func() error{
	"backups":   checkBackups,
	"load":      checkLoad,
	"openfiles": checkOpenFiles,
	"disk":      checkDisk,
}

type response struct {
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Headers any      `json:"headers"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var resp response
	for _, f := range AllChecks {
		if e := f(); e != nil {
			resp.Errors = append(resp.Errors, e.Error())
		}
	}

	resp.Headers = r.Header

	w.Header().Add("Content-Type", "application/json")

	if len(resp.Errors) > 0 {
		resp.Status = bad
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.Status = good
		w.WriteHeader(http.StatusOK)
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	_, _ = w.Write(data)
}
