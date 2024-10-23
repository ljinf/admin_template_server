package model

type SysDictData struct {

	/** 字典编码 */
	DictCode int `json:"dictCode"`

	/** 字典排序 */
	DictSort int `json:"dictSort"`

	/** 字典标签 */
	DictLabel string `json:"dictLabel"`

	/** 字典键值 */
	DictValue string `json:"dictValue"`

	/** 字典类型 */
	DictType string `json:"dictType"`

	/** 样式属性（其他样式扩展） */
	CssClass string `json:"cssClass"`

	/** 表格字典样式 */
	ListClass string `json:"listClass"`

	/** 是否默认（Y是 N否） */
	//@Excel(name = "是否默认", readConverterExp = "Y=是,N=否")
	IsDefault string `json:"isDefault"`

	/** 状态（0正常 1停用） */
	//@Excel(name = "状态", readConverterExp = "0=正常,1=停用")
	Status string `json:"status"`
}
