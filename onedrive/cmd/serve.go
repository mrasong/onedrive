// Copyright © 2021 mrasong <i@mrasong.com>

package cmd

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrasong/onedrive"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		listen := ":8009"
		if len(args) > 0 {
			listen = ":" + args[0]
		}
		serve(listen)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(listen string) {
	log.Printf("ListenAndServe %s", listen)
	r := mux.NewRouter()
	r.HandleFunc("/{action}/s!{token}", serveRedirectHandler)
	log.Fatal(http.ListenAndServe(listen, r))
}

// serveRedirectHandler handler
func serveRedirectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	url := onedrive.New(params["action"], params["token"]).GetDirectLink()
	http.Redirect(w, r, url, 302)
	return
}
