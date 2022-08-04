package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

)

var colors = map[int]string{
    20: "\u001b[32m",
    30: "\u001b[32m",
    40: "\u001b[31m",
    50: "\u001b[31m",
}

type data struct {
    statusCode int
    url        string
}

func generateEndpoints(paths []string) []string{
    var lines []string
    reader := bufio.NewScanner(os.Stdin)
    for reader.Scan(){
        text := reader.Text()
        for _, path := range paths{
            lines = append(lines, text+path)
        }

    }
    return lines
}

func printArray(array []string){
    for _, line := range array{
        fmt.Printf("%s\n", line)
    }
}

func waitForServer(url string, timeFlag int) int{
    var timeout = time.Duration(timeFlag) * time.Minute
    deadline := time.Now().Add(timeout)
    for tries :=0; time.Now().Before(deadline); tries++{
        resp, err := http.Get(url)
        if err ==  nil{
            return resp.StatusCode
        }
        time.Sleep(time.Second << uint(tries))
    }
    return 500
}

func printText(url string, statusCode int, vFlag bool){

    decimalStatusCode := statusCode/10

    if (!vFlag && decimalStatusCode == 20){
        fmt.Println(url, " ",colors[decimalStatusCode], statusCode, "\u001b[0m")
        return
    }
    if (!vFlag){
        return
    }
    fmt.Println(url, " ",colors[decimalStatusCode], statusCode, "\u001b[0m")
    
}


func main(){
    paths := []string{"v1/api/docs", 
                      "v1/api/openapi.json", 
                      "v1/docs", 
                      "v1/openapi.json", 
                      "docs", 
                      "openapi.json", 
                      "api/docs", 
                      "api/openapi.json",
                      "v2/docs", 
                      "v2/openapi.json", 
                      "v2/api/docs", 
                      "v2/api/openapi.json", 
                      "swagger/index.html",
                      "swagger/docs.json",
                      "swagger/openapi.json",
                  }     
    timeFlag := flag.Int("t", 60, "timeout in seconds default is 60")
    gFlag := flag.Bool("g", false, "only generate swagger endpoints")
    vFlag := flag.Bool("v", false, "print all status codes")
    tFlag := flag.Int("t", 10,     "number of threads, default is 10")

    flag.Parse()
    
    threadsChan := make(chan struct{}, *tFlag)

    endpoints := generateEndpoints(paths)    

    cont := 1
    max := len(endpoints)
    printDataChan := make(chan data) 
    
    if *gFlag{
        printArray(endpoints)
        return
    }
    

    for _,url := range endpoints{
        go func(url string){
            var content data
            content.url = url

            defer func(){printDataChan <- content}()
            
            threadsChan <- struct{}{} 
            resp, err  := http.Get(url)
            <-threadsChan
        
            if err == nil{
                content.statusCode = resp.StatusCode
            }
            if err != nil{    
                content.statusCode =  waitForServer(url,*timeFlag)  
            }
        }(url)
    }
    for {
        select{
            case content := <-printDataChan:
                printText(content.url, content.statusCode, *vFlag)
                cont += 1
                if cont == max {
                    return
                }
        }
}

}
