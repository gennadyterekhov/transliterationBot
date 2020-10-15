package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type config struct {
	Port     int    `json:"port"`
	APIKey   string `json:"api_key"`
	Cyrillic string `json:"cyrillic"`
	Latin    string `json:"latin"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getConfig() config {
	configStr, err := ioutil.ReadFile("config.json")
	check(err)

	var configObj config = config{}
	if err := json.Unmarshal([]byte(configStr), &configObj); err != nil {
		panic(err)
	}
	return configObj
}

func transliterate(msg string) string {
	var cyrillic string = getConfig().Cyrillic
	var latin string = getConfig().Latin
	// fmt.Printf("%d\n", len(cyrillic))
	// fmt.Printf("%d\n", len(latin))
	const numberOfChars int = 66
	// return "asdf"
	var cyrillicRunes [numberOfChars]rune
	var latinRunes [numberOfChars]rune
	var trans string = ""

	var correspondence map[rune]rune = make(map[rune]rune)
	// correspondence['я'] = 'ä'
	var j uint8 = 0
	for _, runeValue := range cyrillic {
		cyrillicRunes[j] = runeValue
		// fmt.Printf("%#U", runeValue)

		j += 1
	}
	j = 0
	for _, runeValue := range latin {
		// fmt.Printf("%d = %#U\n", i, runeValue)
		latinRunes[j] = runeValue
		j += 1
	}
	// fmt.Printf("ar: %v", cyrillicRunes)
	// fmt.Printf("\n")
	// fmt.Printf("ar: %v", latinRunes)

	for i := 0; i < numberOfChars; i += 1 {
		correspondence[cyrillicRunes[i]] = latinRunes[i]
	}
	// fmt.Printf("\n%v", correspondence)

	for _, runeValue := range msg {
		if latinChar, ok := correspondence[runeValue]; ok {
			//do something here
			trans += string(latinChar)
		} else {
			trans += string(runeValue)
		}

	}

	return trans

	fmt.Printf("\nmsg:%s\n", msg)
	// return "ss"
	fmt.Printf("\n[latin]:%s\n[cyrillic]:%s\n", latin, cyrillic)
	//
	fmt.Printf("bytes: ")
	for i, runeValue := range msg {
		fmt.Printf("%x = %#U ", i, runeValue)
		// trans += string(runeValue)
		var charIndex int = strings.Index(cyrillic, string(runeValue))
		fmt.Printf("\n[ИНДЕКС]:%d\n", charIndex)
		fmt.Printf("\n[чар по ИНДЕКСУ]:%c\n", cyrillic[120])
		if charIndex != -1 {
			// trans += string(latin[charIndex])
			// fmt.Printf("\n[получен кириллический символ][current char]:%d\n", int32(latin[charIndex]))
		} else {
			// trans += string(runeValue)
			// fmt.Printf("\n[получен некириллический символ][current char]:%d\n", int32(msg[i]))
		}

	}

	//
	// var charIndex int = strings.Index(cyrillic, string(msg[0]))
	// fmt.Printf("\n[ИНДЕКС]:%d\n", charIndex)
	return trans
	for i := 0; i < len(msg); i += 1 {

		var charIndex int = strings.Index(cyrillic, string(msg[i]))
		fmt.Printf("\n[ИНДЕКС]:%d\n", charIndex)
		if charIndex != -1 {
			trans += string(latin[charIndex])
			fmt.Printf("\n[получен кириллический символ][current char]:%d\n", int32(latin[charIndex]))
		} else {
			trans += string(msg[i])
			fmt.Printf("\n[получен некириллический символ][current char]:%d\n", int32(msg[i]))
		}
		fmt.Printf("\n[current message]:%s\n", trans)
	}

	return trans
}

func main() {
	// var buns string = "съешь ещё этих мягких французских булок да выпей же чаю"
	fmt.Println(transliterate("съешь ещё"))
}
