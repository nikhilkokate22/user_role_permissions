package repository

import (
	"user_role_permissions/model"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(app *model.ApplicationDetails) error
	ListApplications( limit, offset int, applicantName string) ([]model.ApplicationDetails, int64, int64, error) 
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicantRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(app *model.ApplicationDetails) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) ListApplications( limit, offset int, applicantName string) ([]model.ApplicationDetails, int64, int64, error) {

    var applications []model.ApplicationDetails
    var total int64
    var filtered int64

    baseQuery := r.db.Model(&model.ApplicationDetails{})

    // total count (no filter)
    if err := r.db.Model(&model.ApplicationDetails{}).
        Count(&total).Error; err != nil {
        return nil, 0, 0, err
    }

    // apply filter
    if applicantName != "" {
        baseQuery = baseQuery.Where(
            "applicant_name LIKE ?",
            "%"+applicantName+"%",
        )
    }

    // filtered count
    if err := baseQuery.Count(&filtered).Error; err != nil {
        return nil, 0, 0, err
    }

    if limit == 0 {
        limit = 10
    }

    if err := baseQuery.
		Preload("ProductionProfile").
        Preload("ProductionCenter").
        Preload("ProductionCenter.Facilities").
        Preload("TurnoverDetails").
        Preload("TurnoverDetails.Years").
        Preload("ProductApplied").
        Preload("PaymentDetails").
        Preload("Declaration").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&applications).Error; err != nil {
        return nil, 0, 0, err
    }

    return applications, total, filtered, nil
}