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
	if err := e.Encode([]interface{}{1, "hello", true}); err != nil {
		log.Fatal(err)
	}

	d := NewDecoder(bytes.NewReader(buf.Bytes()), nil)
	got, err := d.Decode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)

	// Output:
	// a b c
}
