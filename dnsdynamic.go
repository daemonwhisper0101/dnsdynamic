// vim:set sw=2 sts=2:
package dnsdynamic

import (
  "io/ioutil"
  "fmt"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "strings"
  "time"
)

type Domain struct {
  Name string
  IP string
}

type Client struct {
  Email, Pass string
}

func NewClient(email, pass string) *Client {
  return &Client{ Email: email, Pass: pass }
}

func (cl *Client)List(opts ...interface{}) ([]Domain, error) {
  domains := []Domain{}

  httpclient := &http.Client{}
  for _, opt := range opts {
    switch v := opt.(type) {
    case http.Transport: httpclient.Transport = &v
    case *http.Transport: httpclient.Transport = v
    case time.Duration: httpclient.Timeout = v
    default:
    }
  }
  jar, err := cookiejar.New(nil)
  if err != nil {
    return domains, fmt.Errorf("cookiejar: %v", err)
  }
  httpclient.Jar = jar

  data := url.Values{}
  data.Set("email", cl.Email)
  data.Set("pass", cl.Pass)
  auth, err := httpclient.PostForm("https://www.dnsdynamic.org/auth.php", data)
  if err != nil {
    return domains, fmt.Errorf("http.PostForm: auth: %v", err)
  }
  defer auth.Body.Close()

  resp, err := httpclient.Get("https://www.dnsdynamic.org/manage.php?page=domains")
  if err != nil {
    return domains, fmt.Errorf("http.PostForm: manage: %v", err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  /*
    parse
    <tr>
      <td width="300"><span class="detailText">my.dynamic.domain.http01.com</span></td>
      <td width="300"><span class="detailText">aaa.bbb.ccc.ddd</span></td>
      <td width="100"><input type="radio" name="domain" value="my.dynamic.domain.http01.com"></center></td>
    </tr>
  */
  rows := strings.Split(string(body), "<tr")
  for _, row := range rows {
    cols := strings.Split(row, "<td")
    if len(cols) != 4 {
      continue
    }
    getval := func(s string) string {
      c1 := strings.Split(s, `detailText">`) // cut head
      c2 := strings.Split(c1[1], "</span>") // cut tail
      return c2[0]
    }
    name := getval(cols[1])
    ip := getval(cols[2])
    // check ip ~= aaa.bbb.ccc.ddd
    n := strings.Split(ip, ".")
    if len(n) != 4 {
      continue
    }
    domains = append(domains, Domain{ Name: name, IP: ip })
  }

  return domains, nil
}
