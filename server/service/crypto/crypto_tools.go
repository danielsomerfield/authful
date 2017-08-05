package crypto

type BuildPwdHash func(string) []byte

//func ScryptHash(input string) string {
//	N := 16384
//	r := 8
//	p := 1
//	salt := util.GenerateRandomBytes(16)
//	keyLen := 32
//
//	hash, err := scrypt.Key([]byte(input), salt, N, r, p, keyLen)
//
//	if err != nil {
//		log.Panic("Bad configuration of scrypt.")
//	}
//
//	return string(hash)
//}
//
//func ValidateScrypt(data string, hash string) bool {
//
//	re := regexp.MustCompile("\\$scrypt\\$ln=([0-9]+),r=([0-9]+),p=([0-9]+)\\$([a-zA-Z0-9]+)\\$([a-zA-Z0-9\\+]+).*")
//
//	matches := re.FindAllStringSubmatch(hash, -1)
//	if len(matches) != 1 {
//		return false
//	}
//
//	//TODO: handle excessively large input numbers
//	m := matches[0]
//	ln, _ := strconv.ParseInt(m[1], 10, 32)
//	r, _ := strconv.ParseInt(m[2], 10, 32)
//	p, _ := strconv.ParseInt(m[3], 10, 32)
//	salt := m[4]
//	checksum := m[5]
//
//	saltBytes, err := base64.StdEncoding.DecodeString(salt)
//	hashBytes, err := base64.StdEncoding.DecodeString(checksum)
//
//	key, err := scrypt.Key([]byte(data), saltBytes, int(ln), int(r), int(p), 32)
//
//	return err != nil && bytes.Equal(hashBytes, key)
//}
