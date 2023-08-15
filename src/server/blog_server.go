package server

import (
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/server/middleware"
	"WebHome/src/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const articlesPerPage = 10

func BlogListHandler(c *gin.Context) {
	blogsRODao := dao.NewReadOnlyBlogsDao()
	userEntityDao = *dao.NewReadOnlyUserEntityDao()
	res, err := rdb.Get(ctx, "blogCounts").Result()
	var (
		articleCounts        dao.BlogsCounts
		count                int64
		totalPages           int
		articlesProfilesList []dao.BlogProfile
		userIdUsernameMap    map[int64]string
	)
	authorName := c.Query("authorName")
	page, _ := strconv.Atoi(c.Query("page"))
	classification := c.Query("classification")
	title := c.Query("title")
	if err != nil {
		articleCounts = updateBlogsCount(blogsRODao)
	}
	err = json.Unmarshal([]byte(res), &articleCounts)
	if c.Request.Method == http.MethodGet && authorName != "" {
		userId, _ := userEntityDao.SearchUserId(authorName)
		if userId == 0 {
			paginationHTML := generatePaginationHTML(totalPages, page)
			classificationHTML := generateBlogClassificationHTML(model.BlogClassification(classification))
			c.HTML(http.StatusOK, "bloglist.html",
				gin.H{
					"title":              "Blogs",
					"articles":           articlesProfilesList,
					"userIdUsernameMap":  userIdUsernameMap,
					"classificationHTML": template.HTML(classificationHTML),
					"paginationHTML":     template.HTML(paginationHTML),
				})
			return
		}
		userAuth := middleware.GetUserAuth(c)
		var onlyVisible bool
		if userId != userAuth.UserId {
			onlyVisible = true
		} else {
			onlyVisible = false
		}
		count = getBlogCount(articleCounts, userId, classification, onlyVisible)
		articlesProfilesList, userIdUsernameMap = getBlogProfilesList(blogsRODao, userEntityDao, userId, title, classification, page, onlyVisible)
	} else if c.Request.Method == http.MethodGet && authorName == "" {
		count = getBlogCount(articleCounts, 0, classification, false)
		articlesProfilesList, userIdUsernameMap = getBlogProfilesList(blogsRODao, userEntityDao, 0, title, classification, page, false)
	}
	totalPages = int(math.Ceil(float64(count) / float64(articlesPerPage)))
	if page > totalPages {
		page = totalPages
	} else if page == 0 {
		page = 1
	}
	paginationHTML := generatePaginationHTML(totalPages, page)
	classificationHTML := generateBlogClassificationHTML(model.BlogClassification(classification))
	c.HTML(http.StatusOK, "bloglist.html",
		gin.H{
			"title":              "Blogs",
			"articles":           articlesProfilesList,
			"userIdUsernameMap":  userIdUsernameMap,
			"classificationHTML": template.HTML(classificationHTML),
			"paginationHTML":     template.HTML(paginationHTML),
		})
}

func generatePaginationHTML(totalPages, currentPage int) string {
	var buffer bytes.Buffer
	if currentPage > 1 {
		buffer.WriteString(`<a href="blogs?page=`)
		buffer.WriteString(strconv.Itoa(currentPage - 1))
		buffer.WriteString(`">Pre</a>`)
	}

	for i := 1; i <= totalPages; i++ {
		if i == currentPage {
			buffer.WriteString(`<span class="current">`)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(`</span>`)
		} else {
			buffer.WriteString(`<a href="/blogs?page=`)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(`">`)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(`</a>`)
		}
	}

	if currentPage < totalPages {
		buffer.WriteString(`<a href="/blogs?page=`)
		buffer.WriteString(strconv.Itoa(currentPage + 1))
		buffer.WriteString(`">Next</a>`)
	}

	return buffer.String()
}

func generateBlogClassificationHTML(defaultType model.BlogClassification) string {
	var buffer bytes.Buffer
	for i := 1; i < len(blogClassifications); i++ {
		buffer.WriteString(`<option value="`)
		buffer.WriteString(string(blogClassifications[i]))
		if defaultType == blogClassifications[i] {
			buffer.WriteString(`" selected>`)
		} else {
			buffer.WriteString(`">`)
		}
		buffer.WriteString(string(blogClassifications[i]))
		buffer.WriteString(`</option>`)
	}
	return buffer.String()
}

func ReadHandle(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Blog Not Found", "text": "404 - Article does not exist."})
		return
	}
	blogsRODao := dao.NewReadOnlyBlogsDao()
	entity := blogsRODao.GetBlogDetail(title[1:])
	if entity.Id == 0 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Blog Not Found", "text": "404 - Article does not exist."})
		return
	} else if entity.DeletedAt != 0 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Blog Not Found", "text": "404 - Article has been deleted."})
		return
	} else {
		entity.TotalReviews += 1
		entity.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		blogsDao := dao.NewBlogsDao()
		blogsDao.Save(&entity)
		userAuth := middleware.GetUserAuth(c)
		if userAuth.UserId == entity.UserId {
			c.HTML(http.StatusOK, "blogdetail.html", gin.H{"title": entity.Title, "edit": true, "content": string(utils.Base64DecodeString(entity.Content)), "isAnonymous": entity.IsAnonymous})
			return
		} else {
			c.HTML(http.StatusOK, "blogdetail.html", gin.H{"title": entity.Title, "edit": false, "content": string(utils.Base64DecodeString(entity.Content)), "isAnonymous": entity.IsAnonymous})
			return
		}
	}
}

func updateBlogsCount(blogsDao *dao.BlogsDao) dao.BlogsCounts {
	data := blogsDao.Count()
	blogsCount, _ := json.Marshal(data)
	rdb.Set(ctx, "blogCounts", blogsCount, 0)
	return data
}

