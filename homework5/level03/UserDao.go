package main

// 还是无法抛弃 Java 那套思想
type UserDao struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd_lock"`
}
