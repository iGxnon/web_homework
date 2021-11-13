package main

type Skill interface {
	Release()
}

type ThunderBreathe struct {
}

func (b *ThunderBreathe) Release() {
	ReleaseSkill("雷之呼吸·壹之型·霹雳一闪", PrefixMap["全集中·常中"])
}

type DiySkill struct {
	Name string
	Func ReleaseFunc
}

func (d *DiySkill) Release() {
	ReleaseSkill(d.Name, d.Func)
}
