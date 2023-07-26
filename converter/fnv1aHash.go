package converter

import (
	"fmt"
	"hash/fnv"
)

// Fnv1aHash returns a 32 bit FNV-1a hash
// Usage example:
//
//	hash := converter.Fnv1aHash("Vienna")
//	fmt.Println(hash)
//	// prints "712dc882"
func Fnv1aHash(input string) string {
	ha := fnv.New32a()
	_, _ = ha.Write([]byte(input))
	// Achtung! das 08 ist wichtig - muessen genau
	// 8 chars sein z.B. hash := common.Fnv1aHash("uswesampol01")
	// hat eine fuehrende 0!
	return fmt.Sprintf("%08x", ha.Sum32())
}
