package claims

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/ctxfilter"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
	"github.com/gin-gonic/gin"
)

func createClaim(
	dB db.DB,
	claimService services.ClaimService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.ClaimSubmissionForm
		if err := c.BindJSON(&req); err != nil {
			appErr := apperr.NewErrorWithType(err, apperr.BadRequest)
			utils.HandleError(c, appErr)
			return
		}

		// ✅ SubmitClaim runs the full validation pipeline:
		// eligibility → benefit limit → fraud check → approval decision
		result, err := claimService.SubmitClaim(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func getClaim(
	dB db.DB,
	claimService services.ClaimService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		claimID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(err, apperr.BadRequest)
			utils.HandleError(c, appErr)
			return
		}

		claim, err := claimService.GetClaimByID(c.Request.Context(), dB, claimID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, claim)
	}
}

func listClaims(
	dB db.DB,
	claimService services.ClaimService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		memberID := c.Param("memberID")
		if len(strings.TrimSpace(memberID)) < 1 {
			appErr := apperr.NewBadRequest("member id not set")
			utils.HandleError(c, appErr)
			return
		}

		filter, err := ctxfilter.FilterFromContext(c)
		if err != nil {
			appErr := apperr.NewErrorWithType(err, apperr.BadRequest)
			utils.HandleError(c, appErr)
			return
		}

		claimList, err := claimService.GetClaims(c.Request.Context(), dB, memberID, filter)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, claimList)
	}
}