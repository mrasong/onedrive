# OneDrive Direct Link


1. Get you 1drv.ms link url like this `https://1drv.ms/{action}/s!{token}`
  eg: `https://1drv.ms/u/s!Al8SermeNqO5gi1negxQ6wUzyB36?e=3s4Z2p`
  we got `action` =>  `u`, `token` => Al8SermeNqO5gi1negxQ6wUzyB36

2. use as a lib.

```go
package main

import (
	"log"

	"github.com/mrasong/onedrive"
)

func main() {
    // https://1drv.ms/{action}/s!{token}
    onedriveURL := "https://1drv.ms/u/s!Al8SermeNqO5gi1negxQ6wUzyB36?e=3s4Z2p"
    action := "u"
    token := "Al8SermeNqO5gi1negxQ6wUzyB36"

    // get direct link by action and token
    url := onedrive.New(action, token).GetDirectLink()
	log.Println(url)
    
    // get direct link by URL
    url := onedrive.NewFromURL(onedriveURL).GetDirectLink()
	log.Println(url)
}
```

3. use as a server

```bash
go get -u github.com/mrasong/onedrive/onedrive
```

```bash
onedrive serve
```

