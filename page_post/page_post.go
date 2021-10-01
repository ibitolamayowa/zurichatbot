package page_post

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ibitolamayowa/zurichatbot/utils"
)

type PostRequest struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

func PostToPage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page_id := params["page_id"]
	var pq PostRequest
	err := utils.ParseJsonFromRequest(r, &pq)
	if err != nil {
		utils.GetError(err, http.StatusUnprocessableEntity, w)
		return
	}
	method := "POST"
	Url, err := url.Parse("https://graph.facebook.com/" + page_id + "/feed")
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Add("message", pq.Message)
	parameters.Add("access_token", pq.AccessToken)
	Url.RawQuery = parameters.Encode()
	client := &http.Client{}
	req, err := http.NewRequest(method, Url.String(), nil)
	if err != nil {
		utils.GetError(err, http.StatusInternalServerError, w)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		utils.GetError(err, http.StatusInternalServerError, w)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.GetError(err, http.StatusInternalServerError, w)
		return
	}
	utils.GetSuccess("posted succesfully", string(body), w)
}
