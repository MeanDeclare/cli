package debug

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/exercism/cli/utils"
)

var (
	// Verbose determines if debugging output is displayed to the user
	Verbose      bool
	output       io.Writer = os.Stderr
	UnmaskAPIKey bool
)

// Println conditionally outputs a message to Stderr
func Println(args ...interface{}) {
	if Verbose {
		fmt.Fprintln(output, args...)
	}
}

// Printf conditionally outputs a formatted message to Stderr
func Printf(format string, args ...interface{}) {
	if Verbose {
		fmt.Fprintf(output, format, args...)
	}
}

// DumpRequest dumps out the provided http.Request
func DumpRequest(req *http.Request) {
	if !Verbose {
		return
	}

	var bodyCopy bytes.Buffer
	body := io.TeeReader(req.Body, &bodyCopy)
	req.Body = ioutil.NopCloser(body)

	temp := req.Header.Get("Authorization")

	if !UnmaskAPIKey {
		req.Header.Set("Authorization", "Bearer "+utils.Redact(strings.Split(temp, " ")[1]))
	}

	dump, err := httputil.DumpRequest(req, req.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpRequest =========================")
	Println(string(dump))
	Println("========================= END DumpRequest =========================")
	Println("")

	req.Header.Set("Authorization", temp)
	req.Body = ioutil.NopCloser(&bodyCopy)
}

// DumpResponse dumps out the provided http.Response
func DumpResponse(res *http.Response) {
	if !Verbose {
		return
	}

	var bodyCopy bytes.Buffer
	body := io.TeeReader(res.Body, &bodyCopy)
	res.Body = ioutil.NopCloser(body)

	dump, err := httputil.DumpResponse(res, res.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpResponse =========================")
	Println(string(dump))
	Println("========================= END DumpResponse =========================")
	Println("")

	res.Body = ioutil.NopCloser(body)
}
