package Controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"plc-backend/File"
	"plc-backend/Shm"

	"github.com/go-playground/validator/v10"
)

func validate(str settings) []*errorResponse {
	validate := validator.New()
	var errors []*errorResponse

	err := validate.Struct(str)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

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

	contents := new(settings)
	_ = json.NewDecoder(r.Body).Decode(&contents)

	// Body has invalid format
	err := validate(*contents)
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
