package main

import (
    "fmt"
    "net/http"
)

func main() {
    res, err := http.Get("https://www.amazon.com/gp/navigation-country/select-country/ref=?ie=UTF8&preferencesReturnUrl=%2F&language=es_US")
    fmt.Println(res)
    fmt.Println(err)
}
