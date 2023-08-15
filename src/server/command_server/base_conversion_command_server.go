package command_server

import (
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type BaseConversionServer struct {
	BaseCommand
}

func (bcs *BaseConversionServer) ParseCommand(stdin string) {
	bcs.Options = make(map[string]string)
	bcs.Options["sourceBase"] = "10"
	bcs.Options["targetBase"] = "2"
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	var value string
	if len(parts) == 0 {
		bcs.Options["Error"] = "ERROR: Missing values that need to be converted."
		return
	}
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		switch strings.ToUpper(arg) {
		case "-S":
			if i+1 < len(parts) {
				i++
				bcs.Options["sourceBase"] = parts[i]
				continue
			} else {
				bcs.Options["Error"] = "ERROR: Unable to parse the " + arg + "source parameters."
				return
			}
		case "-T":
			if i+1 < len(parts) {
				i++
				bcs.Options["targetBase"] = parts[i]
				continue
			} else {
				bcs.Options["Error"] = "ERROR: Unable to parse the " + arg + "target parameter."
				return
			}
		}
		numStr := strings.ReplaceAll(arg, ",", "")
		value += numStr
	}
	if value == "" {
		return
	}
	valList := strings.Split(value, ".")
	switch len(valList) {
	case 0:
		bcs.Options["Error"] = "ERROR: Missing values that need to be converted."
	case 1:
		intPart := valList[0]
		if strings.HasPrefix(intPart, "+") || strings.HasPrefix(intPart, "-") {
			if utils.IsNumeric(intPart[1:]) {
				bcs.Options["intPart"] = intPart
				return
			}
		} else {
			if utils.IsNumeric(intPart) {
				bcs.Options["intPart"] = intPart
				return
			}
		}
		bcs.Options["Error"] = "ERROR: Numeric parsing exception."
	case 2:
		bcs.Options["Error"] = "ERROR: Floating-point data conversion is not supported for the time being."
		return
	default:
		bcs.Options["Error"] = "ERROR: Numeric parsing exception."
	}
}

func (bcs *BaseConversionServer) ExecuteCommand(c *gin.Context) {
	if bcs.Options["Error"] != "" {
		c.JSON(http.StatusOK, gin.H{"response": bcs.Options["Error"]})
		return
	}
	sourceBase := bcs.Options["sourceBase"]
	sourceBaseVal, err := strconv.ParseInt(sourceBase, 10, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("ERROR: Unresolvable -s values: %s", sourceBase)})
		return
	}
	targetBase := bcs.Options["targetBase"]
	targetBaseVal, err := strconv.ParseInt(targetBase, 10, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("ERROR: Unresolvable -t values: %s", targetBase)})
		return
	}
	if sourceBaseVal < 2 || sourceBaseVal > 36 || targetBaseVal < 2 || targetBaseVal > 36 {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Progressive conversion intervals of 2-36."})
		return
	}
	intPart := bcs.Options["intPart"]
	intPartRes, err := utils.ConvertWithSign(intPart, int(sourceBaseVal), int(targetBaseVal))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "Error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": intPartRes})
}
