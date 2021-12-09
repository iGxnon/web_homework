package dao

import (
	"encoding/json"
	"tree-hollow/model"
)

const (
	DefaultIndexTitle         = "${Prefix}"
	DefaultIndexContent       = "'在从前，当一个人心里有不可告人的秘密时，他会跑到深山里，找一棵树，在树上挖一个洞，将秘密告诉那个树洞，然后再用泥土封起来，这样这个秘密就永远没有人知道了。'\n\n 当然，在这里别人能够聆听到你的秘密哦，但他们不知道你是谁，不过欢迎你打开那个象征着勇气的'公开'开关！"
	DefaultIndexBtnStrFirst   = "聆听树洞里的秘密"
	DefaultIndexBtnStrSecond  = "说出心里的秘密"
	DefaultSecretBtnStrFist   = "聆听下一个秘密"
	DefaultSecretBtnStrSecond = "返回树洞"
	DefaultSaySwitchStr       = "证明'勇气'的公开开关！"
	StrCnt                    = 7
)

func InsertTreeHollow(treeHollow model.TreeHollow) error {
	sqlStr := "INSERT INTO tree_hollow(prefix, model_json, model_serialized_img, loc, form_title, form_texts, animation, particle) values(?, ?, ?, ?);"
	loc, err := json.Marshal(treeHollow.Loc)
	if err != nil {
		return err
	}
	formTexts, err := json.Marshal(treeHollow.FormTexts)
	if err != nil {
		return err
	}
	_, err = dB.Exec(sqlStr, treeHollow.Prefix, treeHollow.ModelJson, treeHollow.ModelSerializedIMG, string(loc), treeHollow.FormTitle, string(formTexts), treeHollow.Animation, treeHollow.Particle)
	return err
}

func UpdateFormTexts(id int, texts model.FormTexts) error {
	sqlStr := "UPDATE tree_hollows SET form_texts=? WHERE id = ?;"
	bytes, err := json.Marshal(texts)
	if err != nil {
		return err
	}
	ret, err := dB.Exec(sqlStr, string(bytes), id)
	if err != nil {
		return err
	}
	_, err = ret.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func SelectFormTextsById(id int) (model.FormTexts, error) {

	var formTexts model.FormTexts

	row := dB.QueryRow("SELECT form_texts FROM tree_hollows WHERE id = ?;", id)
	if row.Err() != nil {
		return formTexts, row.Err()
	}

	var serializedFormTexts string

	err := row.Scan(&serializedFormTexts)
	if err != nil {
		return formTexts, err
	}

	err = json.Unmarshal([]byte(serializedFormTexts), &formTexts)
	if err != nil {
		return formTexts, err
	}

	return formTexts, nil
}
