package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Twitter_like_application/internal/services"
)

type SettingResponse struct {
	Text string `json:"message"`
}

func (s *Service) SettingsTweet(w http.ResponseWriter, r *http.Request) {
	var newSettings Settings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.TweetLength = newSettings.TweetLength
	jsonValue, err := json.Marshal(newSettings)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = s.ChangeSettings("tweet", jsonValue); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := SettingResponse{
		Text: fmt.Sprintf("the maximum tweet length is now %d characters.", newSettings.TweetLength),
	}
	services.ReturnJSON(w, http.StatusOK, msg)
}
func (s *Service) DefaultMaxTweetLength() error {
	defaultLength := &Settings{TweetLength: 400}
	jsonValue, err := json.Marshal(defaultLength)
	if err != nil {
		return err
	}
	if err = s.ChangeSettings("tweet", jsonValue); err != nil {
		return err
	}
	return nil
}
func (s *Service) ChangeSettings(key string, value []byte) error {
	query := `
        INSERT INTO settings (key, value)
        VALUES ($1, $2)
        ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;
    `
	_, err := s.DB.Exec(query, key, value)
	return err
}

func (s *Service) GetSettings(key string) error {
	query := `
        SELECT value
        FROM settings
        WHERE key = $1;
    `
	var settingsJSON []byte
	err := s.DB.QueryRow(query, key).Scan(&settingsJSON)
	if err != nil {
		return err
	}
	var settings Settings
	err = json.Unmarshal(settingsJSON, &settings)
	if err != nil {
		return err
	}
	s.TweetLength = settings.TweetLength
	return nil
}
