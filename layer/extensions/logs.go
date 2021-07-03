package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	apiAddr := os.Getenv("AWS_LAMBDA_RUNTIME_API")

	extensionUrl := fmt.Sprintf("http://%s/2020-01-01/extension", apiAddr)
	logsUrl := fmt.Sprintf("http://%s/2020-08-15/logs", apiAddr)

	id, err := register(extensionUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("ExtensionId", id)

	err = logs(id, logsUrl)
	if err != nil {
		panic(err)
	}

	err = next(id, extensionUrl)
	if err != nil {
		panic(err)
	}

}

func register(baseUrl string) (string, error) {
	b := bytes.NewReader([]byte("{ \"events\": [\"SHUTDOWN\"]}"))

	req, err := http.NewRequest(http.MethodPost, baseUrl+"/register", b)
	if err != nil {
		return "", err
	}
	req.Header.Add("Lambda-Extension-Name", "logs")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected statusCode to be %v, got %v", http.StatusOK, resp.StatusCode)
	}

	id := resp.Header.Get("Lambda-Extension-Identifier")
	return id, nil
}

func next(id string, baseUrl string) error {
	req, err := http.NewRequest(http.MethodGet, baseUrl+"/event/next", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Lambda-Extension-Identifier", id)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected statusCode to be %v, got %v", http.StatusOK, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("next", string(b))

	return nil
}

type SubscribePayload struct {
	Types       []string `json:"types"`
	Destination struct {
		Protocol string `json:"protocol"`
		URI      string `json:"URI"`
	} `json:"destination"`
}

func logs(id string, baseUrl string) error {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			defer r.Body.Close()

			fmt.Println("logs response", string(b))
		}))

		err := http.ListenAndServe("sandbox.localdomain:9090", mux)
		if err != nil {
			panic(err)
		}
	}()

	payload := SubscribePayload{
		Types: []string{"function"},
		Destination: struct {
			Protocol string "json:\"protocol\""
			URI      string "json:\"URI\""
		}{
			Protocol: "HTTP",
			URI:      "http://sandbox.localdomain:9090",
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, baseUrl, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Lambda-Extension-Identifier", id)
	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected statusCode to be %v, got %v", http.StatusOK, resp.StatusCode)
	}

	fmt.Println("Logs response code", resp.StatusCode)

	return nil
}
