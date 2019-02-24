package slackify

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const ApiURL = "https://slack.com/api/files.upload"

type UploadResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	File    struct {
		ID              string `json:"id,omitempty"`
		User            string `json:"user,omitempty"`
		Permalink       string `json:"permalink,omitempty"`
		PermalinkPublic string `json:"permalink_public,omitempty"`
	} `json:"file,omitempty"`
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (result UploadResponse, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				result.Error = err.Error()
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				result.Error = err.Error()
				return
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			result.Error = err.Error()
			return result, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		result.Error = err.Error()
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	//req.Header.Set("Authorization", "Bearer "+token)

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		result.Error = err.Error()
		return
	}

	// Check the response
	/*
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", res.Status)
		}
	*/

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		result.Error = err.Error()
	}

	return
}
