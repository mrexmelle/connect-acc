package grading

import (
	"github.com/mrexmelle/connect-emp/internal/dto"
)

type PostRequestDto struct {
	Ehid      string `json:"ehid"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Grade     string `json:"grade"`
}

type PatchRequestDto struct {
	Fields map[string]interface{} `json:"fields"`
}

type GetResponseDto = dto.HttpResponseWithData[ViewEntity]
type PostResponseDto = dto.HttpResponseWithData[ViewEntity]
type PatchResponseDto = dto.HttpResponseWithoutData
type DeleteResponseDto = dto.HttpResponseWithoutData
