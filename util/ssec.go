package util

import (
	"crypto/md5"
	"encoding/base64"
)

type SSEC struct {
	algorithm string
	key       string
}

func NewSSEC(algo, key string) SSEC {
	return SSEC{
		algorithm: algo,
		key:       key,
	}
}

func (s *SSEC) KeyIsSet() bool {
	return s.key != ""
}

func (s *SSEC) Algorithm() string {
	return s.algorithm
}

func (s *SSEC) Base64Key() string {
	return base64.StdEncoding.EncodeToString([]byte(s.key))
}

func (s *SSEC) Base64KeyMD5() string {
	md5Sum := md5.Sum([]byte(s.key))
	return base64.StdEncoding.EncodeToString(md5Sum[:])
}
