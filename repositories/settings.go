package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
)

func (r *Repository) GetCompanySettings(settings *models.CompanySettings) error {

	err := r.DB.Get(settings, `SELECT * FROM Tbl_Company_Settings LIMIT 1`)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateCompanySettings(tx *sqlx.Tx, input models.CompanyField, logoPath string) error {
	_, err := tx.Exec(`
        UPDATE Tbl_Company_Settings
        SET working_days_per_month=$1, allow_manager_add_leave=$2, company_name = $3, 
		    primary_color = $4, 
		    secondary_color = $5, logo_path = COALESCE(NULLIF($6, ''), logo_path), updated_at=NOW()
    `, input.WorkingDaysPerMonth, input.AllowManagerAddLeave, input.CompanyName, // New field
		input.PrimaryColor, // New field
		input.SecondaryColor,
		logoPath, // New field
	)

	if err != nil {
		return err
	}
	return nil
}
