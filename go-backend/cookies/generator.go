package cookies

import "math/rand/v2"

// Generate random cookie with format "xxxx-xxxx..."
// The cookie can contain characters ranging from a to z
// lenght 3 equals to xxxx-xxxx-xxxx
func GenerateCookie(lenght int) string {

	out := ""
	for i := 0; i < lenght; i++ {
		func() {
			for j := 0; j < 4; j++ {
				// Returns random letters between byte value 97 and 122 (a - z)
				out += string(byte(rand.IntN(122-97) + 97))
			}
			out += "-"

		}()
	}
	// Remove last ugly line
	return out[:len(out)-1]
}
