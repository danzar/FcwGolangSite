package resource

import (
	"net/http"
	"strings"
	"os"
	"bufio"
	"github.com/danzar/fcwServer/common"
)

var debug = false

var themeName = "bs3"

/*
 We are serving up the resource type files here
 So that when a page is loaded and it has some resource type
 this func gets the type and loads it from the passed request type.
*/
func ServerResourceFiles(w http.ResponseWriter, req *http.Request)   {
	path := "public/" + themeName + req.URL.Path
	//log.Println("ResourceRequest:" + path)
	var contentType string

	if strings.HasSuffix(path,".css"){
		contentType = "text/css; charset=utf-8"
	}else if strings.HasSuffix(path,".png"){
		contentType = "image/png; charset=utf-8"
	}else if strings.HasSuffix(path,".jpg"){
		contentType = "image/jpg; charset=utf-8"
	}else if strings.HasSuffix(path,".js"){
		contentType = "application/javascript; charset=utf-8"
	}else {
		contentType = "text/plain; charset=utf-8"
	}



	common.LogDebugData(path, debug)
	f, err := os.Open(path)
	if err == nil{
		defer f.Close()
		w.Header().Add("Content-Type",contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	}else{
		w.WriteHeader(404)
	}
}