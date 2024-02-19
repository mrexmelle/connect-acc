package account

import (
	"fmt"
	"slices"
	"sort"

	"github.com/mrexmelle/connect-emp/internal/career"
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

func (s *Service) RetrieveByEhidOrderByStartDateDesc(ehid string) ([]career.Aggregate, error) {
	aggs := []career.Aggregate{}

	gradings, err := s.GradingService.RetrieveByEhidOrderByStartDate(ehid, grading.OrderDesc)
	if err != nil {
		fmt.Printf("error di grading.RetrieveByEhidOrderByStartDate")
		return []career.Aggregate{}, err
	}
	titlings, err := s.TitlingService.RetrieveByEhidOrderByStartDate(ehid, titling.OrderDesc)
	if err != nil {
		fmt.Printf("error di titling.RetrieveByEhidOrderByStartDate")
		return []career.Aggregate{}, err
	}

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
				return []career.Aggregate{}, err
			}
			startDate = ds.OffsetAndClone(+1).AsString()
		}
		aggs = append(aggs, career.Aggregate{
			Ehid:      ehid,
			StartDate: startDate,
			EndDate:   endDates[i],
		})
	}

	for i, a := range aggs {
		aggsInterval, err := dateinterval.NewFromStrings(a.StartDate, a.EndDate)
		if err != nil {
			return []career.Aggregate{}, err
		}
		for _, g := range gradings {
			gradeInterval, err := dateinterval.NewFromStrings(g.StartDate, g.EndDate)
			if err != nil {
				return []career.Aggregate{}, err
			}
			if gradeInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].Grade = g.Grade
				break
			}
		}
		for _, t := range titlings {
			titleInterval, err := dateinterval.NewFromStrings(t.StartDate, t.EndDate)
			if err != nil {
				return []career.Aggregate{}, err
			}
			if titleInterval.IsEncompassingInterval(aggsInterval) {
				aggs[i].Title = t.Title
				break
			}
		}
	}

	return aggs, nil
}
