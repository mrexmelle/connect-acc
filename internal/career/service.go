package career

import (
	"slices"
	"sort"

	"github.com/mrexmelle/connect-emp/internal/config"
	"github.com/mrexmelle/connect-emp/internal/dateinterval"
	"github.com/mrexmelle/connect-emp/internal/datesort"
	"github.com/mrexmelle/connect-emp/internal/datestr"
	"github.com/mrexmelle/connect-emp/internal/grading"
	"github.com/mrexmelle/connect-emp/internal/titling"
)

type Service struct {
	ConfigService  *config.Service
	GradingService *grading.Service
	TitlingService *titling.Service
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

	return s.mergeGradingsAndTitlings(gradings, titlings)
}

func (s *Service) mergeGradingsAndTitlings(
	gradings []grading.ViewEntity,
	titlings []titling.ViewEntity,
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
	}

	return aggs, nil
}
