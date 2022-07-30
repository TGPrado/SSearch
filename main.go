package main

import (
    "fmt"
    "flag"
    "os"
    "bufio"
    "net/http"
    "time"
)

var colors = map[int]string{
    20: "\u001b[32m",
    30: "\u001b[32m",
    40: "\u001b[31m",
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

func waitForServer(url string) int{
    const timeout = 1 * time.Minute
    deadline := time.Now().Add(timeout)
    for tries :=0; time.Now().Before(deadline); tries++{
        resp, err := http.Get(url)
        if err ==  nil{
            return resp.StatusCode
        }
        time.Sleep(time.Second << uint(tries))
    }
    return 0
}

func printText(url string, statusCode int){
    decimalStatusCode := statusCode/10
    fmt.Println(url, " ",colors[decimalStatusCode], statusCode, "\u001b[0m")
}


func main(){
    paths := []string{"/v1/api/docs", 
                      "/v1/api/openapi.json", 
                      "/v1/docs", 
                      "/v1/openapi.json", 
                      "/docs", 
                      "/openapi.json", 
                      "/api/docs", 
                      "/api/openapi.json",
                      "/v2/docs", 
                      "/v2/openapi.json", 
                      "/v2/api/docs", 
                      "/v2/api/openapi.json", 
                  }     

    gFlag := flag.Bool("g", false, "only generate swagger endpoints")
    flag.Parse()
    

    endpoints := generateEndpoints(paths)    

    countDataChan := make(chan struct{}, len(endpoints))


    if *gFlag{
        printArray(endpoints)
        return
    }
    

    for _,url := range endpoints{
        go func(){
            fmt.Println(url)
            defer func(){countDataChan <- struct{}{}}()
        }()
    }
    for range endpoints{
        <-countDataChan
    }
}
