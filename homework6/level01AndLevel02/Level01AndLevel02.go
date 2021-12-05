package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const MaxAge = 3600

// 凭自己的想法写的，可能会与和已有的设计方式不一样

func main() {

	r := gin.Default()

	// 把文件读入内存
	LoadALlUserDao()
	LoadAllPapers()

	// 开启http服务
	startServer(r)

	err := r.Run()
	if err != nil {
		return
	}
}

func startServer(r *gin.Engine) {
	// 登录注册
	r.POST("/login", loginHandle)

	// 提供个人账户详细信息服务和别人的账户简略信息服务
	// todo 为infoPage添加更多子服务
	infoPage := r.Group("/info")
	infoPage.GET("", authHandle, infoProvideHandle)

	paperPage := r.Group("/paper")
	paperPage.GET("/:path", authHandle, paperProvideHandle)

	// 文件大小限制和保存位置
	r.MaxMultipartMemory = 8 << 20
	r.Static("/homework6/level01AndLevel02/", "./public")

	// 提供文件存储提取服务
	// todo 为storagePage添加更多子服务
	storagePage := r.Group("/storage")
	storagePage.POST("/put", authHandle, storagePutHandle)
	storagePage.GET("/get", authHandle, storageGetHandle)
}

func authHandle(c *gin.Context) {
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

func loginHandle(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password") // 假装这是客户端经过加密后传输给服务端然后解密的结果（
	if dao, ok := UserDaoMap[username]; ok {
		// md5加密后与数据池内密码对比
		h := md5.New()
		h.Write([]byte(password))
		if hex.EncodeToString(h.Sum(nil)) != dao.PwdLock {
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
		return
	}
	c.JSON(403, gin.H{
		"message": "不存在该用户",
	})
}

func infoProvideHandle(c *gin.Context) {
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
		return
	}

	c.JSON(200, gin.H{
		"Name":   dao.Name,
		"Age":    dao.Age,
		"Gender": dao.Gender,
	})

}

func paperProvideHandle(c *gin.Context) {
	path := c.Param("path")
	if paper, ok := PaperMap[path]; ok {
		c.JSON(200, gin.H{
			"title":   paper.Name,
			"content": paper.Content,
		})
		return
	}
	c.JSON(404, gin.H{
		"message": "没有找到你要的文章哦",
	})
}

func storagePutHandle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Set("abruptError", err) // 测试中间件
		return
	}
	err = c.SaveUploadedFile(file, "./homework6/level01AndLevel02/files/"+file.Filename)
	if err != nil {
		c.Set("abruptError", err) // 测试中间件
		return
	}
	c.String(200, "uploaded successfully")
}

func storageGetHandle(c *gin.Context) {
	fileDir := c.Query("fileDir")
	fileName := c.Query("fileName")

	_, err := os.Open("./homework6/level01AndLevel02/files/" + fileDir + "/" + fileName)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "资源不存在",
		})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("./homework6/level01AndLevel02/files/" + fileDir + "/" + fileName)
}
