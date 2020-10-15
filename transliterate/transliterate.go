package transliterate

import (
	"transliteration_bot/config"
)

func Transliterate(msg string) string {
	var cyrillic string = config.GetConfig().Cyrillic
	var latin string = config.GetConfig().Latin

	const numberOfChars int = 66

	var cyrillicRunes [numberOfChars]rune
	var latinRunes [numberOfChars]rune
	var trans string = ""

	var correspondence map[rune]rune = make(map[rune]rune)

	var j uint8 = 0
	for _, runeValue := range cyrillic {
		cyrillicRunes[j] = runeValue
		j += 1
	}
	j = 0
	for _, runeValue := range latin {
		latinRunes[j] = runeValue
		j += 1
	}

	for i := 0; i < numberOfChars; i += 1 {
		correspondence[cyrillicRunes[i]] = latinRunes[i]
	}

	for _, runeValue := range msg {
		if latinChar, ok := correspondence[runeValue]; ok {
			//do something here
			trans += string(latinChar)
		} else {
			trans += string(runeValue)
		}

	}

	return trans

}
