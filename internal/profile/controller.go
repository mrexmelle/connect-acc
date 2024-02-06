package profile

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-authx/pkg/libauthxc"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Controller struct {
	ConfigService     *config.Service
	LocalErrorService *localerror.Service
	AuthxClient       *libauthxc.Client
}

func NewController(cfg *config.Service, les *localerror.Service) *Controller {
	return &Controller{
		ConfigService:     cfg,
		LocalErrorService: les,
		AuthxClient: libauthxc.NewClient(
			cfg.ConfigRepository.GetAuthxHost(),
			cfg.ConfigRepository.GetAuthxPort(),
		),
	}
}

// Get Profiles : HTTP endpoint to get profiles
// @Tags Profiles
// @Description Get a profile
// @Produce json
// @Param ehid path string true "EHID"
// @Success 200 {object} GetResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /profiles/{ehid} [GET]
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
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
