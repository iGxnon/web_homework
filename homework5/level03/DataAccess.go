package level03

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// 把 file 封装一下
type DataPool struct {
	path  string
	Cache []UserDao
}

func (p *DataPool) ReloadAll() {
	p.Cache = p.GetDeserializeList()
}

func (p *DataPool) getPwdLock(name string) string {
	for _, entry := range p.Cache {
		if entry.Name == name {
			return entry.Pwd
		}
	}
	panic("获取用户密码失败!")
}

// todo 存到硬盘里
func (p *DataPool) saveToDisk() {

}

func (p *DataPool) ContainsName(name string) bool {
	for _, entry := range p.Cache {
		if entry.Name == name {
			return true
		}
	}
	return false
}

func (p *DataPool) Contains(dao UserDao) bool {
	return p.ContainsName(dao.Name)
}

func (p *DataPool) getIfContainsName(name string) UserDao {
	for _, entry := range p.Cache {
		if entry.Name == name {
			return entry
		}
	}
	panic("不存在这个Name!")
}

func (p *DataPool) ReloadAndContains(dao UserDao) bool {
	p.ReloadAll()
	return p.Contains(dao)
}

func (p *DataPool) GetDeserializeList() []UserDao {
	daoList := make([]UserDao, 0)
	var str = ""
	b := make([]byte, 1024)
	file, _ := os.Open(p.path)
	num, err := file.Read(b)
	for err != io.EOF {
		str = str + string(b[:num])
		num, err = file.Read(b)
	}
	data := strings.Split(str, ";")
	for _, entry := range data {
		dao := UserDao{}
		err := json.Unmarshal([]byte(entry), &dao)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		daoList = append(daoList, dao)
	}
	err3 := file.Close()
	if err3 != nil {
		fmt.Println(err3.Error())
		return nil
	}
	return daoList
}
