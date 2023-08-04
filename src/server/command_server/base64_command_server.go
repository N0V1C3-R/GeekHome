package command_server

import (
	"WebHome/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Base64Server struct {
	BaseCommand
}

func (bs *Base64Server) ParseCommand(stdin string) {
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) == 0 {
		bs.Options = make(map[string]string)
		bs.Options["warn"] = "WARRING: Please enter the text to be encoded or decoded."
		return
	}
	bs.Required = make(map[string]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if bs.Required["text"] == "" {
			switch strings.ToUpper(arg) {
			case "-D", "-DECODE":
				bs.Required["mode"] = "DECODE"
				continue
			case "-E", "-ENCODE":
				bs.Required["mode"] = "ENCODE"
				continue
			}
		}
		bs.Required["text"] += arg + " "
	}
	bs.Required["text"] = strings.TrimRight(bs.Required["text"], " ")
	if bs.Required["mode"] == "" {
		bs.Required["mode"] = "ENCODE"
	}
}

func (bs *Base64Server) ExecuteCommand(c *gin.Context) {
	if bs.Options["warn"] != "" {
		c.JSON(http.StatusBadRequest, gin.H{"response": bs.Options["warn"]})
		return
	}

	mode := bs.Required["mode"]
	modeMap := map[string]interface{}{
		"ENCODE": base64Encode,
		"DECODE": base64Decode,
	}
	function := modeMap[mode]
	res := utils.CallFunction(function, bs.Required["text"])
	c.JSON(http.StatusOK, gin.H{"response": res})
}

func base64Encode(text string) string {
	encodeRes := utils.Base64EncodeString([]byte(text))
	return encodeRes
}

func base64Decode(text string) string {
	if !utils.IsBase64String(text) {
		return "ERROR: Illegal Base64-encoded string"
	}
	decodeRes := utils.Base64DecodeString(text)
	return string(decodeRes)
}
