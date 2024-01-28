package titling

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
	TitlingService *Service
}

func NewController(cfg *config.Service, svc *Service) *Controller {
	return &Controller{
		ConfigService:  cfg,
		TitlingService: svc,
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
		dtobuilderwithdata.New[Entity](nil, localerror.ErrIdNotInteger).RenderTo(w)
		return
	}
	data, err := c.TitlingService.RetrieveById(id)
	dtobuilderwithdata.New[ViewEntity](data, err).RenderTo(w)
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
		dtobuilderwithdata.New[Entity](nil, localerror.ErrBadJson).RenderTo(w)
		return
	}

	data, err := c.TitlingService.Create(requestBody)
	dtobuilderwithdata.New[ViewEntity](data, err).RenderTo(w)
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
		dtobuilderwithoutdata.New(localerror.ErrIdNotInteger).RenderTo(w)
		return
	}

	var requestBody PatchRequestDto
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		dtobuilderwithoutdata.New(localerror.ErrBadJson).RenderTo(w)
		return
	}

	err = c.TitlingService.UpdateById(requestBody.Fields, id)
	dtobuilderwithoutdata.New(err).RenderTo(w)
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
		dtobuilderwithoutdata.New(localerror.ErrIdNotInteger).RenderTo(w)
		return
	}

	err = c.TitlingService.DeleteById(id)
	dtobuilderwithoutdata.New(err).RenderTo(w)
}
