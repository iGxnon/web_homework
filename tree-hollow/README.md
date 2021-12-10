## Web后端部分

### 技术细节

+ 数据库设计

comment 表

> 储存秘密的评论以及评论的评论
>
> parent_id --> 父节点id (可以是comment，也可以是secret)
>
> comment_type --> 类型 (代表 comment 或者 secret 的评论)

```sql
-- auto-generated definition
create table comment
(
    id           bigint auto_increment
        primary key,
    parent_id    bigint      not null,
    comment_type tinyint(1)  not null,
    content      text        not null,
    snitch_name  varchar(20) not null,
    is_open      tinyint(1)  not null,
    comment_time datetime        not null,
    update_time  datetime        not null
);
```

secret 表

> 储存秘密的表单

```sql
-- auto-generated definition
create table secret
(
    id          bigint auto_increment
        primary key,
    comment_cnt int         default 0  not null,
    content     text                   not null,
    snitch_name varchar(20) default '' not null,
    is_open     tinyint(1)  default 0  not null,
    post_time   datetime                   not null,
    update_time datetime                   not null
);
```

snitch 表

> 储存用户信息(snitch)

```sql
-- auto-generated definition
create table snitch
(
    id       bigint auto_increment,
    name     varchar(20) not null,
    password text        not null,
    primary key (id, name)
);
```

tree_hollow 表

> 储存**树洞**实体和**Form**表单中的信息，与Nukkit相关

```sql
-- auto-generated definition
create table tree_hollow
(
    id                   bigint auto_increment,
    prefix               varchar(20) default '' not null,
    model_json           mediumtext             not null,
    model_serialized_img mediumtext             not null,
    loc                  tinytext               not null,
    form_title           varchar(20) default '' not null,
    form_texts           tinytext               not null,
    animation            tinytext               null,
    particle             tinytext               null,
    primary key (id, prefix)
);
```

+ jwt 鉴权

> 根据网上查到的jwt基本原理手动造了个轮子（主要是不允许用除gin以外第三方库
>
> 两个 token，都用来负责签名授权，其中一个时效比另一个长，可以作为短期间隔登录刷新一组新token对 的 refreshtoken

```go
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
	"tree-hollow/config"
)

// 自己简单地造了个轮子
// 签名算法使用的是 HMAC SHA256

const (
	HeaderPlain = "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
)

// Claims Payload
type Claims struct {
	InterArrivalTime int64  `json:"iat"` // 到达时间
	ExpirationDate   int64  `json:"exp"` // 认证时间
	UserName         string `json:"user_name"`
}

func GenerateTokenPairWithUserName(userName string) (token, refreshToken string, err error) {
	now := time.Now().Unix()
	tokenClaims := Claims{
		InterArrivalTime: now,
		ExpirationDate:   config.Config.JWTTimeOut,
		UserName:         userName,
	}

	// refreshToken时间会比token要长
	refreshTokenClaim := Claims{
		InterArrivalTime: now,
		ExpirationDate:   config.Config.JWTTimeOut * 10,
		UserName:         userName,
	}

	token, err = generateByClaims(tokenClaims)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = generateByClaims(refreshTokenClaim)
	if err != nil {
		return "", "", err
	}
	return
}

func generateByClaims(claims Claims) (string, error) {
	bytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.StdEncoding.EncodeToString([]byte(HeaderPlain))
	payload := base64.StdEncoding.EncodeToString(bytes)

	return SignJWT(header, payload), nil
}

func SignJWT(header, payload string) string {

	hash := hmac.New(sha256.New, []byte(config.Config.JWTKey))
	hash.Write([]byte(header + payload))

	signed := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(hash.Sum(nil))))

	return header + "." + payload + "." + signed
}

func AuthorizeJWT(jwtStr string) (bool, error) {
	claims := Claims{}

	parts := strings.Split(jwtStr, ".")

	payload, _ := base64.StdEncoding.DecodeString(parts[1])
	signed := parts[2]
	dSigned := strings.Split(SignJWT(parts[0], parts[1]), ".")[2]
	if signed != dSigned {
		return false, nil
	}
	err := json.Unmarshal(payload, &claims)
	if err != nil {
		return false, err
	}
	now := time.Now().Unix()

	// 如果现在的时间减去上一次登录时间大于认证时间
	if claims.ExpirationDate < now-claims.InterArrivalTime {
		return false, nil
	}
	return true, nil
}

```

+ 加盐Hash
> 官方有现成的库
> 
> "golang.org/x/crypto/bcrypt"

# TreeHollow (建设中) 

---

> "在从前，当一个人心里有不可告人的秘密时，他会跑到深山里，
> 找一棵树，在树上挖一个洞，将秘密告诉那个树洞，然后再用泥土封起来，
> 这样这个秘密就永远没有人知道了。"
> ————《2046》

# 缘起

---
缘起于一次网校后端培训的作业：做一个登录评论的后端程序，结合《2046》的经典桥段，便诞生了这么一个项目

# 介绍

---

![](https://gitee.com/iGxnon/image-host/raw/master/pic-go/TreeHollow_notransparent.png)

# 功能

---

## Nukkit 插件

也是主打部分，作为一个插件给玩家提供"树洞服务"

## 树洞网页

emmm，跨平台！

## 树洞APP

市面上类似app好像非常多，那肯定竞争不过别人的

# 其他内容

---

[TheWorldTreeHollow 插件介绍]()

# 结语

---

> Hope you enjoy.
