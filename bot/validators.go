package bot

import (
	"regexp"
	"strings"
)

var (
	// Телефон: + или 8, затем 10-15 цифр (с пробелами, скобками, дефисами)
	// Примеры: +7 999 123-45-67, +79991234567, 8 (999) 123 45 67
	phoneRegex = regexp.MustCompile(`^[+8][\d\s\-\(\)]{9,20}$`)

	// Юзернейм: @ + латиница/цифры/подчёркивание, длина 5-32 символа
	// Примеры: @vasya, @vasya_petrov, @vasya123
	usernameRegex = regexp.MustCompile(`^@[a-zA-Z0-9_]{5,32}$`)
)

// isValidContact — проверяет, что контакт похож на телефон или @username
func isValidContact(input string) bool {
	input = strings.TrimSpace(input)

	if input == "" {
		return false
	}

	// Проверяем как телефон
	if phoneRegex.MatchString(input) {
		// Дополнительно: должно быть хотя бы 10 цифр
		digits := countDigits(input)
		if digits >= 10 && digits <= 15 {
			return true
		}
	}

	// Проверяем как юзернейм
	if usernameRegex.MatchString(input) {
		return true
	}

	return false
}

// countDigits — считает количество цифр в строке
func countDigits(s string) int {
	count := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			count++
		}
	}
	return count
}
