package password

import (
	"math/rand"
	"strings"
)

type PasswordGeneratorOptions struct {
	Length    int
	UpperCase bool
	LowerCase bool
	Numbers   bool
	Special   bool
}

type PasswordGenerator struct {
	Options *PasswordGeneratorOptions
	Runes   *availableRunes
}

type availableRunes struct {
	uppercase []rune
	lowercase []rune
	numbers   []rune
	special   []rune
}

func New(options *PasswordGeneratorOptions) *PasswordGenerator {
	g := &PasswordGenerator{
		Options: options,
	}

	g.Runes = g.generateAvailableRunes()

	return g
}

func (g *PasswordGenerator) generateAvailableRunes() *availableRunes {
	runes := &availableRunes{}

	if g.Options.UpperCase {
		runes.uppercase = generateRunesFromRange('A', 'Z')
	}

	if g.Options.LowerCase {
		runes.lowercase = generateRunesFromRange('a', 'z')
	}

	if g.Options.Numbers {
		runes.numbers = generateRunesFromRange('0', '9')
	}

	if g.Options.Special {
		runes.special = generateSpecialRunes()
	}

	return runes
}

func generateRunesFromRange(start, end rune) []rune {
	var runes []rune

	for r := start; r <= end; r++ {
		runes = append(runes, r)
	}

	return runes
}

func generateSpecialRunes() []rune {
	return []rune{'~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '`', '-', '=', '{', '}', '|', '[', ']', ':', '<', '>', '?', ',', '.'}
}

func (g *PasswordGenerator) Generate() string {
	var password strings.Builder

	optionsCount := 0

	if g.Options.UpperCase {
		optionsCount++
	}
	if g.Options.LowerCase {
		optionsCount++
	}
	if g.Options.Numbers {
		optionsCount++
	}
	if g.Options.Special {
		optionsCount++
	}
	if optionsCount == 0 {
		return ""
	}

	equalShare := g.Options.Length / optionsCount

	if g.Options.UpperCase {
		password.WriteRune(g.Runes.uppercase[rand.Intn(len(g.Runes.uppercase))])
		generateChars(&password, g.Runes.uppercase, equalShare)
	}
	if g.Options.LowerCase {
		password.WriteRune(g.Runes.lowercase[rand.Intn(len(g.Runes.lowercase))])
		generateChars(&password, g.Runes.lowercase, equalShare)
	}
	if g.Options.Numbers {
		password.WriteRune(g.Runes.numbers[rand.Intn(len(g.Runes.numbers))])
		generateChars(&password, g.Runes.numbers, equalShare)
	}
	if g.Options.Special {
		password.WriteRune(g.Runes.special[rand.Intn(len(g.Runes.special))])
		generateChars(&password, g.Runes.special, equalShare)
	}

	additionalChars := g.Options.Length - password.Len()
	if additionalChars > 0 {
		availableOptions := make([][]rune, 0)
		if g.Options.UpperCase {
			availableOptions = append(availableOptions, g.Runes.uppercase)
		}
		if g.Options.LowerCase {
			availableOptions = append(availableOptions, g.Runes.lowercase)
		}
		if g.Options.Numbers {
			availableOptions = append(availableOptions, g.Runes.numbers)
		}
		if g.Options.Special {
			availableOptions = append(availableOptions, g.Runes.special)
		}

		for i := 0; i < additionalChars; i++ {
			optionIndex := rand.Intn(len(availableOptions))
			option := availableOptions[optionIndex]
			password.WriteRune(option[rand.Intn(len(option))])
		}
	}

	runes := []rune(password.String())
	for i := len(runes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func generateChars(password *strings.Builder, runes []rune, count int) {
	for i := 0; i < count; i++ {
		password.WriteRune(runes[rand.Intn(len(runes))])
	}
}
