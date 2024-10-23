package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	v1 "poem_server_admin/api/v1"
	"poem_server_admin/internal/handler"
	model "poem_server_admin/internal/model/sys"
	service "poem_server_admin/internal/service/sys"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type MenuHandler interface {
	GetRouters(ctx *gin.Context)
	MenuList(ctx *gin.Context)
	GetInfo(ctx *gin.Context)
	TreeSelect(ctx *gin.Context)
	RoleMenuTreeSelect(ctx *gin.Context)
	MenuTreeSelect(ctx *gin.Context)
	AddMenu(ctx *gin.Context)
	UpdateMenu(ctx *gin.Context)
	DelMenu(ctx *gin.Context)
}

type menuHandler struct {
	*handler.Handler
	menuService service.MenuService
}

func NewMenuHandler(handler *handler.Handler, service service.MenuService) MenuHandler {
	return &menuHandler{
		Handler:     handler,
		menuService: service,
	}
}

// 所有菜单
func (h *menuHandler) MenuTreeSelect(ctx *gin.Context) {
	sysMenus, err := h.menuService.MenuList(ctx, model.SysMenu{Status: "0"})
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, h.buildMenuTreeSelect(sysMenus))
}

// 加载对应角色菜单列表树
func (h *menuHandler) RoleMenuTreeSelect(ctx *gin.Context) {
	id := ctx.Param("roleId")
	roleId, _ := strconv.Atoi(id)
	if roleId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	userId, ok := ctx.Get("userId")
	if !ok {
		v1.HandleError(ctx, http.StatusOK, v1.ErrUnauthorized, nil)
		return
	}

	sysMenus, err := h.menuService.GetMenuInfoByUserId(ctx, userId.(int64))
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	menuIdsByRoleId, err := h.menuService.GetMenuIdsByRoleId(ctx, int64(roleId))
	if err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	res := map[string]interface{}{
		"checkedKeys": menuIdsByRoleId,
		"menus":       h.buildMenuTreeSelect(sysMenus),
	}
	v1.HandleSuccess(ctx, res)
}

// 获取菜单下拉树列表
func (h *menuHandler) TreeSelect(ctx *gin.Context) {

}

// 根据菜单编号获取详细信息
func (h *menuHandler) GetInfo(ctx *gin.Context) {
	id := ctx.Param("menuId")

	menuId, _ := strconv.Atoi(id)
	if menuId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	menuInfo, err := h.menuService.GetMenuInfoById(ctx, int64(menuId))
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, menuInfo)

}

// 获取菜单列表
func (h *menuHandler) MenuList(ctx *gin.Context) {
	var param v1.MenuRequest
	if err := ctx.ShouldBind(&param); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	req := model.SysMenu{
		MenuName: param.MenuName,
		Visible:  param.Visible,
		Status:   "0",
	}

	list, err := h.menuService.MenuList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrInternalServerError, nil)
		return
	}

	res := map[string]interface{}{
		"rows":  list,
		"total": len(list),
	}
	v1.HandleSuccess(ctx, res)
}

// 获取路由信息
func (h *menuHandler) GetRouters(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		v1.HandleError(ctx, http.StatusOK, v1.ErrLoginExpired, nil)
		return
	}

	menus, err := h.menuService.GetRouters(ctx, userId.(int64))
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	routerList := h.buildMenus(menus)
	v1.HandleSuccess(ctx, routerList)
}

