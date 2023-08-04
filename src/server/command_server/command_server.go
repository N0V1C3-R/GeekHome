package command_server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func CommandHandle(c *gin.Context) {
	reqBody, _ := io.ReadAll(c.Request.Body)
	reqMap := &CommandRequest{}
	err := json.Unmarshal(reqBody, reqMap)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(reqMap.Stdin, " ")
	cmd := parts[0]

	commandServer, ok := commands[strings.ToLower(cmd)]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"response": fmt.Sprintf("command not found: %s", cmd)})
		return
	}
	commandEntity := reflect.New(commandServer.Elem())
	service := commandEntity.Interface().(CommandImpl)

	stdin := reqMap.Stdin[len(cmd):]
	stdin = strings.TrimLeft(stdin, " ")
	service.ParseCommand(stdin)
	service.ExecuteCommand(c)
}
