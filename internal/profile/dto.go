package profile

import (
	"github.com/mrexmelle/connect-authx/pkg/libauthxc"
	"github.com/mrexmelle/connect-emp/internal/dto/dtorespwithdata"
)

type GetResponseDto = dtorespwithdata.Class[libauthxc.ProfileEntity]