// 菜单增改删
func (h *menuHandler) AddMenu(ctx *gin.Context) {

	//todo 权限校验

	var parm v1.MenuInfoRequest
	if err := ctx.ShouldBind(&parm); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	req := model.SysMenu{
		MenuName:  parm.MenuName,
		ParentId:  parm.ParentId,
		OrderNum:  parm.OrderNum,
		Path:      parm.Path,
		Component: parm.Component,
		Query:     parm.Query,
		IsFrame:   parm.IsFrame,
		MenuType:  parm.MenuType,
		Visible:   parm.Visible,
		Status:    parm.Status,
		Perms:     parm.Perms,
		Icon:      parm.Icon,
	}

	if err := h.menuService.AddMenu(ctx, &req); err != nil {
		h.Logger.Error(err.Error(), zap.Any("param", parm))
		v1.HandleError(ctx, http.StatusOK, v1.ErrCreateMenuFailed, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *menuHandler) UpdateMenu(ctx *gin.Context) {

	//todo 权限校验

	var parm v1.MenuInfoRequest
	if err := ctx.ShouldBind(&parm); err != nil {
		h.Logger.Error(err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if parm.MenuId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	req := model.SysMenu{
		MenuId:    parm.MenuId,
		MenuName:  parm.MenuName,
		ParentId:  parm.ParentId,
		OrderNum:  parm.OrderNum,
		Path:      parm.Path,
		Component: parm.Component,
		Query:     parm.Query,
		IsFrame:   parm.IsFrame,
		MenuType:  parm.MenuType,
		Visible:   parm.Visible,
		Status:    parm.Status,
		Perms:     parm.Perms,
		Icon:      parm.Icon,
	}

	if err := h.menuService.UpdateMenu(ctx, &req); err != nil {
		h.Logger.Error(err.Error(), zap.Any("param", parm))
		v1.HandleError(ctx, http.StatusOK, v1.ErrEditMenuFailed, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *menuHandler) DelMenu(ctx *gin.Context) {

	//todo 权限校验

	id := ctx.Param("menuId")

	menuId, _ := strconv.Atoi(id)
	if menuId == 0 {
		v1.HandleError(ctx, http.StatusOK, v1.ErrParamError, nil)
		return
	}

	if err := h.menuService.DelMenu(ctx, int64(menuId)); err != nil {
		v1.HandleError(ctx, http.StatusOK, v1.ErrDelMenuFailed, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)

}

// 构建前端路由所需要的菜单
func (h *menuHandler) buildMenus(list []*model.SysMenu) []model.RouterVo {

	var routerList []model.RouterVo
	for _, v := range list {
		router := model.RouterVo{
			Name:      h.getRouteName(*v),
			Path:      h.getRouterPath(*v),
			Component: h.getComponent(*v),
			Query:     v.Query,
			Hidden:    isHidden(v.Visible),
			Meta: model.MetaVo{
				Title:   v.MenuName,
				Icon:    v.Icon,
				NoCache: false,
				Link:    isHttp(v.Path),
			},
		}
		cMenus := v.Children
		if len(cMenus) > 0 && v.MenuType == "M" {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = h.buildMenus(cMenus)
		} else if h.isMenuFrame(*v) {
			router.Meta = model.MetaVo{}
			var childrenList []model.RouterVo
			children := model.RouterVo{
				Path:      v.Path,
				Component: v.Component,
				Name:      v.Path,
				Meta: model.MetaVo{
					Title:   v.MenuName,
					Icon:    v.Icon,
					NoCache: false,
					Link:    isHttp(v.Path),
				},
				Query: v.Query,
			}
			childrenList = append(childrenList, children)
			router.Children = childrenList
		} else if v.ParentId == 0 && h.isInnerLink(*v) {
			router.Meta = model.MetaVo{
				Title: v.MenuName,
				Icon:  v.Icon,
			}
			router.Path = "/"

			var childrenList []model.RouterVo
			children := model.RouterVo{
				Path:      h.innerLinkReplaceEach(v.Path),
				Component: "InnerLink",
				Name:      v.Path,
				Meta: model.MetaVo{
					Title: v.MenuName,
					Icon:  v.Icon,
					Link:  isHttp(v.Path),
				},
				Query: v.Query,
			}
			childrenList = append(childrenList, children)
			router.Children = childrenList
		}
		routerList = append(routerList, router)
	}
	return routerList
}

// 构建前端所需要下拉树结构
func (h *menuHandler) buildMenuTreeSelect(list []*model.SysMenu) []model.TreeSelect {
	menuTree := h.buildMenuTree(list)
	var treeList []model.TreeSelect
	for _, v := range menuTree {
		treeItem := model.TreeSelect{
			Id:       int(v.MenuId),
			Label:    v.MenuName,
			Children: h.menuToTree(v.Children),
		}
		treeList = append(treeList, treeItem)
	}
	return treeList
}

func (h *menuHandler) menuToTree(list []*model.SysMenu) []model.TreeSelect {
	var treeList []model.TreeSelect
	for _, v := range list {
		treeItem := model.TreeSelect{
			Id:       int(v.MenuId),
			Label:    v.MenuName,
			Children: h.menuToTree(v.Children),
		}
		treeList = append(treeList, treeItem)
	}
	return treeList
}

// 构建前端所需要树结构
func (h *menuHandler) buildMenuTree(list []*model.SysMenu) []*model.SysMenu {

	var returnList []*model.SysMenu
	var idsList []int64
	for _, v := range list {
		idsList = append(idsList, v.MenuId)
	}

	for _, v := range list {
		// 如果是顶级节点, 遍历该父节点的所有子节点
		if !exist(idsList, v.ParentId) {
			h.recursionFn(list, v)
			returnList = append(returnList, v)
		}
	}

	if len(returnList) < 1 {
		returnList = list
	}
	return returnList

}

// 获取路由名称
func (h *menuHandler) getRouteName(menu model.SysMenu) string {
	routerName := menu.Path
	// 非外链并且是一级目录（类型为目录）
	if h.isMenuFrame(menu) {
		routerName = ""
	}
	return capitalizeFirstLetter(routerName)
}

// 是否为菜单内部跳转
func (h *menuHandler) isMenuFrame(menu model.SysMenu) bool {
	return menu.ParentId == 0 && menu.MenuType == "C" && menu.IsFrame == "1"
}

// 获取路由地址
func (h *menuHandler) getRouterPath(menu model.SysMenu) string {
	routerPath := menu.Path
	// 内链打开外网方式
	if menu.ParentId != 0 && h.isInnerLink(menu) {
		routerPath = h.innerLinkReplaceEach(routerPath)
	}

	// 非外链并且是一级目录（类型为目录）
	if menu.ParentId == 0 && menu.MenuType == "M" && menu.IsFrame == "1" {
		routerPath = "/" + menu.Path
	} else if h.isMenuFrame(menu) {
		// 非外链并且是一级目录（类型为菜单）
		routerPath = "/"
	}
	return routerPath
}

// 是否为内链组件
func (h *menuHandler) isInnerLink(menu model.SysMenu) bool {
	return menu.IsFrame == "1" && strings.HasPrefix(menu.Path, "http")
}

// 内链域名特殊字符替换
func (h *menuHandler) innerLinkReplaceEach(path string) string {
	//return StringUtils.replaceEach(path, new String[] { Constants.HTTP, Constants.HTTPS, Constants.WWW, ".", ":" },
	////                new String[] { "", "", "", "/", "/" });
	temp := strings.SplitN(path, ":", 2)[1]
	return strings.ReplaceAll(temp, ".", "/")
}

func (h *menuHandler) getComponent(menu model.SysMenu) string {
	component := "Layout"
	if menu.Component != "" && !h.isMenuFrame(menu) {
		component = menu.Component
	} else if menu.Component == "" && menu.ParentId != 0 && h.isInnerLink(menu) {
		component = "InnerLink"
	} else if menu.Component == "" && h.isParentView(menu) {
		component = "ParentView"
	}
	return component
}

// 是否为parent_view组件
func (h *menuHandler) isParentView(menu model.SysMenu) bool {
	return menu.ParentId != 0 && menu.MenuType == "M"
}

func isHttp(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return ""
}

func isHidden(visible string) bool {
	//菜单状态（0显示 1隐藏）
	if visible == "0" {
		return false
	}
	return true
}

// capitalizeFirstLetter takes a string and returns a new string with the first letter capitalized.
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func exist(list []int64, target int64) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// 递归列表
func (h *menuHandler) recursionFn(list []*model.SysMenu, t *model.SysMenu) {

	// 得到子节点列表
	childList := h.getChildList(list, t)
	t.Children = childList
	for _, tChild := range childList {
		if h.hasChild(list, tChild) {
			h.recursionFn(list, tChild)
		}
	}
}

// 得到子节点列表
func (h *menuHandler) getChildList(list []*model.SysMenu, t *model.SysMenu) []*model.SysMenu {

	tlist := []*model.SysMenu{}
	for _, v := range list {
		if v.ParentId == t.MenuId {
			tlist = append(tlist, v)
		}
	}
	return tlist

}

// 判断是否有子节点
func (h *menuHandler) hasChild(list []*model.SysMenu, t *model.SysMenu) bool {
	return len(h.getChildList(list, t)) > 0
}
