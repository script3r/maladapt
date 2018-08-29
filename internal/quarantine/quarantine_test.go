package quarantine

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestComputeSHA256(t *testing.T) {
	q := NewManager(NewZipQuarantiner(""), nil)
	test := "test"

	res := q.computeSHA256([]byte(test))

	fmt.Printf("%x\n", res)

	//Use for urls?
	fmt.Println(base64.RawURLEncoding.EncodeToString(res[:]))
	fmt.Println(base64.URLEncoding.EncodeToString(res[:]))
	fmt.Println(res[:])
	fmt.Println(base64.RawURLEncoding.DecodeString(base64.RawURLEncoding.EncodeToString(res[:])))
}
