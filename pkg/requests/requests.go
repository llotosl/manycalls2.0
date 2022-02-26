package requests

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Request struct with GET, POST and PUT methods.
type Request struct {
	client *http.Client
}

// NewRequest create Request
func NewRequest(client *http.Client) *Request {
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

// MakeBoundary create boundary body for request.
func MakeBoundary(token string, data map[string]string) (string, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	err := w.SetBoundary("----WebKitFormBoundary" + token)
	if err != nil {
		return "", "", err
	}

	for key, value := range data {
		err = w.WriteField(key, value)
		if err != nil {
			return "", "", err
		}
	}
	defer w.Close()

	return b.String(), w.FormDataContentType(), nil
}

// MakeClient create http.Client for Request.
func MakeClient(proxy string) (*http.Client, *http.Client, error) {
	var client *http.Client
	var clientCookie *http.Client

	if proxy != "" {
		// Создание Transport с прокси.
		ProxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, nil, err
		}
		transport := &http.Transport{Proxy: http.ProxyURL(ProxyURL)}
		// Создание jar-файла куки.
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, nil, err
		}

		client = &http.Client{
			Transport: transport,
		}
		clientCookie = &http.Client{
			Transport: transport,
			Jar:       jar,
		}
	} else {
		// Создание jar-файла куки.
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, nil, err
		}

		client = &http.Client{}
		clientCookie = &http.Client{
			Jar: jar,
		}
	}

	return client, clientCookie, nil
}
