package my

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
)

type AddressBookCollection struct {
}

// Create Create address book collection
// @Tags MyAddressBookCollection
// @Summary Create address book collection
// @Description Create address book collection
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book collection information"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection/create [post]
// @Security token
func (abc *AddressBookCollection) Create(c *gin.Context) {
	f := &model.AddressBookCollection{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.UserService.CurUser(c)
	f.UserId = u.Id
	err := service.AllService.AddressBookService.CreateCollection(f)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List List
// @Tags MyAddressBookCollection
// @Summary Address book collection list
// @Description Address book collection list
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=model.AddressBookCollectionList}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection/list [get]
// @Security token
func (abc *AddressBookCollection) List(c *gin.Context) {
	query := &admin.AddressBookCollectionQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	u := service.AllService.UserService.CurUser(c)
	query.UserId = int(u.Id)
	res := service.AllService.AddressBookService.ListCollection(query.Page, query.PageSize, func(tx *gorm.DB) {
		tx.Where("user_id = ?", query.UserId)
	})
	response.Success(c, res)
}

// Update Edit
// @Tags MyAddressBookCollection
// @Summary Edit address book collection
// @Description Edit address book collection
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book collection information"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection/update [post]
// @Security token
func (abc *AddressBookCollection) Update(c *gin.Context) {
	f := &model.AddressBookCollection{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	if f.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	//if f.UserId != u.Id {
	//	response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
	//	return
	//}
	ex := service.AllService.AddressBookService.CollectionInfoById(f.Id)
	if ex.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	if ex.UserId != u.Id {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}

	err := service.AllService.AddressBookService.UpdateCollection(f)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete Delete
// @Tags MyAddressBookCollection
// @Summary Delete address book collection
// @Description Delete address book collection
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book collection information"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection/delete [post]
// @Security token
func (abc *AddressBookCollection) Delete(c *gin.Context) {
	f := &model.AddressBookCollection{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.Id
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	ex := service.AllService.AddressBookService.CollectionInfoById(f.Id)
	if ex.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	if ex.UserId != u.Id {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}
	err := service.AllService.AddressBookService.DeleteCollection(ex)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
}
