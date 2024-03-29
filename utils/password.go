package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
	"sync"
)

var ConfPasswordComplexitys = []string{"lower", "upper", "digit", "spec"}

// complexity contains information about a particular kind of password complexity
type complexity struct {
	ValidChars string
	TrNameOne  string
}

var (
	matchComplexityOnce sync.Once
	validChars          string
	requiredList        []complexity

	charComplexities = map[string]complexity{
		"lower": {
			`abcdefghijklmnopqrstuvwxyz`,
			"form.password_lowercase_one",
		},
		"upper": {
			`ABCDEFGHIJKLMNOPQRSTUVWXYZ`,
			"form.password_uppercase_one",
		},
		"digit": {
			`0123456789`,
			"form.password_digit_one",
		},
		"spec": {
			` !"#$%&'()*+,-./:;<=>?@[\]^_{|}~` + "`",
			"form.password_special_one",
		},
	}
)

// NewComplexity for preparation
func NewComplexity() {
	matchComplexityOnce.Do(func() {
		setupComplexity(ConfPasswordComplexitys)
	})
}

func setupComplexity(values []string) {
	if len(values) != 1 || values[0] != "off" {
		for _, val := range values {
			if v, ok := charComplexities[val]; ok {
				validChars += v.ValidChars
				requiredList = append(requiredList, v)
			}
		}
		if len(requiredList) == 0 {
			// No valid character classes found; use all classes as default
			for _, v := range charComplexities {
				validChars += v.ValidChars
				requiredList = append(requiredList, v)
			}
		}
	}
	if validChars == "" {
		// No complexities to check; provide a sensible default for password generation
		validChars = charComplexities["lower"].ValidChars + charComplexities["upper"].ValidChars + charComplexities["digit"].ValidChars
	}
}

// IsComplexEnough return True if password meets complexity settings
func IsComplexEnough(pwd string) bool {
	NewComplexity()
	if len(validChars) > 0 {
		for _, req := range requiredList {
			if !strings.ContainsAny(req.ValidChars, pwd) {
				return false
			}
		}
	}
	return true
}

// Generate 生成一个n位数随机串
func Generate(n int) (string, error) {
	NewComplexity()
	buffer := make([]byte, n)
	max := big.NewInt(int64(len(validChars)))
	for j := 0; j < n; j++ {
		rnd, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		buffer[j] = validChars[rnd.Int64()]
	}
	return string(buffer), nil
}

func GenRandomSalt() string {
	salt, err := Generate(8)
	if err != nil {
		fmt.Println("GenRandomSalt err ", err)
		return fmt.Sprintf("%x", mrand.Int31())
	}

	return salt

}
