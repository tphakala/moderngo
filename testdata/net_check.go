package testdata

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// --- DeprecatedReverseProxyDirector ---

func checkReverseProxyDirector() {
	// Should trigger: Director field in struct literal
	_ = &httputil.ReverseProxy{ // want: "ReverseProxy.Director is deprecated"
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
		},
	}

	// Should trigger: Director assignment
	proxy := &httputil.ReverseProxy{}
	proxy.Director = func(req *http.Request) { // want: "ReverseProxy.Director is deprecated"
		req.URL.Scheme = "https"
	}
	_ = proxy

	// Should NOT trigger: Rewrite field (the recommended replacement)
	_ = &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
		},
	}

	// FALSE POSITIVE: $proxy.Director = $_ matches ANY .Director assignment,
	// not just httputil.ReverseProxy. This fires on unrelated types.
	type Movie struct{ Director string }
	m := Movie{}
	m.Director = "Spielberg" // Should NOT trigger: not a ReverseProxy
}

// --- FilepathIsLocal (false-positive-prone) ---

func checkFilepathIsLocal(path string) {
	// Should trigger: simple .. check on a path variable
	if strings.Contains(path, "..") { // want: "consider using filepath.IsLocal"
		return
	}

	// FALSE POSITIVE RISK: strings.Contains with ".." in non-path contexts
	// These WILL trigger because the rule matches any strings.Contains($, "..")
	// but they are not about path traversal:

	// Checking version ranges (e.g. "1.0..2.0")
	version := "1.0..2.0"
	if strings.Contains(version, "..") { // want: "consider using filepath.IsLocal"
		_ = version
	}

	// Checking for ellipsis in text
	text := "loading..."
	if strings.Contains(text, "..") { // want: "consider using filepath.IsLocal"
		_ = text
	}

	// Should NOT trigger: different string check
	if strings.Contains(path, "/") {
		return
	}

	// Should NOT trigger: HasPrefix (not Contains)
	if strings.HasPrefix(path, "..") {
		return
	}
}

// --- JoinHostPort (false-positive-prone) ---

func checkJoinHostPort(host string, port int) {
	// Should trigger: fmt.Sprintf for host:port with integer port
	_ = fmt.Sprintf("%s:%d", host, port) // want: "use net.JoinHostPort"

	// Should trigger: %v variant
	_ = fmt.Sprintf("%v:%d", host, port) // want: "use net.JoinHostPort"

	// Should NOT trigger: string port (could be non-network)
	_ = fmt.Sprintf("%s:%s", host, "8080")

	// Should NOT trigger: different format entirely
	_ = fmt.Sprintf("ratio: %d:%d", 16, 9)

	// Should NOT trigger: more than two args
	_ = fmt.Sprintf("%s:%d (attempt %d)", host, port, 3)
}

// --- ErrorBeforeUse ---

func checkErrorBeforeUse() {
	// NOTE: ErrorBeforeUse doesn't fire — ruleguard multi-statement matching limitation
	f, err := os.Open("test.txt")
	_ = f.Name()
	if err != nil {
		return
	}
	_ = f
}
