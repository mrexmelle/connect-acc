package grading

import (
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithoutdata"
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

type GetResponseDto = dtorespwithdata.Class[ViewEntity]
type PostResponseDto = dtorespwithdata.Class[ViewEntity]
type PatchResponseDto = dtorespwithoutdata.Class
type DeleteResponseDto = dtorespwithoutdata.Class
