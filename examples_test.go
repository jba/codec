// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"fmt"
	"log"
)

func Example() {
	var buf bytes.Buffer
	e := NewEncoder(&buf, nil)
	for _, x := range []interface{}{1, "hello", true} {
		if err := e.Encode(x); err != nil {
			log.Fatal(err)
		}
	}

	d := NewDecoder(bytes.NewReader(buf.Bytes()), nil)
	for i := 0; i < 3; i++ {
		got, err := d.Decode()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(got)
	}

	// Output:
	// 1
	// hello
	// true
}
