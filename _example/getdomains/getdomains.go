// vim:set sw=2 sts=2:
package main

import (
  "fmt"
  "net/http"
  "net/url"
  "os"

  "github.com/daemonwhisper0101/dnsdynamic"
)

func main() {
  if len(os.Args) < 4 {
    os.Exit(1)
  }
  u, err := url.Parse(os.Args[1])
  if err != nil {
    fmt.Printf("url.Parse: %v\n", err)
    os.Exit(1)
  }
  client := dnsdynamic.NewClient(os.Args[2], os.Args[3])
  tr := &http.Transport{ Proxy: http.ProxyURL(u) }
  domains, err := client.List(tr)
  for _, domain := range domains {
    fmt.Printf("%s: %s\n", domain.Name, domain.IP)
  }
}
