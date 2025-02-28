package user

import (
	"github.com/gin-gonic/gin"
	coredto "github.com/unusualcodeorg/goserve/arch/dto"
	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/common"
)

type controller struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/profile", authProvider, authorizeProvider),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getPublicProfileHandler)
	private := group.Use(c.Authentication())
	private.GET("/mine", c.getPrivateProfileHandler)
}

func (c *controller) getPublicProfileHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	data, err := c.service.GetUserPublicProfile(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}

func (c *controller) getPrivateProfileHandler(ctx *gin.Context) {
	user := c.MustGetUser(ctx)

	data, err := c.service.GetUserPrivateProfile(user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}