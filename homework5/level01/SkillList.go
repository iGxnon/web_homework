package level01

import "errors"

var SkillMap = make(map[string]Skill)

func RegisterDefaultSkill() {
	SkillMap["雷之呼吸·壹之型·霹雳一闪"] = &ThunderBreathe{}
}

func RegisterSkill(skillName, prefix string) (err error) {
	err = nil
	if check(skillName) {
		err = errors.New("你嘴臭你妈呢？")
		return
	}
	if _, ok := PrefixMap[prefix]; !ok {
		err = errors.New("没有叫这个名字的模板哦！")
		return
	}
	SkillMap[skillName] = &DiySkill{
		Name: skillName,
		Func: PrefixMap[prefix],
	}
	return
}
