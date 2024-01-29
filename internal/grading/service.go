package grading

import (
	"database/sql"
	"time"

	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/localerror"
)

type Service struct {
	ConfigService     *config.Service
	GradingRepository Repository
}

func NewService(
	cfg *config.Service,
	r Repository,
) *Service {
	return &Service{
		ConfigService:     cfg,
		GradingRepository: r,
	}
}

func (s *Service) Create(req PostRequestDto) (*ViewEntity, error) {
	sd, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	var ed sql.NullTime
	if req.EndDate == "" {
		ed.Valid = false
	} else {
		ed.Time, err = time.Parse("2006-01-02", req.EndDate)
		ed.Valid = (err == nil)
	}

	result, err := s.GradingRepository.Create(&Entity{
		Ehid:      req.Ehid,
		StartDate: sd,
		EndDate:   ed,
		Grade:     req.Grade,
	})
	if err != nil {
		return nil, err
	}
	return toViewEntity(result), nil
}

func (s *Service) RetrieveById(id int) (*ViewEntity, error) {
	result, err := s.GradingRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return toViewEntity(result), nil
}

func (s *Service) UpdateById(fields map[string]interface{}, id int) error {
	return s.GradingRepository.UpdateById(fields, id)
}

func (s *Service) DeleteById(id int) error {
	err := s.GradingRepository.DeleteById(id)
	return err
}

func (s *Service) RetrieveByEhidOrderByStartDate(ehid string, orderDir string) ([]ViewEntity, error) {
	if orderDir != OrderAsc && orderDir != OrderDesc && orderDir != OrderNone {
		return []ViewEntity{}, localerror.ErrBadQueryParam
	}
	result, err := s.GradingRepository.FindByEhidOrderByStartDate(ehid, orderDir)
	if err != nil {
		return []ViewEntity{}, err
	}
	return toViewEntitySlice(result), nil
}

func (s *Service) RetrieveCurrentByNodeId(nodeId string) ([]ViewEntity, error) {
	result, err := s.GradingRepository.FindCurrentByNodeId(nodeId)
	if err != nil {
		return []ViewEntity{}, err
	}
	return toViewEntitySlice(result), nil
}

func (s *Service) RetrieveCurrentByEhid(ehid string) ([]ViewEntity, error) {
	result, err := s.GradingRepository.FindCurrentByEhid(ehid)
	if err != nil {
		return []ViewEntity{}, err
	}
	return toViewEntitySlice(result), nil
}
