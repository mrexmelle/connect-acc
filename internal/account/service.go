package account

import (
	"github.com/mrexmelle/connect-authx/pkg/libauthxc"
	"github.com/mrexmelle/connect-emp/internal/career"
	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/localerror"
	"github.com/mrexmelle/connect-emp/internal/profile"
)

type Service struct {
	ConfigService *config.Service
	CareerService *career.Service
	AuthxClient   *libauthxc.Client
}

func NewService(
	cfg *config.Service,
	cs *career.Service,
) *Service {
	return &Service{
		ConfigService: cfg,
		CareerService: cs,
		AuthxClient: libauthxc.NewClient(
			cfg.ConfigRepository.GetAuthxHost(),
			cfg.ConfigRepository.GetAuthxPort(),
		),
	}
}

func (s *Service) RetrieveCareer(ehid string) ([]career.Aggregate, error) {
	return s.CareerService.RetrieveByEhidOrderByStartDateDesc(ehid)
}

func (s *Service) RetrieveProfile(ehid string) (*profile.Aggregate, error) {
	p, err := s.AuthxClient.GetProfileByEhid(ehid)
	if err != nil || p.Error.Code != localerror.ErrSvcCodeNone {
		return nil, err
	}

	agg := &profile.Aggregate{
		Ehid:         ehid,
		EmployeeId:   p.Data.EmployeeId,
		Name:         p.Data.Name,
		EmailAddress: p.Data.EmailAddress,
		Dob:          p.Data.Dob,
	}

	career, err := s.CareerService.RetrieveCurrentByEhid(ehid)
	if err != nil {
		return nil, err
	}
	agg.Grade = career.Grade
	agg.Title = career.Title

	return agg, nil
}
