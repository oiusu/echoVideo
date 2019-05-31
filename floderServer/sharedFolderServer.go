/* Tiny web server in Golang for sharing a folder
Copyright (c) 2010-2014 Alexis ROBERT <alexis.robert@gmail.com>

Contains some code from Golang's http.ServeFile method, and
uses lighttpd's directory listing HTML template. */

package floderServer

import (
	"bytes"
	"container/list"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

const (
	serverUA      = "echoVideo/0.0.1"
	fs_maxbufsize = 4096 // 4096 bits = default page size on OSX
)

/* Go is the first programming language with a templating engine embeddeed11
 * but with no min function. */
func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Manages directory listings
type dirlisting struct {
	Name           string
	RootPath       string
	ParentPath     string
	Children_dir   []string
	Children_files []string
	ServerUA       string
}

func copyToArray(src *list.List) []string {
	dst := make([]string, src.Len())

	i := 0
	for e := src.Front(); e != nil; e = e.Next() {
		dst[i] = e.Value.(string)
		i = i + 1
	}

	return dst
}

func HandleSharedFile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", serverUA)

	//serveFile(filepath, w, req)
	serverVideo(w, req)

	fmt.Printf("\"%s %s %s\" \"%s\" \"%s\"\n",
		req.Method,
		req.URL.Path,
		req.Proto,
		req.Referer(),
		req.UserAgent()) // TODO: Improve this crappy logging
}

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request, trimReqPath string) {
	names, _ := f.Readdir(-1)

	// Otherwise, generate folder content.
	children_dir_tmp := list.New()
	children_files_tmp := list.New()

	for _, val := range names {
		if val.Name()[0] == '.' {
			continue
		} // Remove hidden files from listing

		if val.IsDir() {
			children_dir_tmp.PushBack(val.Name())
		} else {
			children_files_tmp.PushBack(val.Name())
		}
	}

	// And transfer the content to the final array structure
	children_dir := copyToArray(children_dir_tmp)
	children_files := copyToArray(children_files_tmp)
	// 获得父路径11
	var parent_path string
	path_arr := strings.Split(trimReqPath, "/")

	var buffer bytes.Buffer
	buffer.WriteString(*Url_prefix)
	if len(path_arr) > 2 {
		subPath := path_arr[:len(path_arr)-1]
		for _, val := range subPath {
			if val != "" {
				buffer.WriteString("/")
				buffer.WriteString(val)
			}
		}
	}
	buffer.WriteString("/")
	parent_path = buffer.String() // 拼接结果
	//tpl, err := template.New("tpl").Parse(dirlisting_tpl)
	tpl, err := template.ParseFiles(*template_dir + "tpl.html")

	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		return
	}

	data := dirlisting{Name: trimReqPath, RootPath: (*root_folder), ParentPath: parent_path, ServerUA: serverUA,
		Children_dir: children_dir, Children_files: children_files}
	fmt.Printf("parent_path = %s \n", parent_path)
	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

//func getPathSuffix(path string) string{
//	path
//}

func serverVideo(w http.ResponseWriter, req *http.Request) {

	reqPath := path.Clean(req.URL.Path)
	var trimReqPath string
	if reqPath == *Url_prefix {
		trimReqPath = "/"
	} else {
		trimReqPath = strings.Replace(reqPath, *Url_prefix+"/", "/", 1) //todo 手动替换前缀
	}

	// Opening the file handle
	filepath := path.Join((*root_folder), trimReqPath)
	f, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		return
	}

	defer f.Close()

	// Checking if the opened handle is really a file
	statinfo, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		handleDirectory(f, w, req, trimReqPath)
		return
	}

	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
		return
	}
	// Content-Type handling
	query, err := url.ParseQuery(req.URL.RawQuery)

	if err == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		// Fetching file's mimetype and giving it to the browser
		ext := path.Ext(filepath)
		fmt.Printf("ext=%s \n", ext)
		if mimetype := mime.TypeByExtension(ext); mimetype != "" {
			if ext == ".avi" || ext == ".MP4" {
				//todo .avi 需要手动指定
				w.Header().Set("Content-Type", "video/mp4")
				//fmt.Printf("set Content-Type video/mp4 , req.URL.Path=%s \n",req.URL.Path)
			} else {
				w.Header().Set("Content-Type", mimetype)
				//fmt.Printf("set Content-Type %s , req.URL.Path=%s \n",mimetype,req.URL.Path)
			}
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
			//fmt.Printf("set Content-Type application/octet-stream ,req.URL.Path=%s \n",req.URL.Path)
		}
	}

	http.ServeContent(w, req, filepath, time.Now(), f)

}

