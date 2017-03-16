package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"path/filepath"
)

func fileExist(fileName string) bool{
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false;
	}

	return true;
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fileName := req.FormValue("fileName")
	fileContent := req.FormValue("fileContent")
	fmt.Fprint(w, "hello")

	rootPath := req.FormValue("rootPath")

	filePath := fileName
	dir, fileName := filepath.Split(filePath)
	parentDir := filepath.Join(rootPath, dir)
	if fileExist(parentDir){
		log.Println(parentDir + " file exist")
	}else {
		err := os.MkdirAll(parentDir, 755)
		log.Fatalln(err)
	}

	fullFileName := filepath.Join(parentDir, fileName)
	var f *os.File
	if fileExist(fullFileName){
		log.Printf("file exist")
		f, _ = os.OpenFile(fullFileName, os.O_TRUNC|os.O_RDWR , 755)
	}else {
		log.Printf("file not exist")
		f, _ = os.Create(fullFileName)
	}
	defer f.Close()

	f.WriteString(fileContent)

	defer func(){
		log.Println(rootPath + " " + fileName + " sync")
		if e := recover(); e != nil{
			log.Fatal(e)
		}
	}()
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
