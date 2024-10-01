package v1

import (
	"net/http"

	"github.com/KaySar12/NextZenOS/model"
	"github.com/KaySar12/NextZenOS/pkg/utils/common_err"
	"github.com/KaySar12/NextZenOS/service"
	"github.com/gin-gonic/gin"
)

func PostNotifyMessage(c *gin.Context) {
	name := c.Param("name")
	message := make(map[string]interface{})
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.INVALID_PARAMS, Message: err.Error()})
		return
	}

	service.MyService.Notify().SendNotify(name, message)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

func PostSystemStatusNotify(c *gin.Context) {
	message := make(map[string]interface{})
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, model.Result{Success: common_err.INVALID_PARAMS, Message: err.Error()})
		return
	}

	service.MyService.Notify().SettingSystemTempData(message)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
