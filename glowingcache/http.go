package glowingcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	if !strings.HasPrefix(request.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected request. Path " + request.URL.Path)
	}
	p.Log("%s %s", request.Method, request.URL.Path)
	// /<basepath>/<groupname>/<key>
	parts := strings.SplitN(request.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(writer, "Group Not Found", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Write(view.ByteSlice())

}
