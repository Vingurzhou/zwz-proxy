/**
 * Created by zhouwenzhe on 2023/5/25
 */

package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request %s %s %s", request.Method, request.Host, request.RemoteAddr)
	transport := http.DefaultTransport
	outReq := new(http.Request)
	*outReq = *request
	if clientIp, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIp = strings.Join(prior, ",") + "," + clientIp
		}
		outReq.Header.Set("X-Forwarded-For", clientIp)
	}
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		writer.WriteHeader(http.StatusBadGateway)
	}
	for key, value := range res.Header {
		for _, v := range value {
			writer.Header().Add(key, v)
		}
	}
	writer.WriteHeader(res.StatusCode)
	io.Copy(writer, res.Body)
	res.Body.Close()
}

func NewProxy() *Proxy {
	return &Proxy{}
}
