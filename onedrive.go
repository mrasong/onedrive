package onedrive

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// OneDriveDirectLink const
const (
	URLTemplate string = "https://1drv.ms/%s/s!%s"
)

// OneDrive struct
type OneDrive struct {
	URL       string // https://1drv.ms/{:action}/s!{:token}
	DirectURL string

	Client *http.Client
}

// New return *OneDrive
func New(action, token string) *OneDrive {
	return NewFromURL(fmt.Sprintf(URLTemplate, action, token))
}

// NewFromURL return *OneDrive
func NewFromURL(url string) *OneDrive {
	return &OneDrive{
		URL: url,
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.RequestURI)
				return http.ErrUseLastResponse
			},
		},
	}
}

// GetDirectLink return direct link url
func (o *OneDrive) GetDirectLink() string {
	if err := o.handler(); err != nil {
		log.Println(err)
	}

	return o.DirectURL
}

// handler url parse
func (o *OneDrive) handler() error {
	request := func(method, url string) *http.Request {
		req, _ := http.NewRequest(method, url, nil)
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7,ja;q=0.6")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36")
		return req
	}

	log.Println("original url: ", o.URL)

	// step 1
	// get temp url
	r1, err := o.Client.Do(request("HEAD", o.URL))
	if err != nil {
		log.Println("r1: ", err)
		return err
	}
	defer r1.Body.Close()
	log.Println(r1.Header)
	tempURL, err := r1.Location()
	if err != nil {
		log.Println("tempURL error: ", err)
		return err
	}
	log.Println("tempURL: ", tempURL)

	// step 2
	// replace tempURL to download URL
	downloadURL := strings.ReplaceAll(tempURL.String(), "/redir?", "/download?")
	log.Println("downloadURL: ", downloadURL)

	// step 3
	// get direct URL
	r2, err := o.Client.Do(request("HEAD", downloadURL))
	if err != nil {
		log.Println("r2: ", err)
		return err
	}
	defer r2.Body.Close()

	directURL, err := r2.Location()
	if err != nil {
		log.Println("directURL error: ", err)
		return err
	}

	o.DirectURL = directURL.String()
	log.Println(o.DirectURL)

	return nil
}
