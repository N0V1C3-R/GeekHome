package command_server

import (
	"WebHome/src/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"os"
	"strings"
)

type HelpServer struct {
	BaseCommand
}

type ManServer struct {
	BaseCommand
}

type helpFileStruct []commandInfo

type commandInfo struct {
	CommandName string     `json:"commandName"`
	AliasNames  []string   `json:"aliasNames"`
	Profile     string     `json:"profile"`
	Detail      detailInfo `json:"detail"`
}

type detailInfo struct {
	Usage    string   `json:"usage"`
	Desc     string   `json:"desc"`
	Options  options  `json:"options"`
	Examples examples `json:"examples"`
	Notes    notes    `json:"notes"`
}

type options map[string][]string

type examples []string

type notes []string

func (hs *HelpServer) ParseCommand(stdin string) {
	return
}

func (hs *HelpServer) ExecuteCommand(c *gin.Context) {
	file, _ := os.Open("../../../help.json")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	var helpFileStruct helpFileStruct
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&helpFileStruct)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "command not found: help"})
		return
	}
	response := generateCommandInfoHTML(helpFileStruct)
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func (ms *ManServer) ParseCommand(stdin string) {
	ms.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) != 1 {
		ms.Options["Error"] = "ERROR: Please enter the correct command."
		return
	}
	ms.Options["command"] = parts[0]
}

func (ms *ManServer) ExecuteCommand(c *gin.Context) {
	if ms.Options["Error"] != "" {
		c.JSON(http.StatusOK, gin.H{"response": ms.Options["Error"]})
		return
	}
	file, _ := os.Open("../../../help.json")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	var helpFileStruct helpFileStruct
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&helpFileStruct)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "command not found: man"})
		return
	}
	command := ms.Options["command"]
	for _, commandInfo := range helpFileStruct {
		if commandInfo.CommandName == command {
			response := generateCommandDetailHTML(commandInfo)
			c.JSON(http.StatusOK, gin.H{"response": response})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("No manual entry for %s", command)})
}

func generateCommandInfoHTML(helpFileStruct helpFileStruct) string {
	var buffer bytes.Buffer
	buffer.WriteString(`<table><tbody>`)
	for _, commandInfo := range helpFileStruct {
		commandName := commandInfo.CommandName
		commandProfile := commandInfo.Profile
		buffer.WriteString(`<tr><td><span>`)
		buffer.WriteString(commandName)
		buffer.WriteString(`</span></td><td><span>`)
		buffer.WriteString(commandProfile)
		buffer.WriteString(`</span></td></tr>`)
	}
	buffer.WriteString(`</tbody></table>`)
	return buffer.String()
}

func generateCommandDetailHTML(info commandInfo) string {
	var buffer bytes.Buffer
	buffer.WriteString(`<div class="manBox">`)
	buffer.WriteString(`
	<div class="divider">
		<span class="divider-line"></span>
		<span class="divider-text">General Commands Manual</span>
		<span class="divider-line"></span>
	</div>
`)

	commandName := info.CommandName
	profile := info.Profile
	buffer.WriteString(`<div class="title">NAME</div>`)
	buffer.WriteString(`<div class="commandName">`)
	buffer.WriteString(commandName + " - " + profile)
	buffer.WriteString(`</div>`)

	aliasNames := info.AliasNames
	if len(aliasNames) != 0 {
		buffer.WriteString(`<div class="title">ALIAS</div>`)
		buffer.WriteString(`<div class="aliasList">`)
		buffer.WriteString(`<span>`)
		for i := 0; i < len(aliasNames); i++ {
			buffer.WriteString(aliasNames[i])
			if i != len(aliasNames)-1 {
				buffer.WriteString(`,&nbsp;`)
			}
		}
		buffer.WriteString(`</span>`)
		buffer.WriteString(`</div>`)
	}

	usage := info.Detail.Usage
	usage = strings.ReplaceAll(usage, "<", "&lt;")
	usage = strings.ReplaceAll(usage, ">", "&gt;")
	buffer.WriteString(`<div class="title">USAGE</div>`)
	buffer.WriteString(`<div class="usage">`)
	buffer.WriteString(`<span>`)
	buffer.WriteString(usage)
	buffer.WriteString(`</span>`)
	buffer.WriteString(`</div>`)

	desc := info.Detail.Desc
	buffer.WriteString(`<div class="title">DESC</div>`)
	buffer.WriteString(`<div class="desc">`)
	buffer.WriteString(`<span>`)
	buffer.WriteString(desc)
	buffer.WriteString(`</span>`)
	buffer.WriteString(`</div>`)

	options := info.Detail.Options
	if len(options) > 0 {
		buffer.WriteString(`<div class="title">OPTIONS</div>`)
		buffer.WriteString(`<div class="options">`)
		buffer.WriteString(`<table>`)
		for option, instructions := range options {
			buffer.WriteString(`<tr>`)
			buffer.WriteString(`<td class="optionTable">`)
			buffer.WriteString(`<span>`)
			buffer.WriteString(option)
			buffer.WriteString(`</span>`)
			buffer.WriteString(`<br>`)
			buffer.WriteString(`<span class="optionality">`)
			buffer.WriteString(instructions[0])
			buffer.WriteString(`</span>`)
			buffer.WriteString(`</td>`)
			buffer.WriteString(`<td class="optionTable">`)
			for i := 1; i < len(instructions); i++ {
				instruction := instructions[i]
				buffer.WriteString(`<span class="instruction">`)
				buffer.WriteString(instruction)
				buffer.WriteString(`</span>`)
				buffer.WriteString(`<br>`)
			}
			buffer.WriteString(`</td>`)
			buffer.WriteString(`</tr>`)
		}
		buffer.WriteString(`</table>`)
		buffer.WriteString(`</div>`)
	}

	examples := info.Detail.Examples
	if len(examples) > 0 {
		buffer.WriteString(`<div class="title">EXAMPLE</div>`)
		buffer.WriteString(`<div class="examples">`)
		for _, example := range examples {
			buffer.WriteString(`<pre><code>`)
			buffer.WriteString(example)
			buffer.WriteString(`</pre></code>`)
		}
		buffer.WriteString(`</div>`)
	}

	notes := info.Detail.Notes
	if len(notes) > 0 {
		buffer.WriteString(`<div class="title">NOTE</div>`)
		buffer.WriteString(`<ul class="notes">`)
		for _, note := range notes {
			buffer.WriteString(`<li>`)
			buffer.WriteString(note)
			buffer.WriteString(`</li>`)
		}
		buffer.WriteString(`</ul>`)
	}

	buffer.WriteString(`
	<div class="divider">
		<span class="divider-line"></span>
		<span class="divider-text">END</span>
		<span class="divider-line"></span>
	</div>
`)

	buffer.WriteString(`</div>`)

	return buffer.String()
}
