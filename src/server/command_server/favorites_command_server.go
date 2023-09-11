package command_server

import (
	"WebHome/src/database/dao"
	"WebHome/src/server/middleware"
	"WebHome/src/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type operateFavoritesFolderStatus int

const (
	OK          operateFavoritesFolderStatus = 0
	Failed      operateFavoritesFolderStatus = 1
	Exist       operateFavoritesFolderStatus = 2
	NonExistent operateFavoritesFolderStatus = 3
	Ban         operateFavoritesFolderStatus = 4
)

type CdServer struct {
	BaseCommand
}

type LsServer struct {
	BaseCommand
}

type PwdServer struct {
	BaseCommand
}

type MkdirServer struct {
	BaseCommand
}

type RmServer struct {
	BaseCommand
}

type LikeServer struct {
	BaseCommand
}

type MvServer struct {
	BaseCommand
}

func (cs *CdServer) ParseCommand(stdin string) {
	cs.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) > 1 {
		cs.Options["path"] = ""
	} else if len(parts) == 0 {
		cs.Options["path"] = "/"
	} else {
		cs.Options["path"] = parts[0]
	}
}

func (cs *CdServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	path := cs.Options["path"]
	if path == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man cd' command for help."})
		return
	} else if path == "/" {
		userAuth.WorkingPath = path
		c.SetCookie("userAuthorization", utils.SerializationObj(userAuth), 3600, "/", "", true, true)
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
		return
	}
	path = utils.ResolvePath(userAuth.WorkingPath, path)
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	ok := favoritesFolderDao.IsExist(userAuth.UserId, "", path)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("cd: no such file or directory: %s", path)})
		return
	}
	userAuth.WorkingPath = path
	c.SetCookie("userAuthorization", utils.SerializationObj(userAuth), 3600, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"response": "OK!"})
}

func (ls *LsServer) ParseCommand(stdin string) {
	ls.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	switch len(parts) {
	case 0:
		ls.Options["path"] = "\\"
	case 1:
		if strings.ToUpper(parts[0]) == "-L" {
			ls.Options["method"] = "-L"
			ls.Options["path"] = "\\"
		} else {
			ls.Options["path"] = parts[0]
		}
	case 2:
		if strings.ToUpper(parts[0]) == "-L" {
			ls.Options["method"] = "-L"
			ls.Options["path"] = parts[1]
		} else {
			return
		}
	default:
		return
	}
}

func (ls *LsServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	path := ls.Options["path"]
	if path == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man ls' command for help."})
		return
	} else if path == "\\" {
		path = userAuth.WorkingPath
	} else {
		path = utils.ResolvePath(userAuth.WorkingPath, path)
	}
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	ok := favoritesFolderDao.IsExist(userAuth.UserId, "", path)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("ls: %s: No such file or directory", path)})
		return
	}
	res := favoritesFolderDao.GetRecordsByPath(userAuth.UserId, path)
	var response string
	method := ls.Options["method"]
	if method == "" {
		response = generateListHTML(false, res)
	} else {
		response = generateListHTML(true, res)
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func (ps *PwdServer) ParseCommand(stdin string) {
	return
}

func (ps *PwdServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": userAuth.WorkingPath})
}

func (ms *MkdirServer) ParseCommand(stdin string) {
	ms.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) > 1 {
		ms.Options["directoryName"] = ""
	} else if len(parts) == 0 {
		ms.Options["directoryName"] = "\\"
	} else {
		ms.Options["directoryName"] = parts[0]
	}
}

func (ms *MkdirServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	directoryName := ms.Options["directoryName"]
	if directoryName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man mkdir' command for help."})
		return
	} else if directoryName == "\\" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please enter the name and path of the favorite you want to create."})
		return
	}
	path := utils.ResolvePath(userAuth.WorkingPath, directoryName)
	res := createFavoritesFolders(userAuth.UserId, path)
	if res == 0 {
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
	} else if res == 1 {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Favorites already exist"})
	} else {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Favorites already exist"})
	}
}

