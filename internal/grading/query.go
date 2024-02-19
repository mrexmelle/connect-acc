package grading

import (
	"gorm.io/gorm"
)

var (
	FieldsAll = []string{
		"id",
		"ehid",
		"start_date",
		"end_date",
		"grade",
	}

	FieldsAllExceptId = []string{
		"ehid",
		"start_date",
		"end_date",
		"grade",
	}

	FieldsAllExceptIdAndEndDate = []string{
		"ehid",
		"start_date",
		"grade",
	}

	FieldsPatchable = []string{
		"grade",
		"end_date",
	}

	OrderAsc  = "ASC"
	OrderDesc = "DESC"
	OrderNone = ""
)

type Query interface {
	SelectById(fields []string, id int) *gorm.DB
	SelectByEhid(fields []string, ehid string) *gorm.DB
	SelectByEhidOrderByStartDate(fields []string, ehid string, orderDir string) *gorm.DB
	SelectActiveByEhid(fields []string, ehid string) *gorm.DB
	ByEhidAndIntersectingDates(ehid string, startDate string, endDate string) *gorm.DB
	ByEhidAndEndDateIsNull(ehid string) *gorm.DB
}

type QueryImpl struct {
	Db        *gorm.DB
	TableName string
}

func NewQuery(db *gorm.DB, tableName string) Query {
	return &QueryImpl{
		Db:        db,
		TableName: tableName,
	}
}

func (q *QueryImpl) performSelect(fields []string) *gorm.DB {
	return q.Db.
		Select(fields).
		Table(q.TableName)
}

func (q *QueryImpl) SelectById(fields []string, id int) *gorm.DB {
	return q.performSelect(fields).
		Where("id = ?", id)
}

func (q *QueryImpl) SelectByEhid(fields []string, ehid string) *gorm.DB {
	return q.performSelect(fields).
		Where("ehid = ?", ehid)
}

func (q *QueryImpl) SelectByEhidOrderByStartDate(fields []string, ehid string, orderDir string) *gorm.DB {
	if orderDir == OrderNone {
		return q.SelectByEhid(fields, ehid)
	} else {
		return q.SelectByEhid(fields, ehid).
			Order("start_date " + orderDir)
	}
}

func (q *QueryImpl) SelectActiveByNodeId(fields []string, nodeId string) *gorm.DB {
	return q.performSelect(fields).
		Where("node_id = ?", nodeId).
		Where("start_date < NOW()").
		Where("end_date IS NULL OR end_date > NOW()")
}

func (q *QueryImpl) SelectActiveByEhid(fields []string, ehid string) *gorm.DB {
	return q.performSelect(fields).
		Where("ehid = ?", ehid).
		Where("start_date < NOW()").
		Where("end_date IS NULL OR end_date > NOW()")
}

func (q *QueryImpl) ByEhidAndIntersectingDates(ehid string, startDate string, endDate string) *gorm.DB {
	return q.Db.
		Table(q.TableName).
		Where("ehid = ?", ehid).
		Where(
			q.Db.
				Where("start_date <= ?", startDate).
				Where("end_date >= ?", startDate),
		).
		Or(
			q.Db.
				Where("start_date <= ?", endDate).
				Where("end_date >= ?", endDate),
		)
}

func (q *QueryImpl) ByEhidAndEndDateIsNull(ehid string) *gorm.DB {
	return q.Db.
		Table(q.TableName).
		Where("ehid = ?", ehid).
		Where("end_date IS NULL")
}
