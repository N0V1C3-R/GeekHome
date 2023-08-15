package command_server

import (
	"WebHome/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type GenpwdServer struct {
	BaseCommand
}

func (gs *GenpwdServer) ParseCommand(stdin string) {
	gs.Options = make(map[string]string)
	gs.Options["complexities"] = "15"
	gs.Options["passwordLength"] = "12"
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) == 0 {
		return
	}
	var length string
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if strings.HasPrefix(arg, "-") {
			count := 0
			for _, char := range arg {
				switch strings.ToUpper(string(char)) {
				case "L":
					count += 8
				case "U":
					count += 4
				case "N":
					count += 2
				case "S":
					count += 1
				case "-":
					continue
				default:
					gs.Options["Error"] = "ERROR: flag provided but not defined: -" + string(char) + "."
					return
				}
			}
			gs.Options["complexities"] = strconv.Itoa(count)
			continue
		}
		length += arg
	}
	if utils.IsNumeric(length) {
		lengthVal, _ := strconv.ParseInt(length, 10, 0)
		if lengthVal < 6 {
			gs.Options["Error"] = "ERROR: Password length must not be less than 6 digits."
			return
		} else if lengthVal > 99 {
			gs.Options["Error"] = "ERROR: Password length must not be more than 99 digits."
			return
		}
		gs.Options["passwordLength"] = length
	} else if length == "" {
		return
	} else {
		gs.Options["Error"] = "ERROR: Unrecognized passwordLength parameter."
		return
	}
}

func (gs *GenpwdServer) ExecuteCommand(c *gin.Context) {
	if gs.Options["Error"] != "" {
		c.JSON(http.StatusOK, gin.H{"response": gs.Options["Error"]})
		return
	}
	complexitiesVal, _ := strconv.ParseInt(gs.Options["complexities"], 10, 0)
	passwordLengthVal, _ := strconv.ParseInt(gs.Options["passwordLength"], 10, 0)
	password := utils.GeneratePassword(int(complexitiesVal), int(passwordLengthVal))
	c.JSON(http.StatusOK, gin.H{"response": password})
}
