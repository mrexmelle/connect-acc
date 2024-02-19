package account

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-authx/pkg/libauthxc"
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Controller struct {
	ConfigService     *config.Service
	AccountService    *Service
	LocalErrorService *localerror.Service
	AuthxClient       *libauthxc.Client
}

func NewController(
	cfg *config.Service,
	as *Service,
	les *localerror.Service,
) *Controller {
	return &Controller{
		ConfigService:     cfg,
		AccountService:    as,
		LocalErrorService: les,
		AuthxClient: libauthxc.NewClient(
			cfg.ConfigRepository.GetAuthxHost(),
			cfg.ConfigRepository.GetAuthxPort(),
		),
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

	data, err := c.AccountService.RetrieveByEhidOrderByStartDateDesc(ehid)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New[[]career.Aggregate](
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

	data, err := c.AuthxClient.GetProfileByEhid(ehid)
	if err != nil || data.Error.Code != localerror.ErrSvcCodeNone {
		dtorespwithdata.NewError(
			localerror.ErrHttpClient.Error(),
			err.Error(),
		).RenderTo(w, http.StatusInternalServerError)
	}

	dtorespwithdata.New[libauthxc.ProfileEntity](
		data.Data,
		data.Error.Code,
		data.Error.Message,
	).RenderTo(w, http.StatusOK)
}
