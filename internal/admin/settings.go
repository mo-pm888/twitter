package admin

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type SettingRequest struct {
	Text string `json:"message"`
}

func SettingTweetLength(w http.ResponseWriter, r *http.Request, c config.Config) {
	newLength, err := services.StrToInt(mux.Vars(r)["new_length"])
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	c.MaxLengthTweet = strconv.Itoa(newLength)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	err = changeENV(newLength)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	m := SettingRequest{
		Text: fmt.Sprintf("maximum tweet length is %d now", newLength),
	}
	services.ReturnJSON(w, http.StatusOK, m)
}

func changeENV(newLength int) error {
	file, err := os.OpenFile(".env", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening env file:", err)
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	var newLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TWEET_MAX_LENGTH=") {
			newLines = append(newLines, "TWEET_MAX_LENGTH="+strconv.Itoa(newLength))
		} else {
			newLines = append(newLines, line)
		}
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	for _, line := range newLines {
		_, err = fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}
