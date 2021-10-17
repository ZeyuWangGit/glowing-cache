package glowingcache

import (
	"fmt"
	"log"
	"net/http"
)

const defaultBasePath = "/_glowingcache/"

type HTTPPool struct {
	selfPath string
	basePath string
}

func NewHttpPool(path string) *HTTPPool {
	return &HTTPPool{
		selfPath: path,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{})  {
	log.Printf("[Servier %s] %s", p.selfPath, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {

}
