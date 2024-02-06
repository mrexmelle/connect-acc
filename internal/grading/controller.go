package grading

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
	GradingService    *Service
	LocalErrorService *localerror.Service
}

func NewController(cfg *config.Service, svc *Service, les *localerror.Service) *Controller {
	return &Controller{
		ConfigService:     cfg,
		GradingService:    svc,
		LocalErrorService: les,
	}
}

// Get Gradings : HTTP endpoint to get gradings
// @Tags Gradings
// @Description Get a grading
// @Produce json
// @Param id path string true "Grading ID"
// @Success 200 {object} GetResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /gradings/{id} [GET]
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dtorespwithdata.NewError(
			localerror.ErrIdNotInteger.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}
	data, err := c.GradingService.RetrieveById(id)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New[ViewEntity](
		data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Post Gradings : HTTP endpoint to post new gradings
// @Tags Gradings
// @Description Post a new gradings
// @Accept json
// @Produce json
// @Param data body PostRequestDto true "Grading Request"
// @Success 200 {object} PostResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /gradings [POST]
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

	data, err := c.GradingService.Create(requestBody)
	info := c.LocalErrorService.Map(err)
	dtorespwithdata.New[ViewEntity](
		data,
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Patch Gradings : HTTP endpoint to patch a grading
// @Tags Gradings
// @Description Patch a grading
// @Accept json
// @Produce json
// @Param id path string true "Grading ID"
// @Param data body PatchRequestDto true "Grading Patch Request"
// @Success 200 {object} PatchResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /gradings/{id} [PATCH]
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

	err = c.GradingService.UpdateById(requestBody.Fields, id)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}

// Delete Gradings : HTTP endpoint to delete gradings
// @Tags Gradings
// @Description Delete a grading
// @Produce json
// @Param id path string true "Grading ID"
// @Success 200 {object} DeleteResponseDto "Success Response"
// @Failure 400 "BadRequest"
// @Failure 500 "InternalServerError"
// @Router /gradings/{id} [DELETE]
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dtorespwithoutdata.New(
			localerror.ErrIdNotInteger.Error(),
			err.Error(),
		).RenderTo(w, http.StatusBadRequest)
		return
	}

	err = c.GradingService.DeleteById(id)
	info := c.LocalErrorService.Map(err)
	dtorespwithoutdata.New(
		info.ServiceErrorCode,
		info.ServiceErrorMessage,
	).RenderTo(w, info.HttpStatusCode)
}
