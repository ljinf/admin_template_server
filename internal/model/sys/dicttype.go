package model

type SysDictType struct {

	/** 字典主键 */
	DictId int `json:"dictId"`

	/** 字典名称 */
	DictName string `json:"dictName"`

	/** 字典类型 */
	DictType string `json:"dictType"`

	/** 状态（0正常 1停用） */
	Status string `json:"status"`
}
