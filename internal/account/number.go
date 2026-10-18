package account

import (
	"fmt"
	"math/rand/v2"
)

func generateNumber() string {
	num := rand.Uint64N(1e12)

	val := fmt.Sprintf("%012d", num)

	return val
}
