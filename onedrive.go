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
	RequestUA   string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Safari/537.36 Edg/89.0.774.45"
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
	log.Println("original url: ", o.URL)

	// step 1
	// get temp url
	req1, _ := http.NewRequest("HEAD", o.URL, nil)
	req1.Header.Add("User-Agent", RequestUA)
	res1, err := o.Client.Do(req1)
	if err != nil {
		log.Println("res1 URL: ", err)
		return err
	}
	log.Println(res1.Location())

	// step 1
	// get temp url
	r1, err := o.Client.Head(o.URL)
	if err != nil {
		log.Println("r1: ", err)
		return err
	}
	log.Println(r1.Header)
	tempURL, err := res1.Location()
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
	req2, _ := http.NewRequest("HEAD", downloadURL, nil)
	req2.Header.Add("User-Agent", RequestUA)
	r2, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Println("r2: ", err)
		return err
	}

	directURL, err := r2.Location()
	if err != nil {
		log.Println("directURL error: ", err)
		return err
	}

	o.DirectURL = directURL.String()
	log.Println(o.DirectURL)

	return nil
}
