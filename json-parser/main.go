package main

import (
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"os"
	"strconv"
	"strings"
)

const (
	trueText   = "true"
	falseText  = "false"
	nullText   = "null"
	jsonArray  = "array"
	jsonObject = "jsonObject"
)

func main() {
	// Check if the correct number of command-line arguments is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: programName <file_name>")
		os.Exit(1)
	}

	// Get the file name from the command-line arguments
	fileName := os.Args[1]

	fmt.Println(fileName)
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic("no file present")
	}
	jsonText := string(data)
	if output, err := marshal(jsonText); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println(output)
	}
}
func marshal(jsonText string) (any, error) {
	//trim white spaces
	jsonText = strings.TrimSpace(jsonText)
	// check if its a valid json if yes then return the type which can be array or jsonObject
	jsonType, err := checkJSONType(jsonText)
	if err != nil {
		return nil, err
	}
	//remove 1st and last character and extract the string between
	remainingString := strings.TrimSpace(jsonText[1 : len(jsonText)-1])
	if len(remainingString) == 0 {
		if jsonType == jsonArray {
			return make([]interface{}, 0), nil
		} else {
			return make(map[string]any), nil
		}
	}
	// make tokens using commas
	commaSeparatedTokens, err := createCommaSeparatedTokens(remainingString)
	if err != nil {
		return nil, err
	}
	if len(commaSeparatedTokens) == 0 {
		return nil, nil
	}
	if jsonType == jsonObject {
		var resultantMap = make(map[string]any)
		// iterate through comma separated token
		for i := 0; i < len(commaSeparatedTokens); i++ {
			keyValSublist := strings.SplitN(strings.TrimSpace(commaSeparatedTokens[i]), ":", 2)
			if len(keyValSublist) < 2 {
				return nil, errors.New("missing key or val")
			}
			key := strings.TrimSpace(keyValSublist[0])
			val := strings.TrimSpace(keyValSublist[1])
			if !isProperlyDoubleQuoted(key) {
				return nil, errors.New("invalid key")
			}
			//check if key is already declared
			if _, exist := resultantMap[key]; exist {
				return nil, errors.New("key already exist")
			}
			resultantVal, err := evalValue(val)
			if err != nil {
				return nil, err
			}
			resultantMap[key] = resultantVal
		}
		return resultantMap, nil
	} else {
		// iterate through comma separated token
		var resultantArray = make([]interface{}, 0)
		for i := 0; i < len(commaSeparatedTokens); i++ {
			val := strings.TrimSpace(commaSeparatedTokens[i])
			resultantVal, err := evalValue(val)
			if err != nil {
				return nil, err
			}
			resultantArray = append(resultantArray, resultantVal)
		}
		return resultantArray, nil
	}
}

func createColonSeparatedTokens(text string) []string {
	resultantList := make([]string, 0)

	return resultantList
}

func checkJSONType(jsonText string) (string, error) {
	//check length of string is 0
	if len(jsonText) == 0 {
		return "", errors.New("empty string")
	}
	firstChar := jsonText[0]
	lastChar := jsonText[len(jsonText)-1]
	if firstChar == '{' && lastChar == '}' {
		return jsonObject, nil
	}
	if firstChar == '[' && lastChar == ']' {
		return jsonArray, nil
	}
	return "", errors.New("invalid json")
}

func evalValue(val string) (any, error) {
	// check for cases like "abc' or 'abc" or "acv or ac' or ac" or 'ac or " or ' or 'abc'
	if isInvalidQuotation(val) {
		return nil, errors.New("value has invalid quotation")
	}
	if isProperlyDoubleQuoted(val) {
		return val, nil
	}
	// Check if the string represents a boolean
	if val == trueText || val == falseText {
		boolValue, _ := strconv.ParseBool(val)
		return boolValue, nil
	}

	// Check if the string represents a number
	if num, err := strconv.ParseFloat(val, 64); err == nil {
		return num, nil
	}
	// Check if the string is "null"
	if val == nullText {
		return nil, nil
	}

	_, err := checkJSONType(val)
	if err != nil {
		return nil, errors.New("invalid value")
	}
	return marshal(val)
}

func isProperlyDoubleQuoted(str string) bool {
	firstChar := str[0]
	lastChar := str[len(str)-1]
	return firstChar == '"' && lastChar == '"'
}

func isInvalidQuotation(str string) bool {
	// Check if the string is empty
	if len(str) == 0 {
		return true // An empty string is considered to begin and end with a quote
	}

	// Check if the string begins or ends with a single quote or double quote
	firstChar := str[0]
	lastChar := str[len(str)-1]

	return (firstChar == '\'' && lastChar == '\'') || (firstChar != '\'' && lastChar == '\'') || (firstChar == '\'' && lastChar != '\'') || (firstChar == '"' && lastChar != '"') || (firstChar != '"' && lastChar == '"')
}

//
//abc- [abc]
//
//abc,- error
//
//a,b,c -[a,b,c]
//
//a,{},c - [a,string,c]
//a,{a,},c - error
//a,[a,b,c] - [a,string]
//a,{"a":"b,c","c":[a,b,c]} - [a,string]
// [{"a:"b,"},c] -> [string,c]
/*

  "key": "value",
  "key-n": 101,
  "keyo": {
    "innerkey": "inner,,value",
    "inner key1": "inner value"
  },
  "key-l": ["list value"]


-> [ "key": "value","key-n": 101,string, ]

"a","b",{} -> handled
"a","b",{} -> handled
"a",,"b",{} -> handled
"a","b",{}, -> handled
,"a","b",{} -> handled
"ab,c","b",{} -> [ab,c,b,{}]
"ab,c","b",{"abc":"abc","efg":"def"} -> ["ab,c","b",

*/
func createCommaSeparatedTokens(text string) ([]string, error) {
	resultantList := make([]string, 0)
	// check if the string is ending with comma
	if text[0] == ',' || text[len(text)-1] == ',' {
		return nil, errors.New("comma at the end")
	}
	j := 1
	i := 0
	for j < len(text) {
		if text[j] == '[' {
			st := stack.New()
			st.Push('[')
			j++
			for st.Len() != 0 {
				if j == len(text) {
					return nil, errors.New("invalid json")
				}
				if text[j] == ']' && st.Peek() == '[' {
					st.Pop()
				}
				if text[j] == '[' {
					st.Push(']')
				}
				j++
			}
			if len(text) == j {
				break
			}
		}
		if text[j] == '{' {
			st := stack.New()
			st.Push('{')
			j++
			for st.Len() != 0 {
				if j == len(text) {
					return nil, errors.New("invalid json")
				}
				if text[j] == '}' && st.Peek() == '{' {
					st.Pop()
				}
				if text[j] == '{' {
					st.Push('{')
				}
				j++
			}
			if len(text) == j {
				break
			}
		}
		if text[j] == '"' {
			for text[j] != '"' {
				j++
			}
		}
		if text[j] == ',' && text[j-1] == ',' {
			return nil, errors.New("no text between comma")
		}
		if text[j] == ',' {
			resultantList = append(resultantList, text[i:j])
			i = j + 1
		}
		j++
	}
	resultantList = append(resultantList, text[i:j])
	return resultantList, nil
}
