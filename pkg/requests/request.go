package requests

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Request struct with GET, POST and PUT methods.
type Request struct {
	client *http.Client
}

func newRequest(client *http.Client) *Request {
	return &Request{client: client}
}

// Get method.
func (request *Request) Get(url string, headers map[string]string) ([]byte, *http.Response, error) {
	// Создание запроса
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Добавление заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Отправка запроса
	response, err := request.client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	// Проверка на gzip-body в ответе. Если есть, то декодит
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}

		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	return body, response, err
}

// Post method.
func (request *Request) Post(url string, headers map[string]string, data string) ([]byte, *http.Response, error) {
	dataNew := []byte(data)

	// Создание запроса
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataNew))
	if err != nil {
		return nil, nil, err
	}

	// Добавление заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Отправка запроса
	response, err := request.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Проверка на gzip-body в ответе. Если есть, то декодит
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}

		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	return body, response, err
}

// Put method.
func (request *Request) Put(url string, headers map[string]string, data string) ([]byte, *http.Response, error) {
	dataNew := []byte(data)

	// Создание запроса
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(dataNew))
	if err != nil {
		return nil, nil, err
	}

	// Добавление заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Отправка запроса
	response, err := request.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Проверка на gzip-body в ответе. Если есть, то декодит
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}

		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	return body, response, err
}
