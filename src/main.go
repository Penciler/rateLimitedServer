package main

import (
    "sync"
    "time"
    "log"
    "net/http"
    "fmt"
    "net"
    //"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
    //limiter  *rate.Limiter
    limiter *limiter
    lastSeen time.Time
}


type limiter struct {
	numOfReq int
	lastReq time.Time
	leak int
	vol int
}
// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
    go cleanupVisitors()
}

func bucket(leak int, vol int) *limiter {
	l := new(limiter)
	l.numOfReq = 0
	l.lastReq = time.Now();
	l.leak = leak
	l.vol = vol
	return l
}



func (l *limiter) allowed() (flag bool){
	tmp := time.Now().Sub(l.lastReq).Seconds()
	log.Println(int(tmp))
	l.numOfReq -= (int(tmp)/l.leak)
	if(l.numOfReq<0){ 
		l.numOfReq = 0
	}

	log.Println(l.numOfReq)
	if(l.numOfReq<l.vol){
		l.numOfReq += 1
		l.lastReq = time.Now()
		log.Println(l.numOfReq)
		return true
	}
	return false
} 

//func addVisitor(ip string) *rate.Limiter {
func addVisitor(ip string) *limiter {
    //limiter := rate.NewLimiter(1, 2)
    limiter := bucket(1,20)
    mu.Lock()
    // Include the current time when creating a new visitor.
    visitors[ip] = &visitor{limiter, time.Now()}
    mu.Unlock()
    return limiter
}

//func getVisitor(ip string) *rate.Limiter {
func getVisitor(ip string) *limiter {
    mu.Lock()
    //defer mu.Unlock()

    v, exists := visitors[ip]
    if !exists {
    	mu.Unlock()
        return addVisitor(ip)
    }

    // Update the last seen time for the visitor.
    v.lastSeen = time.Now()
    mu.Unlock()
    return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
    for {
        time.Sleep(time.Minute)

        mu.Lock()
        defer mu.Unlock()
        for ip, v := range visitors {
            if time.Now().Sub(v.lastSeen) > 3*time.Minute {
                delete(visitors, ip)
            }
        }
    }
}

func limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            log.Println(err.Error())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        log.Println(ip)
        limiter := getVisitor(ip)
        //limiter := addVisitor(ip)
        //if limiter.Allow() == false {
        if limiter.allowed() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}


//var limiter = rate.NewLimiter(1, 3)

// func limit(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         if limiter.Allow() == false {
//             http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
//             return
//         }

//         next.ServeHTTP(w, r)
//     })
// }


func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            log.Println(err.Error())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
     mu.Lock()
    //defer mu.Unlock()

    v, exists := visitors[ip]
    if !exists {
    	fmt.Fprintf(w, "Number of requests in 1 minutes %d!", 1)     
    }else {
    	fmt.Fprintf(w, "Number of requests in 1 minutes %d!",v.limiter.numOfReq )
    }

    mu.Unlock()
}

func setRoute() http.Handler{
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)
    return mux;	
}

func main() {
    //http.HandleFunc("/", handler)
    //log.Fatal(http.ListenAndServe(":8080", nil))
    //mux := http.NewServeMux()
    //mux.HandleFunc("/", handler)
	mux := setRoute();
    // Wrap the servemux with the limit middleware.
    log.Println("Listening on :8080...")
    http.ListenAndServe(":8080", limit(mux))
}

