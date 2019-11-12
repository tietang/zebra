package main

import (
    "fmt"
    "github.com/go-ini/ini"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    gopath := os.Getenv("GOPATH")
    fmt.Println(gopath)
    f := "../proxy.ini"
    fileReader, err := os.Open(f)
    defer fileReader.Close()

    if err != nil {
        fmt.Println("read file: ")
    }
    p, err := ini.Load(fileReader)

    fileName := "./zebra_const_default.go"
    fmt.Println(os.Remove(fileName))
    fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)

    fmt.Println("const(")
    fd.WriteString("package main\n")
    fd.WriteString("import \"time\"\n")

    sections := p.Sections()
    for _, section := range sections {

        name := strings.TrimSpace(section.Name())
        fmt.Println(name)
        fd.WriteString("     //")
        fd.WriteString(name)
        fd.WriteString("\n")
        fd.WriteString("const(\n")
        keyLines := ""
        valueLines := ""
        for _, kv := range section.Keys() {

            kvName := strings.TrimSpace(kv.Name())
            key := strings.Join([]string{name, kvName}, ".")
            keyValue := "\"" + key + "\""
            key = strings.Replace(key, ".", "_", -1)
            keykey := "     KEY_" + strings.ToUpper(key)
            key = "     DEFAULT_" + strings.ToUpper(key)

            value := strings.TrimSpace(kv.String())
            line := key + " = \"" + value + "\""
            _, err := strconv.ParseFloat(value, 64)
            if err == nil {
                line = key + " = " + value
            }
            d, err := time.ParseDuration(value)
            //s := time.Duration(d.Nanoseconds()/int64(time.Millisecond)) * time.Millisecond
            if err == nil {
                s := "time.Duration(" + strconv.Itoa(int(d.Nanoseconds()/int64(time.Millisecond))) + ") * time.Millisecond"
                line = key + " = " + s
            }
            keyLines += keykey + "=" + keyValue + "\n"
            valueLines += line + "\n"
        }
        fd.WriteString("     //key \n")
        fmt.Println(keyLines)
        fd.WriteString(keyLines)

        fmt.Println(valueLines)
        fd.WriteString("     //default value\n")
        fd.WriteString(valueLines)
        fmt.Println(")")
        fd.WriteString(")\n\n")
    }

    fd.Close()
}
