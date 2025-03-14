package short_link_generator

import "math/rand/v2"

// GenerateShortLink генерирует строку длиной 10 символов, состающую из букв, цифр и символа "_"
func GenerateShortLink() string {
	symbols := []rune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890_")

	shortLink := make([]rune, 10)
	for i := 0; i < 10; i++ {
		shortLink[i] = symbols[rand.IntN(len(symbols))]
	}

	return string(shortLink)
}
