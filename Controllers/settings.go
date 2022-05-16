package Controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"plc-backend/File"
	settings "plc-backend/Settings"
	"plc-backend/Shm"
)

// Wrap sessions handler around request handler
func ReadSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	contents, err := File.Read("settings.json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(contents))
}

func WriteSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	contents := new(settings.Settings)

	_ = json.NewDecoder(r.Body).Decode(&contents)

	// Body has invalid format
	err := contents.Validate()
	if err != nil {
		resp, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	//Body has a valid format

	//Write to shared memory
	jsonObject, _ := json.Marshal(contents)

	shmWriteErr := Shm.Write(os.Getenv("PATH-SHM-W-PATH"), os.Getenv("PATH-SHM-W-LOCK"), string(jsonObject))

	if shmWriteErr != nil {
		w.WriteHeader(http.StatusTooEarly)
		w.Write([]byte(shmWriteErr.Error()))
		return
	}

	//Write to settings file
	writerErr := File.Write("settings.json", string(jsonObject))

	if writerErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(shmWriteErr.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
