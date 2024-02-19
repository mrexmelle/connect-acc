package account

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Controller struct {
	ConfigService     *config.Service
	AccountService    *Service
	CareerService     *career.Service
	LocalErrorService *localerror.Service
}

func NewController(
	cfg *config.Service,
	as *Service,
	cs *career.Service,
	les *localerror.Service,
) *Controller {
	return &Controller{
		ConfigService:     cfg,
		AccountService:    as,
		CareerService:     cs,
		LocalErrorService: les,
	}
}

// Get Career : HTTP endpoint to get the career of an account
// @Tags Accounts
// @Description Get a career
// @Produce json
// @Param ehid path string true "EHID"
// @Success 200 {object} GetCareerResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /accounts/{ehid}/career [GET]
func (c *Controller) GetCareer(w http.ResponseWriter, r *http.Request) {
	ehid := chi.URLParam(r, "ehid")

	data, err := c.CareerService.RetrieveByEhidOrderByStartDateDesc(ehid)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New(
		&data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, http.StatusOK)
}

// Get Profile : HTTP endpoint to get the profile of an account
// @Tags Accounts
// @Description Get a profile
// @Produce json
// @Param ehid path string true "EHID"
// @Success 200 {object} GetProfileResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /accounts/{ehid}/profile [GET]
func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) {
	ehid := chi.URLParam(r, "ehid")
	data, err := c.AccountService.RetrieveProfile(ehid)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New(
		data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}
