package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type SettingRequest struct {
	Text string `json:"message"`
}

func (s *Service) SettingTweetLength(w http.ResponseWriter, r *http.Request, c config.Config) {
	newLength, err := services.StrToInt(mux.Vars(r)["new_length"])
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	c.MaxLengthTweet = strconv.Itoa(newLength)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	//err = changeENV(newLength, s.DB)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	m := SettingRequest{
		Text: fmt.Sprintf("maximum tweet length is %d now", newLength),
	}
	services.ReturnJSON(w, http.StatusOK, m)
}

func changeTweetLength(newLength int, db *sql.DB) error {

	return nil
}

func (s *Service) InsertSettings(key string, value []byte) error {
	query := `
        INSERT INTO settings (key, value)
        VALUES ($1, $2);
    `
	_, err := s.DB.Exec(query, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetSettings(key string) (*Settings, error) {
	query := `
        SELECT value
        FROM settings
        WHERE key = $1;
    `
	var settingsJSON []byte
	err := s.DB.QueryRow(query, key).Scan(&settingsJSON)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.Unmarshal(settingsJSON, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
