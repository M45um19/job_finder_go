package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, u *User) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUser := `INSERT INTO users (id, name, email, password, role, userphoto, refresh_token) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at`
	err = tx.QueryRow(ctx, queryUser, u.Id, u.Name, u.Email, u.Password, u.Role, u.UserPhoto, u.RefreshToken).Scan(&u.CreatedAt)
	if err != nil {
		return err
	}

	if u.Role == "employee" {
		queryProfile := `INSERT INTO seeker_profiles (user_id) VALUES ($1)`
		_, err = tx.Exec(ctx, queryProfile, u.Id)
	} else if u.Role == "employer" {
		queryProfile := `INSERT INTO employer_profiles (user_id) VALUES ($1)`
		_, err = tx.Exec(ctx, queryProfile, u.Id)
	}

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.role, u.userphoto, u.refresh_token, u.created_at, 
		       sp.address, sp.experience, sp.skills, sp.education, sp.certification,
		       ep.company_name, ep.company_desc, ep.address
		FROM users u
		LEFT JOIN seeker_profiles sp ON u.id = sp.user_id
		LEFT JOIN employer_profiles ep ON u.id = ep.user_id
		WHERE u.email = $1
	`

	user := &User{}
	var spAddress, spExperience, spSkills, spEducation, spCertification *string
	var epCompanyName, epCompanyDesc, epAddress *string

	err := ur.db.QueryRow(ctx, query, email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.UserPhoto, &user.RefreshToken, &user.CreatedAt,
		&spAddress, &spExperience, &spSkills, &spEducation, &spCertification,
		&epCompanyName, &epCompanyDesc, &epAddress,
	)
	if err != nil {
		return nil, err
	}

	if user.Role == "employee" {
		user.SeekerProfile = &SeekerProfile{
			Address:       spAddress,
			Experience:    spExperience,
			Skills:        spSkills,
			Education:     spEducation,
			Certification: spCertification,
		}
	} else if user.Role == "employer" {
		user.EmployerProfile = &EmployerProfile{
			CompanyName: epCompanyName,
			CompanyDesc: epCompanyDesc,
			Address:     epAddress,
		}
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.role, u.userphoto, u.refresh_token, u.created_at, 
		       sp.address, sp.experience, sp.skills, sp.education, sp.certification,
		       ep.company_name, ep.company_desc, ep.address
		FROM users u
		LEFT JOIN seeker_profiles sp ON u.id = sp.user_id
		LEFT JOIN employer_profiles ep ON u.id = ep.user_id
		WHERE u.id = $1
	`

	user := &User{}
	var spAddress, spExperience, spSkills, spEducation, spCertification *string
	var epCompanyName, epCompanyDesc, epAddress *string

	err := ur.db.QueryRow(ctx, query, id).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.UserPhoto, &user.RefreshToken, &user.CreatedAt,
		&spAddress, &spExperience, &spSkills, &spEducation, &spCertification,
		&epCompanyName, &epCompanyDesc, &epAddress,
	)
	if err != nil {
		return nil, err
	}

	if user.Role == "employee" {
		user.SeekerProfile = &SeekerProfile{
			Address:       spAddress,
			Experience:    spExperience,
			Skills:        spSkills,
			Education:     spEducation,
			Certification: spCertification,
		}
	} else if user.Role == "employer" {
		user.EmployerProfile = &EmployerProfile{
			CompanyName: epCompanyName,
			CompanyDesc: epCompanyDesc,
			Address:     epAddress,
		}
	}

	return user, nil
}

func (ur *UserRepository) UpdatePhoto(ctx context.Context, userId int64, photoURL string) error {
	query := `UPDATE users SET userphoto=$1, updated_at=NOW() WHERE id=$2`

	_, err := ur.db.Exec(ctx, query, photoURL, userId)
	return err
}

func (ur *UserRepository) UpdateRefreshToken(ctx context.Context, userId int64, token string) error {
	query := `UPDATE users SET refresh_token=$1, updated_at=NOW() WHERE id=$2`

	_, err := ur.db.Exec(ctx, query, token, userId)
	return err
}

func (ur *UserRepository) UpdateEmployeeProfile(ctx context.Context, userId int64, name string, address, experience, skills, education, certification *string) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE users SET name=$1, updated_at=NOW() WHERE id=$2`, name, userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE seeker_profiles SET address=$1, experience=$2, skills=$3, education=$4, certification=$5, updated_at=NOW() WHERE user_id=$6`, address, experience, skills, education, certification, userId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (ur *UserRepository) UpdateEmployerProfile(ctx context.Context, userId int64, name string, companyName, companyDesc, address *string) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE users SET name=$1, updated_at=NOW() WHERE id=$2`, name, userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE employer_profiles SET company_name=$1, company_desc=$2, address=$3, updated_at=NOW() WHERE user_id=$4`, companyName, companyDesc, address, userId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