//func serveFile(filepath string, w http.ResponseWriter, req *http.Request) {
//	// Opening the file handle
//	f, err := os.Open(filepath)
//	if err != nil {
//		http.Error(w, "404 Not Found : Error while opening the file.", 404)
//		return
//	}
//
//	defer f.Close()
//
//	// Checking if the opened handle is really a file
//	statinfo, err := f.Stat()
//	if err != nil {
//		http.Error(w, "500 Internal Error : stat() failure.", 500)
//		return
//	}
//
//	if statinfo.IsDir() { // If it's a directory, open it !
//		handleDirectory(f, w, req)
//		return
//	}
//
//	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
//		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
//		return
//	}
//
//	// Manages If-Modified-Since and add Last-Modified (taken from Golang code)
//	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
//		w.WriteHeader(http.StatusNotModified)
//		return
//	}
//	w.Header().Set("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))
//
//	// Content-Type handling
//	query, err := url.ParseQuery(req.URL.RawQuery)
//
//	if err == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
//		w.Header().Set("Content-Type", "application/octet-stream")
//	} else {
//		// Fetching file's mimetype and giving it to the browser
//		if mimetype := mime.TypeByExtension(path.Ext(filepath)); mimetype != "" {
//			w.Header().Set("Content-Type", mimetype)
//		} else {
//			w.Header().Set("Content-Type", "application/octet-stream")
//		}
//	}
//
//	// Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
//	if req.Header.Get("Range") != "" {
//		start_byte := parseRange(req.Header.Get("Range"))
//
//		if start_byte < statinfo.Size() {
//			f.Seek(start_byte, 0)
//		} else {
//			start_byte = 0
//		}
//
//		w.Header().Set("Content-Range",
//			fmt.Sprintf("bytes %d-%d/%d", start_byte, statinfo.Size()-1, statinfo.Size()))
//	}
//
//	// Manage gzip/zlib compression
//	output_writer := w.(io.Writer)
//
//	is_compressed_reply := false
//
//	if (*uses_gzip) == true && req.Header.Get("Accept-Encoding") != "" {
//		encodings := parseCSV(req.Header.Get("Accept-Encoding"))
//
//		for _, val := range encodings {
//			if val == "gzip" {
//				w.Header().Set("Content-Encoding", "gzip")
//				output_writer = gzip.NewWriter(w)
//
//				is_compressed_reply = true
//
//				break
//			} else if val == "deflate" {
//				w.Header().Set("Content-Encoding", "deflate")
//				output_writer = zlib.NewWriter(w)
//
//				is_compressed_reply = true
//
//				break
//			}
//		}
//	}
//
//	if !is_compressed_reply {
//		// Add Content-Length
//		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))
//	}
//
//	// Stream data out !
//	buf := make([]byte, min(fs_maxbufsize, statinfo.Size()))
//	n := 0
//	for err == nil {
//		n, err = f.Read(buf)
//		output_writer.Write(buf[0:n])
//	}
//
//	// Closes current compressors
//	switch output_writer.(type) {
//	case *gzip.Writer:
//		output_writer.(*gzip.Writer).Close()
//	case *zlib.Writer:
//		output_writer.(*zlib.Writer).Close()
//	}
//
//	f.Close()
//}

//func parseCSV(data string) []string {
//	splitted := strings.SplitN(data, ",", -1)
//
//	data_tmp := make([]string, len(splitted))
//
//	for i, val := range splitted {
//		data_tmp[i] = strings.TrimSpace(val)
//	}
//
//	return data_tmp
//}

//func parseRange(data string) int64 {
//	stop := (int64)(0)
//	part := 0
//	for i := 0; i < len(data) && part < 2; i = i + 1 {
//		if part == 0 { // part = 0 <=> equal isn't met.
//			if data[i] == '=' {
//				part = 1
//			}
//
//			continue
//		}
//
//		if part == 1 { // part = 1 <=> we've met the equal, parse beginning
//			if data[i] == ',' || data[i] == '-' {
//				part = 2 // part = 2 <=> OK DUDE.
//			} else {
//				if 48 <= data[i] && data[i] <= 57 { // If it's a digit ...
//					// ... convert the char to integer and add it!
//					stop = (stop * 10) + (((int64)(data[i])) - 48)
//				} else {
//					part = 2 // Parsing error! No error needed : 0 = from start.
//				}
//			}
//		}
//	}
//
//	return stop
//}
