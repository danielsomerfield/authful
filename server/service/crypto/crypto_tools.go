package crypto

import "github.com/danielsomerfield/authful/util"
import (
	"golang.org/x/crypto/scrypt"
	"log"
	"regexp"
	"encoding/base64"

	"bytes"
	"strconv"
	"math"
	"fmt"
)

type BuildPwdHash func(string) string

func ScryptHash(input string) string {
	N := 16
	r := 8
	p := 1
	salt := util.GenerateRandomBytes(16)
	keyLen := 32

	hash, err := scrypt.Key([]byte(input), salt, int(math.Pow(2, float64(N))), r, p, keyLen)

	if err != nil {
		log.Panic("Bad configuration of scrypt.")
	}

	encodedSalt := base64.RawStdEncoding.EncodeToString([]byte(salt))
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$scrypt$ln=%d,r=%d,p=%d$%s$%s", N, r, p, encodedSalt, encodedHash)
}

func ValidateScrypt(data string, hash string) bool {

	re := regexp.MustCompile("\\$scrypt\\$ln=([0-9]+),r=([0-9]+),p=([0-9]+)\\$([a-zA-Z0-9/]+)\\$([a-zA-Z0-9+/]+).*")

	matches := re.FindAllStringSubmatch(hash, -1)
	if len(matches) != 1 {
		return false
	}

	//TODO: handle excessively large input numbers

	m := matches[0]
	ln, _ := strconv.ParseInt(m[1], 10, 32)
	r, _ := strconv.ParseInt(m[2], 10, 32)
	p, _ := strconv.ParseInt(m[3], 10, 32)
	salt := m[4]
	checksum := m[5]

	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {return false}
	checksumBytes, err := base64.RawStdEncoding.DecodeString(checksum)
	if err != nil {return false}
	key, err := scrypt.Key([]byte(data), saltBytes, int(math.Pow(2, float64(ln))), int(r), int(p), 32)
	if err != nil {return false}
	equal := bytes.Equal(checksumBytes, key)

	return err == nil && equal
}
