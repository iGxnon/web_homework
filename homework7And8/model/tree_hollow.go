package model

// TreeHollow
// remain for nukkit as an entity
type TreeHollow struct {
	Id                 int       `json:"id"`
	Prefix             string    `json:"prefix"`               // 实体名称
	ModelJson          string    `json:"model_json"`           // 模型文件
	ModelSerializedIMG string    `json:"model_serialized_img"` // 序列化的模型贴图
	Loc                Location  `json:"loc"`                  // 实体位置
	FormTitle          string    `json:"form_title"`           // Form标题
	FormTexts          FormTexts `json:"form_texts"`           // Form内容
	Animation          string    `json:"animation"`            // TODO 实体动画
	Particle           string    `json:"particle"`             // TODO 粒子效果
}

type FormTexts struct {
	IndexTitle         string `json:"index_title"`
	IndexContent       string `json:"index_content"`
	IndexBtnStrFirst   string `json:"index_btn_str_first"`
	IndexBtnStrSecond  string `json:"index_btn_str_second"`
	SecretBtnStrFist   string `json:"secret_btn_str_fist"`
	SecretBtnStrSecond string `json:"secret_btn_str_second"`
	SaySwitchStr       string `json:"say_switch_str"`
}

type Location struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Yaw       int     `json:"yaw"`
	LevelName string  `json:"level_name"`
}
