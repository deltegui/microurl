package main

import (
	"encoding/json"
	"fmt"
	"microurl/internal"
	"net/http"
	"os"
	"strings"
	"time"
)

type api struct {
	server, token string
}

var client = &http.Client{Timeout: 40 * time.Second}

func (api api) endpoint(path string) string {
	return fmt.Sprintf("%s/api%s", api.server, path)
}

func (api api) endpointWithID(path string, id uint) string {
	return fmt.Sprintf("%s/api%s/%d", api.server, path, id)
}

func (api api) login(user, password string) string {
	payload := fmt.Sprintf("{\"name\": \"%s\", \"password\": \"%s\"}", user, password)
	res, err := client.Post(api.endpoint("/user/login"), "application/json", strings.NewReader(payload))
	if err != nil {
		errd("Error while creating request: %s\n", err.Error())
	}
	defer res.Body.Close()
	checkAPIError(res)
	var r internal.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		errd("Unkown error while decoding login response: %s\n", err.Error())
	}
	return r.Token.Value
}

func (api api) getAll() []internal.URLResponse {
	api.ensureHaveToken()
	res := api.authReq(http.MethodGet, api.endpoint("/url/all"), "")
	defer res.Body.Close()
	checkAPIError(res)
	var r internal.AllURLsResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		errd("Unkown error while decoding multiple URL response: %s\n", err.Error())
	}
	return r.URLs
}

func (api api) delete(id uint) internal.URLResponse {
	api.ensureHaveToken()
	res := api.authReq(http.MethodDelete, api.endpointWithID("/url", id), "")
	return tryDecodeURL(res)
}

func (api api) create(name, value string) internal.URLResponse {
	api.ensureHaveToken()
	payload := fmt.Sprintf("{\"name\": \"%s\", \"url\": \"%s\"}", name, value)
	res := api.authReq(http.MethodPost, api.endpoint("/url"), payload)
	return tryDecodeURL(res)
}

func (api api) generateQR(id uint) string {
	api.ensureHaveToken()
	res := api.authReq(http.MethodPost, api.endpointWithID("/url/genqr", id), "")
	return tryDecodeURL(res).QR
}

func (api api) ensureHaveToken() {
	if len(api.token) == 0 {
		errd("You need to log in\n")
	}
}

func (api api) authReq(method, endpoint, payload string) *http.Response {
	req, err := http.NewRequest(method, endpoint, strings.NewReader(payload))
	if err != nil {
		errd("Error while getting all urls: %s\n", err.Error())
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", api.token)},
		"Content-Type":  {"application/json"},
	}
	res, err := client.Do(req)
	if err != nil {
		errd("Error while connecting to server: %s\n", err.Error())
	}
	return res
}

func tryDecodeURL(res *http.Response) internal.URLResponse {
	defer res.Body.Close()
	checkAPIError(res)
	var r internal.URLResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		errd("Unkown error while decoding URL response: %s\n", err.Error())
	}
	return r
}

func checkAPIError(res *http.Response) {
	if res.StatusCode == http.StatusOK {
		return
	}
	var e internal.UseCaseError
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		errd("Error while decoding api error: %s\n", err.Error())
	}
	fmt.Println("Error:", e.Reason)
	os.Exit(1)
}
