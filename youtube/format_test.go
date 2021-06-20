package youtube_test

import (
   "github.com/89z/mech/youtube"
   "testing"
)

const desc = "Provided to YouTube by Epitaph\n\nSnowflake · Kate Bush\n\n" +
"50 Words For Snow\n\n" +
"℗ Noble & Brite Ltd. trading as Fish People, under exclusive license to Anti Inc.\n\n" +
"Released on: 2011-11-22\n\nMusic  Publisher: Noble and Brite Ltd.\n" +
"Composer  Lyricist: Kate Bush\n\nAuto-generated by YouTube."

func TestDesc(t *testing.T) {
   p, err := youtube.NewPlayer("XeojXq6ySs4")
   if err != nil {
      t.Fatal(err)
   }
   if p.Description() != desc {
      t.Fatalf("%+v\n", p)
   }
   if p.ViewCount() == 0 {
      t.Fatalf("%+v\n", p)
   }
}

func TestFormat(t *testing.T) {
   p, err := youtube.NewPlayer("eAzIAjTBGgU")
   if err != nil {
      t.Fatal(err)
   }
   // this should fail
   f, err := p.NewFormat(247)
   if err == nil {
      t.Fatalf("%+v\n", f)
   }
}