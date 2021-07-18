// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pgsql

import (
	"fmt"
	"math/big"
	"testing"
)


func TestContain(t *testing.T)  {
	bigint := big.NewInt(123)
	bigstr := bigint.String()
	fmt.Println(bigstr)

	// fmt.Println(strings.Count("date,time", "time"))
	// fmt.Println(strings.IndexAny("hi,中国", "国"))
}