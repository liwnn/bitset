# BitSet
A fast bit array implement in go

## Usage
``` go
package main

import (
	"fmt"

	"github.com/liwnn/bitset"
)

func main() {
	b := bitset.NewSize(8)
	b.Set(1)
	b.Set(100)
	if b.Get(1) {
		fmt.Println("1 is set!")
	}
	if b.Get(100) {
		fmt.Println("100 is set!")
	}
	b.Clear(1)
	if !b.Get(1) {
		fmt.Println("1 is clear!")
	}
}
```
Result:
```
1 is set!
100 is set!
1 is clear!
```