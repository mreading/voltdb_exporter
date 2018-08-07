package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
)

type Stats struct {
    cpu []byte
}

func initializeClient() {
    request := fmt.Sprintf("http://%s/api/1.0/?Procedure=@SystemInformation&Parameters=['DEPLOYMENT']&admin=false&User=%s&Password=%s", 
                           addr, user, pass)
    resp, err := http.Get(request)
    if err != nil {
        log.Fatal(err)
    } else if resp.StatusCode != http.StatusOK {
        log.Fatal(fmt.Sprintf("Failed to connect to client at %s", addr))
    }
    log.Print("Successfully connected to client(s)")
}

func get(path string) ([]byte, error) {
    resp, err := http.Get(path)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    buf, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return buf, nil
}

func setPaths() []string {
    return []string {
        fmt.Sprintf("http://%s/api/1.0/?Procedure=@Statistics&Parameters=['CPU',0]&admin=false&User=%s&Password=%s", addr, user, pass),
    }
}

func scrapeData(paths []string) (data [][]byte, err error) {
    for _, path := range paths {
        buf, err := get(path)
        if err != nil {
            return nil, err
        }
        data = append(data, buf)
    }
    return data, nil
}

func getStats() (*Stats, error) {
    paths := setPaths()

    data, err := scrapeData(paths)
    if err != nil {
        return nil, err
    }

    // attribute data corresponds with data[<index of attribute path in paths>]
    stats := Stats {
                cpu: data[0], 
             }

    return &stats, nil
}
