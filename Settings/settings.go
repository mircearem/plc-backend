package Utils

import (
	"encoding/json"
	"net/http"
	"os"
	"plc-backend/File"
	"plc-backend/Shm"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Settings struct {
	Auto  *bool   `json:"Auto" validate:"required"`
	Ratio float32 `json:"Ratio" validate:"required,numeric,min=1,max=100"`
	Kp    float32 `json:"Kp" validate:"required,numeric"`
	Tn    float32 `json:"Tn" validate:"required,numeric"`
	Tv    float32 `json:"Tv" validate:"required,numeric"`
}

type SettingsHandler struct {
	sync.Mutex
	Settings Settings
}

func ValidateSettings(set Settings) []*errorResponse {
	validate := validator.New()
	var errors []*errorResponse

	err := validate.Struct(set)

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

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (s *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	response, _ := json.Marshal(s.Settings)
	defer s.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (s *SettingsHandler) Set(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	settings := new(Settings)
	_ = json.NewDecoder(r.Body).Decode(&settings)

	valid_err := ValidateSettings(*settings)

	// Invalid settings format, exit
	if valid_err != nil {
		resp, _ := json.Marshal(valid_err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	// Valid settings
	s.Lock()
	s.Settings = *settings
	defer s.Unlock()

	// Write settings to  settings file
	json, marshal_err := json.Marshal(s.Settings)

	if marshal_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshal_err.Error()))
		return
	}

	// Write settings to shared memory
	shm_write_err := Shm.Write(os.Getenv("SHM-W-DATA"), os.Getenv("SHM-W-LOCK"), string(json))

	if shm_write_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(shm_write_err.Error()))
		return
	}

	// Write settings file
	writer_err := File.Write("settings.json", string(json))

	if writer_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(writer_err.Error()))
		return
	}
}
