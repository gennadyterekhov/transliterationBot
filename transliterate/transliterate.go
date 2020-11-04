package transliterate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"transliteration_bot/helper"
)

func transliterateFromApi(text string) string {
	var requestStruct helper.RequestJson = helper.RequestJson{
		From: "cyrillic",
		To:   "latin",
		Text: text,
	}

	requestStructBytes, err := json.Marshal(requestStruct)
	helper.Check(err)
	// var requestStructString string = string(requestStructBytes)

	response, err := http.Post(os.Getenv("transliterator_api_url"), "application/json", bytes.NewBuffer(requestStructBytes))
	helper.Check(err)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var resposeStruct helper.ResponseJson = helper.GetResponseJson(body)
	return resposeStruct.Result
}

func Transliterate(msg string) string {
	// local
	// var cyrillic string = config.GetConfig().Cyrillic
	// var latin string = config.GetConfig().Latin
	// prod
	return transliterateFromApi(msg)
	// var cyrillic string = os.Getenv("cyrillic")
	// var latin string = os.Getenv("latin")

	// const numberOfChars int = 66

	// var cyrillicRunes [numberOfChars]rune
	// var latinRunes [numberOfChars]rune
	// var trans string = ""

	// var correspondence map[rune]rune = make(map[rune]rune)

	// var j uint8 = 0
	// for _, runeValue := range cyrillic {
	// 	cyrillicRunes[j] = runeValue
	// 	j += 1
	// }
	// j = 0
	// for _, runeValue := range latin {
	// 	latinRunes[j] = runeValue
	// 	j += 1
	// }

	// for i := 0; i < numberOfChars; i += 1 {
	// 	correspondence[cyrillicRunes[i]] = latinRunes[i]
	// }

	// for _, runeValue := range msg {
	// 	if latinChar, ok := correspondence[runeValue]; ok {
	// 		//do something here
	// 		trans += string(latinChar)
	// 	} else {
	// 		trans += string(runeValue)
	// 	}

	// }

	// return trans

}
