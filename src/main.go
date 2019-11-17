package main

import (
    "sync"
    "time"
    "log"
    "net/http"
    "fmt"
    "net"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
    limiter *limiter
    lastSeen time.Time
}


type limiter struct {
	numOfReq int //number of requests in the bucket
	lastReq time.Time //the time when last request happend
	leak int //handled request per second
	vol int //volume of the bucket
}
// The map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
    go cleanupVisitors()
}

// Create bucket(limiter)
func bucket(leak int, vol int) *limiter {
	l := new(limiter)
	l.numOfReq = 0
	l.lastReq = time.Now();
	l.leak = leak
	l.vol = vol
	return l
}

// Check if the new request is allowed by Leaky bucket algorithm
func (l *limiter) allowed() (flag bool){
	tmp := time.Now().Sub(l.lastReq).Seconds()
	l.numOfReq -= (int(tmp)/l.leak)
	if(l.numOfReq<0){ 
		l.numOfReq = 0
	}

	if(l.numOfReq<l.vol){
		l.numOfReq += 1
		l.lastReq = time.Now()
		return true
	}
	return false
} 

//Add new visitor when get request from new ip
func addVisitor(ip string) *limiter {
    limiter := bucket(1,60)
    mu.Lock()
    // Include the current time when creating a new visitor.
    visitors[ip] = &visitor{limiter, time.Now()}
    mu.Unlock()
    return limiter
}

func getVisitor(ip string) *limiter {
    mu.Lock()

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
            if time.Now().Sub(v.lastSeen) > 1*time.Minute {
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
        if limiter.allowed() == false {
            //http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            w.WriteHeader(http.StatusTooManyRequests)
            fmt.Fprintf(w, "Error, too many request per minutes")
            return
        }

        next.ServeHTTP(w, r)
    })
}


func handler(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            log.Println(err.Error())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
     mu.Lock()

    v, exists := visitors[ip]
    w.WriteHeader(http.StatusOK)
    if !exists {
    	fmt.Fprintf(w, "Number of requests in 1 minutes: %d", 1)     
    }else {
    	fmt.Fprintf(w, "Number of requests in 1 minutes: %d",v.limiter.numOfReq )
    }
    mu.Unlock()
}

func setRoute() http.Handler{
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)
    return mux;	
}

func main() {
	mux := setRoute();
    // Wrap the servemux with the limit middleware.
    log.Println("Listening on :8080...")
    http.ListenAndServe(":8080", limit(mux))
}

