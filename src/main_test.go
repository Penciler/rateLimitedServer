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


func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
   req, _ := http.NewRequest(method, path, nil)
   w := httptest.NewRecorder()
   r.ServeHTTP(w, req)
   return w
}

func TestMain( t *testing.T){
	mux := setRoute()
	//http.ListenAndServe(":8080", limit(mux))
	svr := httptest.NewServer(limit(mux))

	for i := 0;i<10;i++ {
		res, err := http.Get(svr.URL)
		if err != nil {
			log.Fatal(err)
		}
		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", greeting)
		time.Sleep(10*time.Millisecond);
	}

	// res2, err2 := http.Get(svr.URL)
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// greeting2, err2 := ioutil.ReadAll(res2.Body)
	// res2.Body.Close()
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	// fmt.Printf("%s", greeting2)

	// res3, err3 := http.Get(svr.URL)
	// if err3 != nil {
	// 	log.Fatal(err3)
	// }
	// greeting3, err3 := ioutil.ReadAll(res3.Body)
	// res3.Body.Close()
	// if err3 != nil {
	// 	log.Fatal(err3)
	// }

	// fmt.Printf("%s", greeting3)

	// res24, err24 := http.Get(svr.URL)
	// if err24 != nil {
	// 	log.Fatal(err24)
	// }
	// greeting24, err24 := ioutil.ReadAll(res24.Body)
	// res24.Body.Close()
	// if err24 != nil {
	// 	log.Fatal(err24)
	// }

	// fmt.Printf("%s", greeting24)
	// w1 := performRequest(mux, "GET", "/")
	// w2 := performRequest(mux, "GET", "/")
	// status1 := w1.Code
	// log.Println(status1)
	// body1, _ := ioutil.ReadAll(w1.Result().Body)
	// fmt.Println(string(body1))

	// status2 := w2.Code
	// log.Println(status2)
	// body2, _ := ioutil.ReadAll(w2.Result().Body)
	// fmt.Println(string(body2))
}