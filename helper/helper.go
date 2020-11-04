package helper

import "encoding/json"

type RequestJson struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type ResponseJson struct {
	Ok     bool   `json:"ok"`
	From   string `json:"from"`
	To     string `json:"to"`
	Result string `json:"result"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
func GetRequestJson(jsonString []byte) RequestJson {
	var requestStruct RequestJson = RequestJson{}
	if err := json.Unmarshal(jsonString, &requestStruct); err != nil {
		panic(err)
	}
	return requestStruct
}

func GetResponseJson(jsonString []byte) ResponseJson {
	var responseStruct ResponseJson = ResponseJson{}
	if err := json.Unmarshal(jsonString, &responseStruct); err != nil {
		panic(err)
	}
	return responseStruct
}
func GenerateResponseBody(from string, to string, result string) string {
	var responseStruct ResponseJson = ResponseJson{
		Ok:     true,
		From:   from,
		To:     to,
		Result: result,
	}
	resultBytes, err := json.Marshal(responseStruct)
	Check(err)
	var resultString string = string(resultBytes)
	return resultString
}
