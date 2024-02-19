package titling

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithoutdata"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Controller struct {
	ConfigService     *config.Service
	LocalErrorService *localerror.Service
	TitlingService    *Service
}

func NewController(cfg *config.Service, les *localerror.Service, svc *Service) *Controller {
	return &Controller{
		ConfigService:     cfg,
		LocalErrorService: les,
		TitlingService:    svc,
	}
}

// Get Titlings : HTTP endpoint to get titlings
// @Tags Titlings
// @Description Get a titling
// @Produce json
// @Param id path string true "Titling ID"
// @Success 200 {object} GetResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /titlings/{id} [GET]
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dtorespwithdata.NewError(
			localerror.ErrIdNotInteger.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}
	data, err := c.TitlingService.RetrieveById(id)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New(
		data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Post Titlings : HTTP endpoint to post new titlings
// @Tags Titlings
// @Description Post a new titlings
// @Accept json
// @Produce json
// @Param data body PostRequestDto true "Titling Request"
// @Success 200 {object} PostResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /titlings [POST]
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	var requestBody PostRequestDto
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtorespwithdata.NewError(
			localerror.ErrBadJson.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	data, err := c.TitlingService.Create(requestBody)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New(
		data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Patch Titlings : HTTP endpoint to patch a titling
// @Tags Titlings
// @Description Patch a titling
// @Accept json
// @Produce json
// @Param id path string true "Titling ID"
// @Param data body PatchRequestDto true "Titling Patch Request"
// @Success 200 {object} PatchResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /titlings/{id} [PATCH]
func (c *Controller) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrIdNotInteger.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	var requestBody PatchRequestDto
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrBadJson.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	err = c.TitlingService.UpdateById(requestBody.Fields, id)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Delete Titlings : HTTP endpoint to delete titlings
// @Tags Titlings
// @Description Delete a titling
// @Produce json
// @Param id path string true "Titling ID"
// @Success 200 {object} DeleteResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /titlings/{id} [DELETE]
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrIdNotInteger.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	err = c.TitlingService.DeleteById(id)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}
