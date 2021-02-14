#!/usr/bin/env bash
# Copyright 2021 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# TODO: replace with go:embed when go 1.17 comes out.

for f in *.tmpl; do
  base=$(basename $f .tmpl)
  outfile=${base}_body.go
  cat > $outfile  <<END
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

const ${base}Body = \`
END
  cat $f >> $outfile
  echo "\`" >> $outfile
done
