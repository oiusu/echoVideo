package main

import (
	"echoVideo/floderServer"
	"flag"
	"fmt"
	"net/http"
)

/**
 * Created by chenc on 2018/9/21
 */

//var mux map[string]func(http.ResponseWriter, *http.Request)
//
//type Myhandler struct{}

func main() {

	//Command line parsing
	bind := flag.String("bind", ":1718", "Bind address")

	flag.Parse()

	http.HandleFunc("/echoVideo", floderServer.HandleSharedFile)
	//http.Handle("/echoVideo", http.HandlerFunc(floderServer.HandleSharedFile))
	//http.Handle("/uploadVideo", http.HandlerFunc(videoToImg.UploadVideo))

	http.ListenAndServe((*bind), nil)
	fmt.Println("server start success ... ")







	//server := http.Server{
	//	Addr:        ":1718",
	//	Handler:     &Myhandler{},
	//	ReadTimeout: 10 * time.Second,
	//}
	//mux = make(map[string]func(http.ResponseWriter, *http.Request))
	//mux["/"] = videoToImg.HandleSharedFile
	//mux["/uploadVideo"] = videoToImg.UploadVideo
	////mux["/video-annotation"] = StaticServer
	//server.ListenAndServe()
	//fmt.Printf("start success \n")

}

//func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if h, ok := mux[r.URL.String()]; ok {
//		h(w, r)
//		return
//	}
//	http.StripPrefix("/", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
//
//}
