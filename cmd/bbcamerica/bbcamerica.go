package main

import (
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/dash"
   "github.com/89z/mech"
   "github.com/89z/mech/bbcamerica"
   "github.com/89z/mech/widevine"
   "net/http"
   "net/url"
   "os"
)

func (d *downloader) setKey() error {
   privateKey, err := os.ReadFile(d.pem)
   if err != nil {
      return err
   }
   clientID, err := os.ReadFile(d.client)
   if err != nil {
      return err
   }
   kID, err := d.Protection().KID()
   if err != nil {
      return err
   }
   mod, err := widevine.NewModule(privateKey, clientID, kID)
   if err != nil {
      return err
   }
   site, err := bbcamerica.NewCrossSite()
   if err != nil {
      return err
   }
   play, err := site.Playback(d.Meta.ID)
   if err != nil {
      return err
   }
   keys, err := mod.Post(play.DRM.Widevine.LicenseServer, nil)
   if err != nil {
      return err
   }
   d.key = keys.Content().Key
   return nil
}

type downloader struct {
   *dash.Period
   *url.URL
   base string
   client string
   info bool
   key []byte
   pem string
}

func (d downloader) DASH(nid, video int64) error {
   auth, err := bbcamerica.NewUnauth()
   if err != nil {
      return err
   }
   play, err := auth.Playback(nid)
   if err != nil {
      return err
   }
   source := play.DASH()
   fmt.Println("GET", source.Src)
   res, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   d.URL = res.Request.URL
   d.Period, err = dash.NewPeriod(res.Body)
   if err != nil {
      return err
   }
   d.base = play.Body.Data.PlaybackJsonData.Name
   return d.download(video, dash.Video)
}

func (d *downloader) download(band int64, fn dash.PeriodFunc) error {
   if band == 0 {
      return nil
   }
   reps := d.Represents(fn)
   rep := reps.Represent(band)
   if d.info {
      for _, each := range reps {
         if each.Bandwidth == rep.Bandwidth {
            fmt.Print("!")
         }
         fmt.Println(each)
      }
   } else {
      if d.key == nil {
         err := d.setKey()
         if err != nil {
            return err
         }
      }
      ext, err := mech.ExtensionByType(rep.MimeType)
      if err != nil {
         return err
      }
      file, err := os.Create(d.base + ext)
      if err != nil {
         return err
      }
      defer file.Close()
      init, err := rep.Initialization(d.URL)
      if err != nil {
         return err
      }
      fmt.Println("GET", init)
      res, err := http.Get(init.String())
      if err != nil {
         return err
      }
      defer res.Body.Close()
      if _, err := file.ReadFrom(res.Body); err != nil {
         return err
      }
      media, err := rep.Media(d.URL)
      if err != nil {
         return err
      }
      pro := format.ProgressChunks(file, len(media))
      for _, addr := range media {
         res, err := http.Get(addr.String())
         if err != nil {
            return err
         }
         pro.AddChunk(res.ContentLength)
         if err := dash.Decrypt(pro, res.Body, d.key); err != nil {
            return err
         }
         if err := res.Body.Close(); err != nil {
            return err
         }
      }
   }
   return nil
}