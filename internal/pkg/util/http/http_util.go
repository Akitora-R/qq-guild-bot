package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func Get(url string, param map[string]string, header http.Header) (*[]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	for k, v := range param {
		query.Set(k, v)
	}
	if header != nil {
		req.Header = header
	}
	req.URL.RawQuery = query.Encode()
	//parse, err := url2.Parse("http://localhost:1082")
	client := &http.Client{}
	//client.Transport = &http.Transport{Proxy: http.ProxyURL(parse)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return &respBody, nil
}

func BatchGet(urls []string, param map[string]string) <-chan *[]byte {
	respCh := make(chan *[]byte)
	for _, url := range urls {
		go func(url string, param *map[string]string) {
			get, _ := Get(url, *param, nil)
			respCh <- get
			defer close(respCh)
		}(url, &param)
	}
	return respCh
}

func PostMultipartForm(headers, data map[string]string, fileData map[string][]byte, targetUrl string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for fieldName, b := range fileData {
		formFileWriter, err := bodyWriter.CreateFormFile(fieldName, "file")
		if err != nil {
			fmt.Println("boterr writing to buffer")
			return nil, err
		}
		if _, err = formFileWriter.Write(b); err != nil {
			return nil, err
		}
	}

	for f, v := range data {
		field, err := bodyWriter.CreateFormField(f)
		if err != nil {
			fmt.Println("boterr writing to buffer")
			return nil, err
		}
		if _, err = field.Write([]byte(v)); err != nil {
			return nil, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	req, _ := http.NewRequest("POST", targetUrl, bodyBuf)
	req.Header.Set("Content-Type", contentType)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func PostJson(json *[]byte, url string) (*[]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(*json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func SetProxy() {
	_ = os.Setenv("HTTP_PROXY", "http://localhost:1082")
	_ = os.Setenv("HTTPS_PROXY", "http://localhost:1082")
}
