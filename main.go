package main

import (
    "fmt"
    "flag"
    "os"
    "bufio"
)



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

func main(){
    paths := []string{"/v1/api/docs", 
                      "/v1/api/openapi.json", 
                      "v1/docs", 
                      "v1/openapi.json", 
                      "/docs", 
                      "/openapi.json", 
                      "/api/docs", 
                      "/api/openapi.json",
                      "/v2/docs", 
                      "v2/openapi.json", 
                      "/v2/api/docs", 
                      "/v2/api/openapi.json", 
                  }     

    gFlag := flag.Bool("g", false, "only generate swagger endpoints")
    vFlag := flag.Bool("v", false, "print all status codes")
    flag.Parse()
    
    endpoints := generateEndpoints(paths)    
    fmt.Println(endpoints)
    if *gFlag{
        fmt.Println(*gFlag, *vFlag)
        return
    }

}
