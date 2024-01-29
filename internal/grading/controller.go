package grading

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dto/dtobuilderwithdata"
	"github.com/mrexmelle/connect-emp/internal/dto/dtobuilderwithoutdata"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Controller struct {
	ConfigService  *config.Service
	GradingService *Service
}

func NewController(cfg *config.Service, svc *Service) *Controller {
	return &Controller{
		ConfigService:  cfg,
		GradingService: svc,
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
		dtobuilderwithdata.New[Entity](nil, localerror.ErrIdNotInteger).RenderTo(w)
		return
	}
	data, err := c.GradingService.RetrieveById(id)
	dtobuilderwithdata.New[ViewEntity](data, err).RenderTo(w)
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
		dtobuilderwithdata.New[Entity](nil, localerror.ErrBadJson).RenderTo(w)
		return
	}

	data, err := c.GradingService.Create(requestBody)
	dtobuilderwithdata.New[ViewEntity](data, err).RenderTo(w)
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
		dtobuilderwithoutdata.New(localerror.ErrIdNotInteger).RenderTo(w)
		return
	}

	var requestBody PatchRequestDto
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtobuilderwithoutdata.New(localerror.ErrBadJson).RenderTo(w)
		return
	}

	err = c.GradingService.UpdateById(requestBody.Fields, id)
	dtobuilderwithoutdata.New(err).RenderTo(w)
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
		dtobuilderwithoutdata.New(localerror.ErrIdNotInteger).RenderTo(w)
		return
	}

	err = c.GradingService.DeleteById(id)
	dtobuilderwithoutdata.New(err).RenderTo(w)
}
