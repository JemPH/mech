package youtube

import (
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "sort"
   "strings"
)

const chunk = 10_000_000

func numberFormat(val float64, met []string) string {
   var key int
   for val >= 1000 {
      val /= 1000
      key++
   }
   if key >= len(met) {
      return ""
   }
   return fmt.Sprintf("%.1f", val) + met[key]
}

type Format struct {
   Bitrate bitrate
   ContentLength contentLength `json:"contentLength,string"`
   Height int
   Itag int
   MimeType string
   URL string
}

func (f Format) Write(w io.Writer) error {
   req, err := http.NewRequest("GET", f.URL, nil)
   if err != nil {
      return err
   }
   var pos contentLength
   dump, err := httputil.DumpRequest(req, false)
   if err != nil {
      return err
   }
   os.Stdout.Write(dump)
   for pos < f.ContentLength {
      bytes := fmt.Sprintf("bytes=%d-%d", pos, pos+chunk-1)
      req.Header.Set("Range", bytes)
      fmt.Printf("%d%% %v\n", 100*pos/f.ContentLength, bytes)
      // this sometimes redirects, so cannot use http.Transport
      res, err := new(http.Client).Do(req)
      if err != nil {
         return err
      }
      defer res.Body.Close()
      if res.StatusCode != http.StatusPartialContent {
         return fmt.Errorf("status %v", res.Status)
      }
      if _, err := io.Copy(w, res.Body); err != nil {
         return err
      }
      pos += chunk
   }
   return nil
}

type FormatSlice []Format

func (f FormatSlice) Filter(keep func(Format)bool) FormatSlice {
   var forms FormatSlice
   for _, form := range f {
      if keep(form) {
         forms = append(forms, form)
      }
   }
   return forms
}

func (f FormatSlice) Sort(less ...func(a, b Format) bool) {
   if less == nil {
      less = []func(a, b Format) bool{
         func(a, b Format) bool {
            return b.Height < a.Height
         },
         func(a, b Format) bool {
            f, s := strings.Index, "/mp4;"
            return f(a.MimeType, s) < f(b.MimeType, s)
         },
         func(a, b Format) bool {
            return b.Bitrate < a.Bitrate
         },
      }
   }
   sort.Slice(f, func(a, b int) bool {
      fa, fb := f[a], f[b]
      for _, fn := range less {
         if fn(fa, fb) {
            return true
         }
         if fn(fb, fa) {
            break
         }
      }
      return false
   })
}

type bitrate int64

func (b bitrate) String() string {
   met := []string{"", " kb/s", " mb/s", " gb/s"}
   return numberFormat(float64(b), met)
}

type contentLength int64

func (c contentLength) String() string {
   met := []string{"", " kB", " MB", " GB"}
   return numberFormat(float64(c), met)
}
