package account

import (
	"github.com/mrexmelle/connect-authx/pkg/libauthxc"
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
)

type GetProfileResponseDto = dtorespwithdata.Class[libauthxc.ProfileEntity]
type GetCareerResponseDto = dtorespwithdata.Class[[]career.Aggregate]
