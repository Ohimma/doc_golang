package main

import (
	"fmt"
	"net/http"
	"time"
	"log"
	"flag"
	"os"
	"io"
	"strings"
	"strconv"
	"io/ioutil"
	"html/template"
_	"mime"
_	"path/filepath"
_	"crypto/rand"
)

const maxUploadSize = 20 * 1024 * 1024 // 2 mb

type HelloHandler struct{}
func (h HelloHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello Handler!")
}

func hello2 (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello HandlerFunc")
	// n, _ := io.WriteString(w, "Hello HandlerFunc")
}

func uploadFileHandler(dir string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, nil)
			return
		}
		
		// 解析 request 请求头 
		fmt.Println("request method: ", r.Method)
		// RequestURI是被客户端发送到服务端的请求的请求行中未修改的请求URI
		fmt.Println("request RequestURI: ", r.RequestURI)
		//URL类型,下方分别列出URL的各成员
		fmt.Println("request URL_scheme: ", r.URL.Scheme)
		fmt.Println("request URL_opaque: ", r.URL.Opaque)
		fmt.Println("request URL_user: ", r.URL.User.String())
		fmt.Println("request URL_host: ", r.URL.Host)
		fmt.Println("request URL_path: ", r.URL.Path)
		fmt.Println("request URL_RawQuery: ", r.URL.RawQuery)
		fmt.Println("request URL_Fragment: ", r.URL.Fragment)
		//协议版本
		fmt.Println("request proto: ", r.Proto)
		fmt.Println("request protomajor: ", r.ProtoMajor)
		fmt.Println("request protominor: ", r.ProtoMinor)
		fmt.Println("request RemoteAddr: ", r.RemoteAddr)
		fmt.Println("request ContentLength: ", r.ContentLength)
		fmt.Println("request Close: ", r.Close)
		fmt.Println("request host: ", r.Host)

		//HTTP请求的头域
		fmt.Println("request header: ", r.Header)
		for k, v := range r.Header {
			// fmt.Println("Header key:" + k)
			for _, vv := range v {
				fmt.Println("request header key: " + k + ": " + vv)
			}
		}
		//判断是否multipart方式
		is_multipart := false
		for _, v := range r.Header["Content-Type"] {
			if strings.Index(v, "multipart/form-data") != -1 {
				is_multipart = true
			}
		}

		//解析body
		if is_multipart == true {
			err := r.ParseMultipartForm(maxUploadSize)
			if err != nil {
				fmt.Printf("Could not parse multipart form: %v\n", err)
				renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
				return
			}

			// ParseMultipartForm 将请求的主体作为multipart/form-data解析
			fmt.Println("解析方式:ParseMultipartForm = ", r.MultipartForm)
			files := r.MultipartForm.File
			for k, v := range files {
				fmt.Println("tmp1 = ",k, v)
				for kk, vv := range v {
					fmt.Println("tmp2 = ",kk, vv)
					fmt.Println("request ParseMultipartForm: " + vv.Filename)
				}
			}
		} else {
			// ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段
			r.ParseForm()
			fmt.Println("解析方式:ParseForm")
			// FormValue返回key为键查询r.Form字段得到结果[]string切片的第一个值
			fmt.Println("Form", r.Form)
			// PostFormValue返回key为键查询r.PostForm字段得到结果[]string切片的第一个值
			fmt.Println("PostForm", r.PostForm)
		}
		
		fmt.Println("====================================")
		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fmt.Println("FormFile head: ", fileHeader)
		fmt.Println("FormFile file: ", file)
		fmt.Println("FormFile size: ", fileHeader.Size)
		fmt.Println("FormFile name: ", fileHeader.Filename)

		if fileHeader.Size > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		fmt.Println("====================================")
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// 创建文件 
		filename := dir + fileHeader.Filename
		newFile, err := os.Create(filename)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		// write file
		defer newFile.Close() // idempotent, okay to call twice
		_, err = newFile.Write(fileBytes)
		if err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}


func FileDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("download url=%s (%v) [%v]\n", r.RequestURI, r.RequestURI[10:], strings.Split(r.RequestURI, "/download/")[1:])
	filename := r.RequestURI[10:]

	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println("os.Stat err =",err)
		return
	}
	fmt.Println("name =",info.Name())
	fmt.Println("size =",info.Size())
	fmt.Println("mode =",info.Mode())
	fmt.Println("modtime =",info.ModTime())
	fmt.Println("isDir =",info.IsDir())
	fmt.Println("sys =",info.Sys())

	if info.IsDir() {
		fmt.Println("is dir")
		// a := http.FileServer(http.Dir("/Users/admin/me/git/Otools/golang/test"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fileInfos,err:=ioutil.ReadDir(filename)
		if err!=nil{
			fmt.Println("ReadDir failed,error:%v\n", err)
			return
		}
		fmt.Println("a=", fileInfos)
	} else {
		fmt.Println("is file")
		// 1. 一次性读取
		// content, err := ioutil.ReadFile(filename)
		// if err != nil {
		// 	fmt.Println("read file err ", err)
		// }
		// fmt.Printf("One File contents: %v\n", content)
		// fmt.Printf("Two File contents: %v\n", string(content))
		// fmt.Printf("Three File contents: %s\n", content
	

		// 2. 打开文件
		file, err := os.Open(filename)
		if err != nil {
		    fmt.Println("open file err", err)
		}
		defer file.Close()

		// 3. 分块读
		fileHeader := make([]byte, 4*1024)
		file.Read(fileHeader)
	
		fileStat, _ := file.Stat()
	
		w.Header().Set("Content-Disposition", "attachment; filename=" + filename)
		w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
		w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	
		file.Seek(0, 0)
		io.Copy(w, file)

		fmt.Printf("file = %v\n", w)
		return
	}
}

func main() {

	a, _ := os.Getwd()

	var dir string
	var port int
    flag.StringVar(&dir, "dir", a, "下载路径")
    flag.IntVar(&port, "port", 8079, "监听端口")
    flag.Parse()

	fmt.Printf("start http, 下载页面为 = %v \n", dir)
    
	s := &http.Server {
		Addr: fmt.Sprintf(":%d", port),
	    ReadTimeout:    10 * time.Second,
	    WriteTimeout:   10 * time.Second,
	    MaxHeaderBytes: 1 << 20,
	}

	helloHandler := HelloHandler{}
    http.Handle("/hello1", helloHandler)
	http.HandleFunc("/hello2", hello2)
	http.HandleFunc("/quick/", FileDownload)

	http.HandleFunc("/upload", uploadFileHandler(dir))

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/download/", http.StripPrefix("/download/", fs))
	
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("LsitenServer = ", err)
	}
	log.Print("ListenServer  = ")
}