func getBlogCount(blogsCounts dao.BlogsCounts, userId int64, classification string, onlyVisible bool) int64 {
	var blogCount int64
	switch userId {
	case 0:
		if classification == "" {
			blogCount = blogsCounts.TotalCount
		} else {
			blogCount = blogsCounts.ClassificationCounts[classification]
		}
	default:
		if onlyVisible && classification != "" {
			blogCount = blogsCounts.AuthorVisibleClassificationCounts[strconv.FormatInt(userId, 10)+"-"+classification]
		} else if onlyVisible && classification == "" {
			blogCount = blogsCounts.AuthorVisibleCounts[strconv.FormatInt(userId, 10)]
		} else if !onlyVisible && classification != "" {
			blogCount = blogsCounts.AuthorClassificationCounts[strconv.FormatInt(userId, 10)+"-"+classification]
		} else {
			blogCount = blogsCounts.AuthorCounts[strconv.FormatInt(userId, 10)]
		}
	}
	return blogCount
}

func getBlogProfilesList(blogsRODao *dao.BlogsDao, userRODao dao.UserEntityDao, userId int64, title, classification string, page int, onlyVisible bool) ([]dao.BlogProfile, map[int64]string) {
	blogProfilesList := blogsRODao.GetBlogProfiles(userId, title, classification, page-1, onlyVisible)
	var userIdList []int64
	deDuplicateMap := make(map[int64]bool)
	for _, profile := range blogProfilesList {
		if !deDuplicateMap[profile.UserId] {
			deDuplicateMap[profile.UserId] = true
			userIdList = append(userIdList, profile.UserId)
		}
	}
	userList := userRODao.FindUserListByUserId(userIdList)
	userIdUsernameMap := make(map[int64]string)
	for _, user := range userList {
		userIdUsernameMap[user.Id] = user.Username
	}
	return blogProfilesList, userIdUsernameMap
}

func UploadImageHandle(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := c.FormFile("editormd-image-file")
		guid := c.Query("guid")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error(), "url": ""})
			return
		}
		if file.Size > 2<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "File size cannot exceed 2MB!", "url": ""})
			return
		}
		err = c.SaveUploadedFile(file, fmt.Sprintf("%s%s", blogImagePath, guid))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": 0, "message": err.Error(), "url": ""})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": 1, "message": "", "url": fmt.Sprintf("/img/%s", guid)})
	}()
	wg.Wait()
}

func LoadLocalImageHandle(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		filename := c.Param("guid")
		imagePath := fmt.Sprintf("%s%s", blogImagePath, filename)
		c.File(imagePath)
	}()
	wg.Wait()
}

func EditArticle(c *gin.Context) {
	title := c.Param("title")
	if title != "" {
		userAuth := middleware.GetUserAuth(c)
		blogsDao := dao.NewBlogsDao()
		title = title[1:]
		article := blogsDao.SearchTitle(title)
		if article.UserId == userAuth.UserId {
			classificationHTML := generateBlogClassificationHTML(article.Classification)
			c.HTML(http.StatusOK, "blogedit.html", gin.H{"title": title, "content": string(utils.Base64DecodeString(article.Content)), "isAnonymous": article.IsAnonymous, "classificationHTML": template.HTML(classificationHTML)})
			return
		} else if article.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Can't get the blog information, please check if the blog name is correct."})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This user does not have edit permissions."})
			return
		}
	} else {
		classificationHTML := generateBlogClassificationHTML("")
		c.HTML(http.StatusOK, "blogedit.html", gin.H{"title": "", "content": "", "isAnonymous": false, "classificationHTML": template.HTML(classificationHTML)})
	}
}

type saveRequest struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	Classification string `json:"classification"`
	IsAnonymous    bool   `json:"isAnonymous"`
	PreURL         string `json:"preURL"`
}

func SaveArticle(c *gin.Context) {
	var wg sync.WaitGroup
	blogsDao := dao.NewBlogsDao()
	reqBody, _ := io.ReadAll(c.Request.Body)
	reqMap := &saveRequest{}
	err := json.Unmarshal(reqBody, reqMap)
	title := c.Param("title")[1:]
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"response": "ERROR: No user information was obtained, please login and operate again."})
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Requests that could not be parsed."})
			return
		}

		blogDetail := blogsDao.GetBlogDetail(title)
		if blogDetail.Id == 0 && title != "" {
			blogDetail = *model.NewBlogs()
			blogDetail.UserId = userAuth.UserId
			decodeTitle, _ := url.QueryUnescape(reqMap.Title)
			blogDetail.Title = decodeTitle
			blogDetail.Content = utils.Base64EncodeString([]byte(reqMap.Content))
			blogDetail.Classification = model.BlogClassification(reqMap.Classification)
			blogDetail.IsAnonymous = reqMap.IsAnonymous
			blogsDao.Save(&blogDetail)
			c.JSON(http.StatusOK, gin.H{"response": decodeTitle})
			return
		}
		if blogDetail.UserId != userAuth.UserId {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		decodeTitle, _ := url.QueryUnescape(reqMap.Title)
		blogDetail.Title = decodeTitle
		blogDetail.Content = utils.Base64EncodeString([]byte(reqMap.Content))
		blogDetail.Classification = model.BlogClassification(reqMap.Classification)
		blogDetail.IsAnonymous = reqMap.IsAnonymous
		blogDetail.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		blogsDao.Save(&blogDetail)
		c.JSON(http.StatusOK, gin.H{"response": decodeTitle})
	}()
	wg.Wait()
	updateBlogsCount(blogsDao)
}
