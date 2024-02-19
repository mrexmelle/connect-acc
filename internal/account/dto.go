package account

import (
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
)

type GetProfileResponseDto = dtorespwithdata.Class[career.Aggregate]
type GetCareerResponseDto = dtorespwithdata.Class[[]career.Aggregate]
