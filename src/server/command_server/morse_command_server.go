package command_server

import (
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type MorseServer struct {
	BaseCommand
}

type EsromServer struct {
	BaseCommand
}

type MorseCode string

const (
	Dot         MorseCode = "."
	Dash        MorseCode = "-"
	Space       MorseCode = " "
	LetterSpace MorseCode = "   "
	WordSpace   MorseCode = "       "
)

var CharToMorseCodeMap = map[rune]MorseCode{
	'A':  Dot + Dash,
	'B':  Dash + Dot + Dot + Dot,
	'C':  Dash + Dot + Dash + Dot,
	'D':  Dash + Dot + Dot,
	'E':  Dot,
	'F':  Dot + Dot + Dash + Dot,
	'G':  Dash + Dash + Dot,
	'H':  Dot + Dot + Dot + Dot,
	'I':  Dot + Dot,
	'J':  Dot + Dash + Dash + Dash,
	'K':  Dash + Dot + Dash,
	'L':  Dot + Dash + Dot + Dot,
	'M':  Dash + Dash,
	'N':  Dash + Dot,
	'O':  Dash + Dash + Dash,
	'P':  Dot + Dash + Dash + Dot,
	'Q':  Dash + Dash + Dot + Dash,
	'R':  Dot + Dash + Dot,
	'S':  Dot + Dot + Dot,
	'T':  Dash,
	'U':  Dot + Dot + Dash,
	'V':  Dot + Dot + Dot + Dash,
	'W':  Dot + Dash + Dash,
	'X':  Dash + Dot + Dot + Dash,
	'Y':  Dash + Dot + Dash + Dash,
	'Z':  Dash + Dash + Dot + Dot,
	'1':  Dot + Dash + Dash + Dash + Dash,
	'2':  Dot + Dot + Dash + Dash + Dash,
	'3':  Dot + Dot + Dot + Dash + Dash,
	'4':  Dot + Dot + Dot + Dot + Dash,
	'5':  Dot + Dot + Dot + Dot + Dot,
	'6':  Dash + Dot + Dot + Dot + Dot,
	'7':  Dash + Dash + Dot + Dot + Dot,
	'8':  Dash + Dash + Dash + Dot + Dot,
	'9':  Dash + Dash + Dash + Dash + Dot,
	'0':  Dash + Dash + Dash + Dash + Dash,
	'.':  Dot + Dash + Dot + Dash + Dot + Dash,
	':':  Dash + Dash + Dash + Dot + Dot + Dot,
	',':  Dash + Dash + Dot + Dot + Dash + Dash,
	';':  Dash + Dot + Dash + Dot + Dash + Dot,
	'?':  Dot + Dot + Dash + Dash + Dot + Dot,
	'=':  Dash + Dot + Dot + Dot + Dash,
	'\'': Dot + Dash + Dash + Dash + Dash + Dot,
	'/':  Dash + Dot + Dot + Dash + Dot,
	'!':  Dash + Dot + Dash + Dot + Dash + Dash,
	'-':  Dash + Dot + Dot + Dot + Dot + Dash,
	'"':  Dot + Dash + Dot + Dot + Dash + Dot,
	'(':  Dash + Dot + Dash + Dash + Dot,
	')':  Dash + Dot + Dash + Dash + Dot + Dash,
	'$':  Dot + Dot + Dot + Dash + Dot + Dot + Dash,
	'&':  Dot + Dash + Dot + Dot + Dot,
	'@':  Dot + Dash + Dash + Dot + Dash + Dot,
	'+':  Dot + Dash + Dot + Dash + Dot,
	'_':  Dot + Dot + Dash + Dash + Dot + Dash,
}

var MorseCodeToCharMap = map[MorseCode]rune{
	Dot + Dash:                                'A',
	Dash + Dot + Dot + Dot:                    'B',
	Dash + Dot + Dash + Dot:                   'C',
	Dash + Dot + Dot:                          'D',
	Dot:                                       'E',
	Dot + Dot + Dash + Dot:                    'F',
	Dash + Dash + Dot:                         'G',
	Dot + Dot + Dot + Dot:                     'H',
	Dot + Dot:                                 'I',
	Dot + Dash + Dash + Dash:                  'J',
	Dash + Dot + Dash:                         'K',
	Dot + Dash + Dot + Dot:                    'L',
	Dash + Dash:                               'M',
	Dash + Dot:                                'N',
	Dash + Dash + Dash:                        'O',
	Dot + Dash + Dash + Dot:                   'P',
	Dash + Dash + Dot + Dash:                  'Q',
	Dot + Dash + Dot:                          'R',
	Dot + Dot + Dot:                           'S',
	Dash:                                      'T',
	Dot + Dot + Dash:                          'U',
	Dot + Dot + Dot + Dash:                    'V',
	Dot + Dash + Dash:                         'W',
	Dash + Dot + Dot + Dash:                   'X',
	Dash + Dot + Dash + Dash:                  'Y',
	Dash + Dash + Dot + Dot:                   'Z',
	Dot + Dash + Dash + Dash + Dash:           '1',
	Dot + Dot + Dash + Dash + Dash:            '2',
	Dot + Dot + Dot + Dash + Dash:             '3',
	Dot + Dot + Dot + Dot + Dash:              '4',
	Dot + Dot + Dot + Dot + Dot:               '5',
	Dash + Dot + Dot + Dot + Dot:              '6',
	Dash + Dash + Dot + Dot + Dot:             '7',
	Dash + Dash + Dash + Dot + Dot:            '8',
	Dash + Dash + Dash + Dash + Dot:           '9',
	Dash + Dash + Dash + Dash + Dash:          '0',
	Dot + Dash + Dot + Dash + Dot + Dash:      '.',
	Dash + Dash + Dash + Dot + Dot + Dot:      ':',
	Dash + Dash + Dot + Dot + Dash + Dash:     ',',
	Dash + Dot + Dash + Dot + Dash + Dot:      ';',
	Dot + Dot + Dash + Dash + Dot + Dot:       '?',
	Dash + Dot + Dot + Dot + Dash:             '=',
	Dot + Dash + Dash + Dash + Dash + Dot:     '\'',
	Dash + Dot + Dot + Dash + Dot:             '/',
	Dash + Dot + Dash + Dot + Dash + Dash:     '!',
	Dash + Dot + Dot + Dot + Dot + Dash:       '-',
	Dot + Dash + Dot + Dot + Dash + Dot:       '"',
	Dash + Dot + Dash + Dash + Dot:            '(',
	Dash + Dot + Dash + Dash + Dot + Dash:     ')',
	Dot + Dot + Dot + Dash + Dot + Dot + Dash: '$',
	Dot + Dash + Dot + Dot + Dot:              '&',
	Dot + Dash + Dash + Dot + Dash + Dot:      '@',
	Dot + Dash + Dot + Dash + Dot:             '+',
	Dot + Dot + Dash + Dash + Dot + Dash:      '_',
}

func (ms *MorseServer) ParseCommand(stdin string) {
	ms.Options = make(map[string]string)
	isValid := utils.CheckStringWithRegex(strings.ToUpper(stdin), "^[ A-Z0-9.:,;?='/!-\"()$&@+_]*$")
	if !isValid {
		ms.Options["error"] = "ERROR: Illegal input, the supported input is English letters and . :,;? ='/! -\"()$&@+_"
		return
	}
	ms.Options["stdin"] = stdin
}

func (ms *MorseServer) ExecuteCommand(c *gin.Context) {
	if ms.Options["error"] != "" {
		c.JSON(http.StatusOK, gin.H{"response": ms.Options["error"]})
		return
	}
	encodeRes := convertToMorseCode(ms.Options["stdin"])
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("Input: %s"+"<br>"+"Output: %s", ms.Options["stdin"], encodeRes)})
}

