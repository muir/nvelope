package nvelope_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/muir/nchi"
	"github.com/muir/nvelope"
)

// ExampleMinimalErrorHandler shows how
// nape, and nvelope.
func ExampleMinimalErrorHandler() {
	mux := nchi.NewRouter()
	mux.Use(nvelope.MinimalErrorHandler)
	mux.Get("/example/:param", func(w http.ResponseWriter, params nchi.Params) error {
		value := params.ByName("param")
		if value == "good" {
			_, _ = w.Write([]byte("okay"))
			return nil
		}
		return fmt.Errorf("ooh, %s", value)
	})
	ts := httptest.NewServer(mux)
	client := ts.Client()
	doGet := func(url string, authHeader string) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", ts.URL+url, nil)
		if err != nil {
			fmt.Println("request error:", err)
			return
		}
		req.Header.Set("Authentication", authHeader)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("response error:", err)
			return
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		res.Body.Close()
		fmt.Println(res.StatusCode, "->"+string(b))
	}
	doGet("/example/good", "good")
	doGet("/example/bad", "bad")
	// Output: 200 ->okay
	// 500 ->ooh, bad
}
