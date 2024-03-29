package main

import (
  "fyne.io/fyne/v2/app"
  "github.com/ipv6rslimited/configurator"
  "os"
  "fmt"
)

func main() {
  myApp := app.New()

  if len(os.Args) < 2 {
    fmt.Printf("Usage: %s path/to/config.json", os.Args[0])
    os.Exit(1)
  }

  configFilePath := os.Args[1]

  configurator.NewWindow(myApp, configFilePath, "IPv6rs", "https://ipv6.rs")
  myApp.Run()
}
