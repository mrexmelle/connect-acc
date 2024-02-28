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

	var sd = ""
	if gsd.IsAfter(tsd) {
		sd = gsd.AsString()
	} else {
		sd = tsd.AsString()
	}
	var ed = ""
	if ged.IsBefore(ted) {
		ed = ged.AsString()
	} else {
		ed = ted.AsString()
	}

	return &Aggregate{
		StartDate: sd,
		EndDate:   ed,
		Grade:     g.Grade,
		Title:     t.Title,
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
	membership, err := s.OrgClient.GetMembershipHistoryByEhidOrderByStartDateDesc(ehid)
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
	for i, g := range gradings {
		idx := slices.Index(endDates, g.EndDate)
		if idx == -1 {
			endDates = append(endDates, g.EndDate)
		}

		if i == 0 {
			earliestStartDate = g.StartDate
		} else {
			if g.StartDate < earliestStartDate {
				earliestStartDate = g.StartDate
			}
		}
	}
	for _, t := range titlings {
		idx := slices.Index(endDates, t.EndDate)
		if idx == -1 {
			endDates = append(endDates, t.EndDate)
		}

		if t.StartDate < earliestStartDate {
			earliestStartDate = t.StartDate
		}
	}

	for _, m := range memberships {
		idx := slices.Index(endDates, m.EndDate)

		if idx == -1 {
			endDates = append(endDates, m.EndDate)
		}

		if m.StartDate < earliestStartDate {
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
			titleInterval, err := dateinterval.NewFromStrings(m.StartDate, m.EndDate)
			if err != nil {
				return []Aggregate{}, err
			}
			if titleInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].OrganizationNode = m.NodeId
				break
			}
		}
	}

	return aggs, nil
}
