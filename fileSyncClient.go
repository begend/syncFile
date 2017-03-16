package main
import(
    	"fmt"
	"github.com/rjeczalik/notify"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"github.com/olebedev/config"
	"strings"
)
type FileInfo struct {
	rootPath string
	fileContent  string
	fileName string
}

func httpPost(f *FileInfo, targetUrl string) {
    resp, err := http.PostForm(targetUrl,
        url.Values{"fileName" : {f.fileName}, "rootPath" : {f.rootPath}, "fileContent" : {f.fileContent}})
    if err != nil {
        fmt.Println(err)
    }

    defer resp.Body.Close()
    log.Printf("succ %s", f.fileName)
}

func main()  {
    	fmt.Println("just for test")

	cfg, err := config.ParseJsonFile("D:\\work\\go\\test\\client.conf")
	if err != nil {
		log.Println(err.Error())
		return
	}
	host, _ := cfg.String("host")
	log.Println(host)
	srcRootPath, _ := cfg.String("singleMap.from")
	dstRootPath, _ := cfg.String("singleMap.to")
	c := make(chan notify.EventInfo, 1)

	for{
		if err := notify.Watch(srcRootPath, c, notify.All); err != nil {
			log.Fatal(err)
		}
		defer notify.Stop(c)

		ei := <-c
		fileInfo := new(FileInfo)

		startIndex := strings.Index(ei.Path(), srcRootPath) + len(srcRootPath)
		realName := ei.Path()[startIndex + 1:]


		fileInfo.fileName = realName
		fileInfo.rootPath = dstRootPath

		fileContentBytes, _ := ioutil.ReadFile(ei.Path())
		fileInfo.fileContent = string(fileContentBytes)
		httpPost(fileInfo, host)
		defer func(){
			log.Println(realName + " " + ei.Path())
			if e := recover(); e != nil{
				log.Fatal(e)
			}
		}()
	}
}