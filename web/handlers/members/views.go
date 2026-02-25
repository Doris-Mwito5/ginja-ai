package members

import (
	"net/http"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
	"github.com/gin-gonic/gin"
)

func createMember(
	dB db.DB,
	memberService services.MemberService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.Member
		err := c.BindJSON(&req)
		if err != nil {
			utils.HandleError(c, apperr.NewErrorWithType(err, apperr.BadRequest))
			return
		}

		member, err := memberService.CreateMember(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
				return
			}
	
			c.JSON(http.StatusCreated, member)
	}
}
