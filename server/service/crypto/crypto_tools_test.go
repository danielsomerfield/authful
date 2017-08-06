package crypto

import (
	"testing"
	"github.com/danielsomerfield/authful/util"
	"strings"
)

//func TestScryptHash_buildsValidHash(t *testing.T) {
//	hash := ScryptHash("foo")
//	util.AssertTrue(ValidateScrypt("foo", hash), "the hash is valid", t)
//}
//
//func TestValidateHash_failsOnInvalidHash(t *testing.T) {
//	hash := ScryptHash("foo")
//
//	util.AssertFalse(ValidateScrypt(hash), "the hash is valid", t)
//}
//
//func TestScryptHash_buildsDifferentHash(t *testing.T) {
//	hash1 := ScryptHash("foo")
//	hash2 := ScryptHash("foo")
//
//	util.AssertFalse(string(hash1) == string(hash2), "The hashes are equal", t)
//}

func TestScryptValidator(t *testing.T) {
	hash := "$scrypt$ln=16,r=8,p=1$aM15713r3Xsvxbi31lqr1Q$nFNh2CVHVjNldFVKDHDlm4CbdRSCdEBsjjJxD+iCs5E"
	util.AssertTrue(ValidateScrypt("password", hash), "The hash is valid", t)
	util.AssertFalse(ValidateScrypt("password2", hash), "The hash is valid", t)
	util.AssertFalse(ValidateScrypt("password", strings.Replace(hash, "jjJ", "MMM", -1)), "The hash is valid", t)
}
