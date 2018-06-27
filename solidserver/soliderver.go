package solidserver

import (
  "github.com/alexissavin/gorequest"
  "encoding/base64"
  "crypto/tls"
  "crypto/x509"
  "net/http"
  "net/url"
  "fmt"
  "log"
  "io/ioutil"
)

type SOLIDserver struct {
  Host                     string
  Username                 string
  Password                 string
  BaseUrl                  string
  SSLVerify                bool
  AdditionalTrustCertsFile string
}

func NewSOLIDserver(host string, username string, password string, sslverify bool, certsfile string) *SOLIDserver {
  s := &SOLIDserver{
    Host:                     host,
    Username:                 username,
    Password:                 password,
    BaseUrl:                  "https://" + host,
    SSLVerify:                sslverify,
    AdditionalTrustCertsFile: certsfile,
  }

  return s
}

func (s *SOLIDserver) Request(method string, service string, parameters *url.Values) (*http.Response, string, []error) {
  // Get the SystemCertPool, continue with an empty pool on error
  log.Printf("[DEBUG] AdditionalTrustCertsFile = %s", s.AdditionalTrustCertsFile)
  rootCAs, err := x509.SystemCertPool()
  if rootCAs == nil || err != nil {
    rootCAs = x509.NewCertPool()
  }
  if s.AdditionalTrustCertsFile != "" {
    certs, err := ioutil.ReadFile(s.AdditionalTrustCertsFile)
    log.Printf("[DEBUG] Certificates = %s", certs)
    if err != nil {
      log.Fatalf("Failed to append %q to RootCAs: %v", s.AdditionalTrustCertsFile, err)
    }

    log.Printf("[DEBUG] Cert Subjects Before Append = %d", len(rootCAs.Subjects()))
    if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
      log.Println("No certs appended, using system certs only")
    }
    log.Printf("[DEBUG] Cert Subjects After Append = %d", len(rootCAs.Subjects()))
  }
  apiclient := gorequest.New()

  switch method {
  case "post":
    return apiclient.Post(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "put":
    return apiclient.Put(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "delete":
    return apiclient.Delete(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "get":
    return apiclient.Get(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  default:
    return nil, "", nil
  }
}
