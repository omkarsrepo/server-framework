// Unpublished Work © 2024

package password

import "github.com/sethvargo/go-password/password"

func Generate() string {
	secret, _ := password.Generate(50, 10, 0, false, false)
	return secret
}
