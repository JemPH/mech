package paramount

import (
   "github.com/89z/mech/widevine"
   "os"
   "testing"
)

func Test_Session(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/mech/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_ID, err := os.ReadFile(home + "/mech/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   key_ID, err := widevine.Key_ID("3be8be937c98483184b294173f9152af")
   if err != nil {
      t.Fatal(err)
   }
   mod, err := widevine.New_Module(private_key, client_ID, key_ID)
   if err != nil {
      t.Fatal(err)
   }
   sess, err := New_Session(tests[test_type{episode, dash_cenc}])
   if err != nil {
      t.Fatal(err)
   }
   keys, err := mod.Post(sess)
   if err != nil {
      t.Fatal(err)
   }
   if keys.Content().String() != "44f12639c9c4a5a432338aca92e38920" {
      t.Fatal(keys)
   }
}