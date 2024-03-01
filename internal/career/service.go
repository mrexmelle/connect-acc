package career

import (
	"slices"
	"sort"

	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dateinterval"
	"github.com/mrexmelle/connect-emp/internal/datesort"
	"github.com/mrexmelle/connect-emp/internal/datestr"
	"github.com/mrexmelle/connect-emp/internal/grading"
	"github.com/mrexmelle/connect-emp/internal/localerror"
	"github.com/mrexmelle/connect-emp/internal/titling"
	"github.com/mrexmelle/connect-org/pkg/liborgc"
)

type Service struct {
	ConfigService  *config.Service
	GradingService *grading.Service
	TitlingService *titling.Service
	OrgClient      *liborgc.Client
}

func NewService(
	cfg *config.Service,
	gs *grading.Service,
	ts *titling.Service,
) *Service {
	return &Service{
		ConfigService:  cfg,
		GradingService: gs,
		TitlingService: ts,
		OrgClient: liborgc.NewClient(
			cfg.ConfigRepository.GetOrgHost(),
			cfg.ConfigRepository.GetOrgPort(),
		),
	}
}

func (s *Service) RetrieveCurrentByEhid(ehid string) (*Aggregate, error) {
	g, err := s.GradingService.RetrieveCurrentByEhid(ehid)
	if err != nil {
		return nil, err
	}

	t, err := s.TitlingService.RetrieveCurrentByEhid(ehid)
	if err != nil {
		return nil, err
	}

	m, err := s.OrgClient.GetMemberNodesByEhid(ehid)
	if err != nil || len(*m.Data) == 0 {
		return nil, err
	}

	gsd, err := datestr.NewFromString(g.StartDate)
	if err != nil {
		return nil, err
	}

	ged, err := datestr.NewFromString(g.EndDate)
	if err != nil {
		return nil, err
	}

	tsd, err := datestr.NewFromString(t.StartDate)
	if err != nil {
		return nil, err
	}

	ted, err := datestr.NewFromString(t.EndDate)
	if err != nil {
		return nil, err
	}

	msd, err := datestr.NewFromString((*m.Data)[0].StartDate)
	if err != nil {
		return nil, err
	}

	med, err := datestr.NewFromString((*m.Data)[0].EndDate)
	if err != nil {
		return nil, err
	}

	sd := s.maxDate([]datestr.Class{*gsd, *tsd, *msd})
	ed := s.minDate([]datestr.Class{*ged, *ted, *med})

	return &Aggregate{
		StartDate:        (*sd).AsString(),
		EndDate:          (*ed).AsString(),
		Grade:            g.Grade,
		Title:            t.Title,
		OrganizationNode: (*m.Data)[0].NodeId,
	}, nil
}

func (s *Service) RetrieveByEhidOrderByStartDateDesc(ehid string) ([]Aggregate, error) {
	gradings, err := s.GradingService.RetrieveByEhidOrderByStartDate(ehid, grading.OrderDesc)
	if err != nil {
		return []Aggregate{}, err
	}
	titlings, err := s.TitlingService.RetrieveByEhidOrderByStartDate(ehid, titling.OrderDesc)
	if err != nil {
		return []Aggregate{}, err
	}
	membership, err := s.OrgClient.GetMemberHistoryByEhidOrderByStartDateDesc(ehid)
	if err != nil || membership.Error.Code != localerror.ErrSvcCodeNone {
		return nil, err
	}

	return s.mergeHistories(gradings, titlings, *membership.Data)
}

func (s *Service) mergeHistories(
	gradings []grading.ViewEntity,
	titlings []titling.ViewEntity,
	memberships []liborgc.MembershipViewEntity,
) ([]Aggregate, error) {
	aggs := []Aggregate{}

	endDates := []string{}
	earliestStartDate := ""
	for _, g := range gradings {
		idx := slices.Index(endDates, g.EndDate)
		if idx == -1 {
			endDates = append(endDates, g.EndDate)
		}

		if earliestStartDate == "" || g.StartDate < earliestStartDate {
			earliestStartDate = g.StartDate
		}
	}

	for _, t := range titlings {
		idx := slices.Index(endDates, t.EndDate)
		if idx == -1 {
			endDates = append(endDates, t.EndDate)
		}

		if earliestStartDate == "" || t.StartDate < earliestStartDate {
			earliestStartDate = t.StartDate
		}
	}

	for _, m := range memberships {
		idx := slices.Index(endDates, m.EndDate)

		if idx == -1 {
			endDates = append(endDates, m.EndDate)
		}

		if earliestStartDate == "" || m.StartDate < earliestStartDate {
			earliestStartDate = m.StartDate
		}
	}

	sort.Sort(sort.Reverse(datesort.DateStringSlice(endDates)))
	for i := 0; i < len(endDates); i++ {
		var startDate = ""
		if i == len(endDates)-1 {
			startDate = earliestStartDate
		} else {
			ds, err := datestr.NewFromString(endDates[i+1])
			if err != nil {
				return []Aggregate{}, err
			}
			startDate = ds.OffsetAndClone(+1).AsString()
		}
		aggs = append(aggs, Aggregate{
			StartDate: startDate,
			EndDate:   endDates[i],
		})
	}

	for i, a := range aggs {
		aggsInterval, err := dateinterval.NewFromStrings(a.StartDate, a.EndDate)
		if err != nil {
			return []Aggregate{}, err
		}
		for _, g := range gradings {
			gradeInterval, err := dateinterval.NewFromStrings(g.StartDate, g.EndDate)
			if err != nil {
				return []Aggregate{}, err
			}
			if gradeInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].Grade = g.Grade
				break
			}
		}
		for _, t := range titlings {
			titleInterval, err := dateinterval.NewFromStrings(t.StartDate, t.EndDate)
			if err != nil {
				return []Aggregate{}, err
			}
			if titleInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].Title = t.Title
				break
			}
		}

		for _, m := range memberships {
			membershipInterval, err := dateinterval.NewFromStrings(m.StartDate, m.EndDate)
			if err != nil {
				return []Aggregate{}, err
			}
			if membershipInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].OrganizationNode = m.NodeId
				break
			}
		}
	}

	return aggs, nil
}

func (s *Service) maxDate(dates []datestr.Class) *datestr.Class {
	if len(dates) == 0 {
		return nil
	}

	maxPtr := &dates[0]
	for i := 1; i < len(dates); i++ {
		if dates[i].IsAfterOrEquals(maxPtr) {
			maxPtr = &dates[i]
		}
	}

	return maxPtr
}

func (s *Service) minDate(dates []datestr.Class) *datestr.Class {
	if len(dates) == 0 {
		return nil
	}

	minPtr := &dates[0]
	for i := 1; i < len(dates); i++ {
		if dates[i].IsBeforeOrEquals(minPtr) {
			minPtr = &dates[i]
		}
	}

	return minPtr
}
