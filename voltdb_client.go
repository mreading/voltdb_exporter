package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
)

type Stats struct {
    state   []byte
    cpu     []byte
    txns    []byte
    latency []byte
    ram     []byte
}

func initializeClient() {
    request := fmt.Sprintf("http://%s/api/1.0/?Procedure=@Ping&admin=false&User=%s&Password=%s", addr, user, pass)
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
    basepath := "http://%s/api/1.0/?Procedure=%s&Parameters=%s&admin=false&User=%s&Password=%s"
    return []string {
        fmt.Sprintf(basepath, addr, "@SystemInformation", "['OVERVIEW']", user, pass),
        fmt.Sprintf(basepath, addr, "@Statistics", "['CPU',0]", user, pass),
        fmt.Sprintf(basepath, addr, "@Statistics", "['LATENCY',0]", user, pass),
        fmt.Sprintf(basepath, addr, "@Statistics", "['LATENCY',0]", user, pass),
        fmt.Sprintf(basepath, addr, "@Statistics", "['MEMORY',0]", user, pass),
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
                state:   data[0],
                cpu:     data[1],
                txns:    data[2],
                latency: data[3],
                ram:     data[4],
             }

    return &stats, nil
}
