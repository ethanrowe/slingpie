package main

import (
    "fmt"
    "log"
    "os"
    "github.com/ethanrowe/slingpie/venv"
)

func showHelp() {
    fmt.Println("Provide the path to the source virtualenv.")
}

func main() {
    args := os.Args[1:]
    if len(args) == 1 {
        src, err := venv.WrapVenv(args[0])
        if err == nil {
            err = src.Construct()
            if err == nil {
                err = src.Stream(os.Stdout)
                // src.Destroy()
            }
        }
        if err != nil {
            log.Fatalln("An error occurred:", err)
        }
    } else {
        showHelp()
    }
}

