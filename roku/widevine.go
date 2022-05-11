package roku

import (
   "bytes"
   "encoding/base64"
   "encoding/hex"
   "github.com/89z/format/json"
   "net/http"
   "strings"
)

var pssh = []byte{
   0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
   // Widevine UUID:
   0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
   0, 0, 0, 0,
   // length + KID:
   8, 0, 0, 0, 0, 0, 0, 0,
}

type CrossSite struct {
   cookie *http.Cookie // has own String method
   token string
}

func NewCrossSite() (*CrossSite, error) {
   // this has smaller body than www.roku.com
   req, err := http.NewRequest("GET", "https://therokuchannel.roku.com", nil)
   if err != nil {
      return nil, err
   }
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var site CrossSite
   for _, cook := range res.Cookies() {
      if cook.Name == "_csrf" {
         site.cookie = cook
      }
   }
   scan, err := json.NewScanner(res.Body)
   if err != nil {
      return nil, err
   }
   scan.Split = []byte("\tcsrf:")
   scan.Scan()
   scan.Split = nil
   if err := scan.Decode(&site.token); err != nil {
      return nil, err
   }
   return &site, nil
}

func (c CrossSite) Playback(id string) (*Playback, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(map[string]string{
      "mediaFormat": "mpeg-dash",
      "rokuId": id,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://therokuchannel.roku.com/api/v3/playback", buf,
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "CSRF-Token": {c.token},
      "Content-Type": {"application/json"},
   }
   req.AddCookie(c.cookie)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   play := new(Playback)
   if err := json.NewDecoder(res.Body).Decode(play); err != nil {
      return nil, err
   }
   return play, nil
}

type Playback struct {
   DRM struct {
      Widevine struct {
         LicenseServer string
      }
   }
}

func (p Playback) Widevine() (*Widevine, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(map[string]string{
      "buildInfo": "",
      "license": p.DRM.Widevine.LicenseServer,
      "pssh": base64.StdEncoding.EncodeToString(pssh),
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", "https://getwvkeys.cc/api", buf)
   if err != nil {
      return nil, err
   }
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   vine := new(Widevine)
   if err := json.NewDecoder(res.Body).Decode(vine); err != nil {
      return nil, err
   }
   return vine, nil
}

type Widevine struct {
   Keys []struct {
      Key string
   }
}

func (w Widevine) Key() ([]byte, error) {
   for _, each := range w.Keys {
      _, key, ok := strings.Cut(each.Key, ":")
      if ok {
         return hex.DecodeString(key)
      }
   }
   return nil, errorString(`":" not found`)
}