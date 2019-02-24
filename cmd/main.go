/*
 * $env:SLACKIFY_APP_TOKEN="xoxb-206192232134-559147467091-GNnTZczoxXPaN0AVgjhGpQVM"
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/lucasepe/slackify"
)

func main() {

	token := flag.String("token", "", "Authentication token bearing required scopes (required)")
	file := flag.String("file", "", "The file you want to upload (optional)")
	channels := flag.String("channels", "", "Comma-separated list of channel names or IDs where the file will be shared (required)")
	comment := flag.String("comment", "", "A comment introducing the file to the specified channels (optional)")

	flag.Parse()

	if *token == "" {
		*token = os.Getenv("SLACKIFY_APP_TOKEN")
		if *token == "" {
			flag.PrintDefaults()
			os.Exit(1)
			return
		}
	}

	if *channels == "" {
		flag.PrintDefaults()
		os.Exit(1)
		return
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 && *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
		return
	}

	var filetype string
	if len(flag.Args()) > 0 {
		filetype = strings.TrimSpace(flag.Args()[0])
		if filetype == "" {
			filetype = "auto"
		}
	}

	values := map[string]io.Reader{
		"token":    strings.NewReader(*token),
		"channels": strings.NewReader(*channels),
		"filetype": strings.NewReader(filetype),
	}

	if *comment != "" {
		values["initial_comment"] = strings.NewReader(*comment)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		values["content"] = bufio.NewReader(os.Stdin)
	} else {
		r, err := os.Open(*file)
		if err != nil {
			panic(err)
		}
		values["file"] = r
	}

	client := &http.Client{}

	res, err := slackify.Upload(client, slackify.ApiURL, values)
	if err != nil {
		panic(err)
	}

	if !res.Success {
		fmt.Printf("*** Error: %s\n", res.Error)
		return
	}

	fmt.Printf("%s\n", res.File.Permalink)
}
