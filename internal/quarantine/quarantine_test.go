package quarantine

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestComputeSHA256(t *testing.T) {
	q := NewQuarantine("nil", NewZipQuarantiner(""))
	test := "test"

	res := q.computeSHA256([]byte(test))

	fmt.Printf("%x\n", res)

	//Use for urls?
	fmt.Println(base64.URLEncoding.EncodeToString(res[:]))
}
