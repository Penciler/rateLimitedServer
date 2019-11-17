package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"log"
	"fmt"
	"io/ioutil"
	"time"
)


func performRequest(r http.Handler, method, path string, srcIp string) *httptest.ResponseRecorder {
   req, _ := http.NewRequest(method, path, nil)
   req.RemoteAddr = srcIp
   w := httptest.NewRecorder()
   r.ServeHTTP(w, req)
   return w
}

func TestMain( t *testing.T){
	mux := setRoute()
	//svr := httptest.NewServer(limit(mux))
	r := limit(mux)

	//Test with ok loading
	for i := 0;i<60;i++ {
		w := performRequest(r, "GET", "/", "192.168.12.1:8080")
   		res := w.Result()

		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", greeting)
		if(w.Code != 200){
			t.Errorf("Server returned wrong status code: got %d want %d", w.Code, 200)
		}
		time.Sleep(10*time.Millisecond);
	}

	//Test with over loading
	for i := 0;i<61;i++ {
		w := performRequest(r, "GET", "/", "192.168.12.2:8080")
   		res := w.Result()

		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", greeting)

		if(i == 61 && w.Code != 429){
			t.Errorf("Server returned wrong status code: got %d want %d", w.Code, 429)
		}
		time.Sleep(10*time.Millisecond);
	}
}