package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type Prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of GO reverse proxy
	proxy *httputil.ReverseProxy
	// default route path
	defaultPath   *httputil.ReverseProxy
	routePatterns []*regexp.Regexp
}

func New(target string, defaultRoutePath string) *Prox {
	url, _ := url.Parse(target)
	urlDefault, _ := url.Parse(defaultRoutePath)
	fmt.Println("[Transporter] proxy initialized")
	fmt.Println("[Transporter] default path : ", urlDefault)
	fmt.Println("[Transporter] redirection path : ", url)
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url), defaultPath: httputil.NewSingleHostReverseProxy(urlDefault)}
}

func (p *Prox) parseWhiteList(r *http.Request) bool {
	for _, regexp := range p.routePatterns {
		if regexp.MatchString(r.URL.Path) {
			return true
		}
	}
	return false
}

func rewriteURL(r *http.Request, subpath string) {
	switch subpath {
	case "<ROUTE_URI>":
		r.RequestURI = ""
		r.URL, _ = url.Parse("<NEW_URL>")
	default:
		fmt.Println("No change required baby!")
	}
}

func (p *Prox) getSubpath(srcURL string) string {
	for _, regexp := range p.routePatterns {
		if regexp.MatchString(srcURL) {
			return regexp.String()
		}
	}
	return srcURL
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")

	fmt.Println("[Transporter] a request reached")
	fmt.Println("[Transporter] request : ", r)
	fmt.Println("[Transporter] evaluating redirection path")

	if p.routePatterns == nil || p.parseWhiteList(r) {
		// url rewrite
		subpath := p.getSubpath(r.URL.String())
		fmt.Println("[Transporter] " + subpath + " subpath detected")
		rewriteURL(r, subpath)
		r.Host = r.URL.Host
		fmt.Println("[Transporter] rerouting request : ", r)
		// end url rewrite
		p.proxy.ServeHTTP(w, r)
	} else {
		fmt.Println("[Transporter] default path detected")
		http.DefaultServeMux.ServeHTTP(w, r)
	}
}

func main() {
	const (
		defaultPort             = ":8001"
		defaultPortUsage        = "default server port, ':8001', ':8080'..."
		defaultTarget           = "<TARGET_URL>"
		defaultTargetUsage      = "default redirect url, '<TARGET_URL>'"
		defaultServer           = "<DEFAULT_URL>"
		defaultServerUsage      = "default server is at, '<DEFAULT_URL>'"
		defaultWhiteRoutes      = `<ROUTE_URI>`
		defaultWhiteRoutesUsage = "list of white route as regexp, '/path1*,/path2*...."
	)

	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)
	defaultServerURL := flag.String("*defaultServer", defaultServer, defaultServerUsage)
	routesRegexp := flag.String("routes", defaultWhiteRoutes, defaultWhiteRoutesUsage)

	flag.Parse()

	fmt.Println("[Transporter] server will run on : ", *port)
	fmt.Println("[Transporter] Default server : ", *defaultServerURL)
	fmt.Println("[Transporter] redirecting to : ", *url)
	fmt.Println("[Transporter] accepted routes : ", *routesRegexp)

	reg, _ := regexp.Compile(*routesRegexp)
	routes := []*regexp.Regexp{reg}

	// proxy
	proxy := New(*url, *defaultServerURL)
	proxy.routePatterns = routes

	// server
	http.HandleFunc("/", proxy.handle)
	fmt.Println("[Transporter] handlers registered to : " + *port + "/")
	http.ListenAndServe(*port, nil)
}
