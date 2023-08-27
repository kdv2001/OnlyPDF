package gotemberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	convertUrl = "/forms/libreoffice/convert"
)

type Repo struct {
	client *http.Client
	url    string
}

func NewRepo(client *http.Client, url string) *Repo {
	return &Repo{
		client: client,
		url:    url,
	}
}

func (r *Repo) ConvertFiles(ctx context.Context, files []string, resultFileName string, needMerge bool) error {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)

	for _, val := range files {
		fw, err := writer.CreateFormFile("files", val)
		if err != nil {
			return err
		}

		fd, err := os.Open(val)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, fd)
		if err != nil {
			return err
		}

		fd.Close()
	}

	if needMerge {
		formField, err := writer.CreateFormField("merge")
		if err != nil {
			log.Fatal(err)
		}

		_, err = formField.Write([]byte("true"))
		if err != nil {
			return err
		}
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprint(r.url, convertUrl), form)
	if err != nil {
		return err
	}

	fmt.Println(writer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	fmt.Println(writer.FormDataContentType())
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status is %d, %s", resp.StatusCode, bodyText)
	}

	f, err := os.Create(resultFileName)
	if err != nil {
		return err
	}

	_, err = f.Write(bodyText)
	if err != nil {
		return err
	}

	return nil
}
