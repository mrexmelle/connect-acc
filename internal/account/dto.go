package account

import (
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
	"github.com/mrexmelle/connect-emp/internal/profile"
)

type GetProfileResponseDto = dtorespwithdata.Class[profile.Aggregate]
type GetCareerResponseDto = dtorespwithdata.Class[[]career.Aggregate]