func (rs *RmServer) ParseCommand(stdin string) {
	rs.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	if len(parts) > 1 {
		rs.Options["directoryName"] = ""
	} else if len(parts) == 0 {
		rs.Options["directoryName"] = "\\"
	} else {
		rs.Options["directoryName"] = parts[0]
	}
}

func (rs *RmServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	directoryName := rs.Options["directoryName"]
	if directoryName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man rm' command for help."})
		return
	} else if directoryName == "\\" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please enter the name and path of the favorite you want to remove."})
		return
	}
	path := utils.ResolvePath(userAuth.WorkingPath, directoryName)
	res := removeFavoritesFolders(userAuth.UserId, path)
	if res == OK {
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
	} else if res == Ban {
		c.JSON(http.StatusOK, gin.H{"response": "WARNING: Disable operation of the directory!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("rm: %s: No such file or directory", directoryName)})
	}
}

func (ls *LikeServer) ParseCommand(stdin string) {
	ls.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	switch len(parts) {
	case 2:
		ls.Options["nickname"] = parts[0]
		ls.Options["path"] = parts[1]
	case 3:
		ls.Options["directory"] = parts[0]
		ls.Options["nickname"] = parts[1]
		ls.Options["path"] = parts[2]
	default:
		return
	}
}

func (ls *LikeServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	nickname := ls.Options["nickname"]
	if nickname == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man like' command for help"})
		return
	}
	var specialSymbols = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	if specialSymbols.FindString(nickname) != "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Do not use special characters for naming."})
		return
	}
	directory := ls.Options["directory"]
	upper := userAuth.WorkingPath
	if directory != "" {
		upper = utils.ResolvePath(upper, directory)
	}
	path := ls.Options["path"]
	creatStatus := createFavoritesFolders(userAuth.UserId, upper)
	if creatStatus == OK || creatStatus == Exist {
		if !creatFavoriteWebPage(userAuth.UserId, nickname, upper, path) {
			c.JSON(http.StatusOK, gin.H{"response": "ERROR: Favorite pages already exist."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "ERROR: Handling exceptions, it is recommended to contact the administrator(x.admin@plac3bo.dev) to check the program problems."})
}

func (ms *MvServer) ParseCommand(stdin string) {
	ms.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	switch len(parts) {
	case 3:
		if strings.ToUpper(parts[0]) == "-F" {
			ms.Options["type"] = "Folder"
			ms.Options["oldName"] = parts[1]
			ms.Options["newName"] = parts[2]
		} else if strings.ToUpper(parts[0]) == "-N" {
			ms.Options["type"] = "Nickname"
			ms.Options["oldName"] = parts[1]
			ms.Options["newName"] = parts[2]
		} else if strings.ToUpper(parts[0]) == "-P" {
			ms.Options["type"] = "Path"
			ms.Options["nickname"] = parts[1]
			ms.Options["path"] = parts[2]
		}
	default:
		return
	}
}

func (ms *MvServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	tp := ms.Options["type"]
	switch tp {
	case "Folder":
		c.JSON(http.StatusOK, gin.H{"response": "Thank you for using, this feature is not yet developed, please look forward to the subsequent development support!"})
		return
	case "Nickname":
		oldName := ms.Options["oldName"]
		newName := ms.Options["newName"]
		if oldName == "." {
			c.JSON(http.StatusOK, gin.H{"response": "WARNING: Prohibit operation of the record."})
			return
		}
		var specialSymbols = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
		if specialSymbols.FindString(newName) != "" {
			c.JSON(http.StatusOK, gin.H{"response": "ERROR: Do not use special characters for naming."})
			return
		}
		if !updatePage(userAuth.UserId, oldName, newName, userAuth.WorkingPath, "") {
			c.JSON(http.StatusOK, gin.H{"response": "ERROR: Failed to update or favorite page does not exist.(Currently only supports changes under the current favorites, please look forward to subsequent updates!)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
	case "Path":
		nickname := ms.Options["nickname"]
		if nickname == "." {
			c.JSON(http.StatusOK, gin.H{"response": "WARNING: Prohibit operation of the record."})
			return
		}
		path := ms.Options["path"]
		if !updatePage(userAuth.UserId, nickname, "", userAuth.WorkingPath, path) {
			c.JSON(http.StatusOK, gin.H{"response": "ERROR: Failed to update or favorite page does not exist.(Currently only supports changes under the current favorites, please look forward to subsequent updates!)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": "OK!"})
	default:
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported inputs, use the 'man mv' command for help."})
	}
}

func generateListHTML(detail bool, records []dao.FavoritesProfile) string {
	var buffer bytes.Buffer

	if detail {
		buffer.WriteString(`<table><tbody>`)
		for _, record := range records {
			isFolder := record.IsFolder
			nickName := record.Nickname
			path := record.Path
			if isFolder {
				buffer.WriteString(`<tr><td><a href="#" class="disabled-link">/`)
				buffer.WriteString(nickName)
				buffer.WriteString(`</a>`)
			} else {
				buffer.WriteString(`<tr><td><a href="`)
				buffer.WriteString(path)
				buffer.WriteString(`" class="enabled-link" target="_blank">`)
				buffer.WriteString(nickName)
				buffer.WriteString(`</a>`)
			}
			buffer.WriteString(`</td><td><span>`)
			buffer.WriteString(path)
			buffer.WriteString(`</span></td></tr>`)
		}
		buffer.WriteString(`</tbody></table>`)
	} else {
		buffer.WriteString(`<div id="favorites_list">`)
		for _, record := range records {
			isFolder := record.IsFolder
			nickName := record.Nickname
			path := record.Path
			if nickName == "." {
				continue
			}
			if isFolder {
				buffer.WriteString(`<a href="#" class="disabled-link">/`)
			} else {
				buffer.WriteString(`<a href="`)
				buffer.WriteString(path)
				buffer.WriteString(`" class="enabled-link" target="_blank">`)
			}
			buffer.WriteString(nickName)
			buffer.WriteString(`</a>`)
		}
		buffer.WriteString(`</div>`)
	}

	return buffer.String()
}

func createFavoritesFolders(userId int64, path string) operateFavoritesFolderStatus {
	res := OK
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	pathParts := utils.GetPathParts(path)
	upper := "/"
	for _, part := range pathParts {
		upperPath := filepath.Join(upper, part)
		ok := favoritesFolderDao.IsExist(userId, ".", upperPath)
		if !ok {
			success := favoritesFolderDao.CreatRecord(userId, true, part, upper, upper)
			if !success {
				res = Failed
				return res
			}
			if !favoritesFolderDao.IsExist(userId, ".", upper) {
				favoritesFolderDao.CreatRecord(userId, false, ".", upper, upper)
			}
		}
		upper = filepath.Join(upper, part)
	}
	if favoritesFolderDao.IsExist(userId, ".", upper) {
		res = Exist
		return res
	}
	favoritesFolderDao.CreatRecord(userId, false, ".", upper, upper)
	return res
}

func removeFavoritesFolders(userId int64, path string) operateFavoritesFolderStatus {
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	dir, file := filepath.Split(path)
	dir, _ = filepath.Abs(dir)
	if file == "" && dir == "/" {
		return Ban
	}
	if favoritesFolderDao.IsExist(userId, ".", path) {
		favoritesFolderDao.RemoveRecords(userId, file, dir)
		favoritesFolderDao.RemoveRecords(userId, "", path)
		favoritesFolderDao.RemoveRecords(userId, "", path+"/%")
		return OK
	}
	res := favoritesFolderDao.GetRecordByNicknameAndUpper(userId, file, dir)
	if res.Id != 0 {
		favoritesFolderDao.Delete(&res)
	} else {
		return NonExistent
	}
	return OK
}

func creatFavoriteWebPage(userId int64, nickname, upper, path string) bool {
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	if favoritesFolderDao.CreatRecord(userId, false, nickname, upper, path) {
		return true
	}
	return false
}

func updatePage(userId int64, oldName, newName, upper, path string) bool {
	favoritesFolderDao := dao.NewFavoritesFolderDao()
	record := favoritesFolderDao.GetRecordByNicknameAndUpper(userId, oldName, upper)
	if record.Id == 0 {
		return false
	}
	if path == "" {
		record.Nickname = newName
	} else {
		record.Path = path
	}
	record.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
	favoritesFolderDao.Save(&record)
	return true
}
