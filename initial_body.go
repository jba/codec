// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

const initialBody = `
«/*» Template body for the beginning of the file. «*/»

// Code generated by the codec package. DO NOT EDIT.

package «.Package»

import (
	«range .StdImports»
		«with .ID»«.»«end» "«.Path»"
	«- end»

	«range .OtherImports»
		«with .ID»«.»«end» "«.Path»"
	«- end»
)
`