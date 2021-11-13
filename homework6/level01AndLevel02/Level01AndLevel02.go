package main

import (
	"github.com/gin-gonic/gin"
)

const MaxAge = 3600

// 凭自己的想法写的，可能会与和已有的设计方式不一样

func main() {
	r := gin.Default()

	// 定义几个服务分组
	loadInfoPage(r)
	loadPaperPage(r)
	loadStoragePage(r)

	// 把文件读入内存
	LoadALlUserDao()
	LoadAllPapers()

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if dao, ok := UserDaoMap[username]; ok {
			if password != dao.PwdLock {
				c.JSON(403, gin.H{
					"message": "密码错误",
				})
				return
			}

			cookieStr := GenerateRandomSid()
			c.SetCookie("gin_cookie", cookieStr, MaxAge, "/", "", false, true)
			session := &Session{
				sid:    cookieStr,
				maxAge: MaxAge,
				id:     username,
			}
			PutSessionIfAbsence(session)

			c.JSON(200, gin.H{
				"message": "登录成功!",
			})
		} else {
			c.JSON(403, gin.H{
				"message": "不存在该用户",
			})
		}
	})
	r.Run()
}

// 中间件=>验证Session
func auth(c *gin.Context) {

	cookie, err := c.Cookie("gin_cookie")
	if err != nil {
		c.JSON(403, gin.H{
			"message": "认证失败,没有cookie",
		})
		c.Abort()
		return
	}
	session, ok := SessionsMap[cookie]
	if !ok {
		c.JSON(403, gin.H{
			"message": "认证失败,cookie无效!",
		})
		c.Abort()
		return
	}
	c.Set("cookie", cookie)
	c.Set("session", session)
	dao, ok := UserDaoMap[session.id]
	if !ok {
		c.JSON(403, gin.H{
			"message": "认证失败,无该用户!",
		})
		c.Abort()
		return
	}
	user := &User{
		Name:    dao.Name,
		Age:     dao.Age,
		Gender:  dao.Gender,
		NpyName: dao.NpyName,
	}
	c.Set("user", user)
	c.Next()
	_, exists := c.Get("abruptError")
	if exists {
		c.JSON(403, gin.H{
			"message": "突发恶疾",
		})
		c.Abort()
	}
}

// 提供个人账户详细信息服务和别人的账户简略信息服务
func loadInfoPage(r *gin.Engine) {
	infoPage := r.Group("/info")
	infoPage.GET("", auth, func(c *gin.Context) {
		user, _ := c.Get("user")
		queryName := c.DefaultQuery("name", user.(*User).Name)
		dao, ok := UserDaoMap[queryName]
		if !ok {
			c.JSON(403, gin.H{
				"message": "不存在这个用户",
			})
			return
		}
		if queryName == user.(*User).Name {
			c.JSON(200, gin.H{
				"Name":   user.(*User).Name,
				"Age":    user.(*User).Age,
				"Gender": user.(*User).Gender,
				"Npy":    user.(*User).NpyName,
			})
		} else {

			c.JSON(200, gin.H{
				"Name":   dao.Name,
				"Age":    dao.Age,
				"Gender": dao.Gender,
			})
		}
	})
}

// 提供文章获取服务
func loadPaperPage(r *gin.Engine) {
	paperPage := r.Group("/paper")
	paperPage.GET("/:path", auth, func(c *gin.Context) {
		path := c.Param("path")
		if paper, ok := PaperMap[path]; ok {
			c.JSON(200, gin.H{
				"title":   paper.Name,
				"content": paper.Content,
			})
		} else {
			c.JSON(404, gin.H{
				"message": "没有找到你要的文章哦",
			})
		}
	})
}

// 提供文件存储读取服务
func loadStoragePage(r *gin.Engine) {
	r.MaxMultipartMemory = 8 << 20
	r.Static("/", "./public")
	storagePage := r.Group("/storage")
	storagePage.POST("/put", auth, func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.Set("abruptError", err)
		}
		err1 := c.SaveUploadedFile(file, file.Filename)
		if err1 != nil {
			c.Set("abruptError", err)
		}
		c.String(200, "uploaded successfully")
	})
	storagePage.GET("/get", auth, func(c *gin.Context) {

	})
}