func (es *EsromServer) ParseCommand(stdin string) {
	es.Options = make(map[string]string)
	isValid := utils.CheckStringWithRegex(stdin, "^[ .-]*$")
	if !isValid {
		es.Options["error"] = "ERROR: Illegal Morse code input"
		return
	}
	es.Options["stdin"] = stdin
}

func (es *EsromServer) ExecuteCommand(c *gin.Context) {
	if es.Options["error"] != "" {
		c.JSON(http.StatusOK, gin.H{"response": es.Options["error"]})
		return
	}
	decodeRes := convertFromMorseCode(es.Options["stdin"])

	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("Input: %s"+"<br>"+"Output: %s", es.Options["stdin"], decodeRes)})
}

func convertToMorseCode(message string) MorseCode {
	var morseCode MorseCode
	message = strings.ToUpper(message)
	for _, char := range message {
		if code, ok := CharToMorseCodeMap[char]; ok {
			morseCode += code + LetterSpace
		} else if char == ' ' {
			morseCode += WordSpace
		}
	}
	return MorseCode(strings.TrimRight(string(morseCode), " "))
}

func convertFromMorseCode(morseCode string) string {
	var res string
	chars := splitMorseCode(morseCode)
	for _, char := range chars {
		if val, ok := MorseCodeToCharMap[MorseCode(char)]; ok {
			res += string(val)
		} else if char == string(Space) {
			res += " "
		}
	}
	return res
}

func splitMorseCode(morseCode string) []string {
	var chars []string
	words := strings.Split(morseCode, string(WordSpace))
	for _, word := range words {
		letters := strings.Split(word, string(LetterSpace))
		for _, letter := range letters {
			spaceCheck := utils.CheckStringWithRegex(letter, string(Space))
			if spaceCheck {
				deepLetters := strings.Split(letter, string(Space))
				for _, char := range deepLetters {
					chars = append(chars, char)
				}
			}
			chars = append(chars, letter)
		}
		chars = append(chars, " ")
	}
	return chars[:len(chars)-1]
}
