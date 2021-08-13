package onedrive

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/goh-chunlin/go-onedrive/onedrive"
)

type Response struct {
	ID                   string                    `json:"id"`
	Etag                 string                    `json:"eTag"`
	LastModifiedDateTime time.Time                 `json:"lastModifiedDateTime"`
	Size                 int64                     `json:"size"`
	CreatedBy            CreatedBy                 `json:"createdBy"`
	File                 *onedrive.DriveItemFile   `json:"file"`
	Folder               *onedrive.DriveItemFolder `json:"folder"`
}

type CreatedBy struct {
	User onedrive.User `json:"user"`
}

//Because Onedrive only use OAuth2.0

//https://docs.microsoft.com/en-us/graph/auth/auth-concepts

func getObject(ctx context.Context, client *onedrive.Client, absPath string) (*onedrive.DriveItem, error) {
	if !path.IsAbs(absPath) {
		return nil, errors.New("not an absoulute path")
	}

	apiURL := "me/drive/root:" + url.PathEscape(absPath)
	//fmt.Println(apiURL)

	req, err := client.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var driveItem *onedrive.DriveItem
	err = client.Do(ctx, req, false, &driveItem)
	if err != nil {
		return nil, err
	}

	return driveItem, nil
}

func Download(ctx context.Context, client *onedrive.Client, absPath string) (io.ReadCloser, error) {
	i, err := getObject(ctx, client, absPath)
	if err != nil {
		return nil, err
	}

	var body io.Reader

	req, err := http.NewRequest("GET", i.DownloadURL, body)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Delete(ctx context.Context, client *onedrive.Client, absPath string) (*onedrive.DriveItem, error) {
	if !path.IsAbs(absPath) {
		return nil, errors.New("not an absoulute path")
	}

	apiURL := "me/drive/root:/" + url.PathEscape(absPath) + ":"

	req, err := client.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var driveItem *onedrive.DriveItem
	err = client.Do(ctx, req, false, &driveItem)
	if err != nil {
		return nil, err
	}

	return driveItem, nil
}

func Upload(ctx context.Context, absPath string, client *onedrive.Client, fileSize int64, r io.Reader) (*onedrive.DriveItem, error) {

	if !path.IsAbs(absPath) {
		return nil, errors.New("Please provide the  abspath")
	}

	//_, err := getObject(ctx, client, absPath)
	//fmt.Println(item, err)

	apiURL := "me/drive/root:" + url.PathEscape(absPath) + ":/content"
	//fmt.Println(apiURL)

	//req, err := client.NewRequest("PUT", apiURL, r)
	buffer := make([]byte, fileSize)
	r.Read(buffer)
	fileReader := bytes.NewReader(buffer)

	req, err := client.NewFileUploadRequest(apiURL, "application/octet-stream", fileReader)
	//fmt.Println(r)
	//req.Header.Add("Content-Type", "text/plain")
	//fmt.Println(item.File.MIMEType)
	if err != nil {
		return nil, err
	}

	var response *onedrive.DriveItem
	err = client.Do(ctx, req, false, &response)
	//fmt.Println(err)
	if err != nil {
		return nil, err
	}

	//fmt.Println(response)
	return response, nil
}

func getSomeObjectPart(ctx context.Context, client *onedrive.Client, absPath string) (*Response, error) {
	if !path.IsAbs(absPath) {
		return nil, errors.New("not an absoulute path")
	}

	apiURL := "me/drive/root:" + url.PathEscape(absPath)
	//fmt.Println(apiURL)

	req, err := client.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var r *Response
	err = client.Do(ctx, req, false, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
