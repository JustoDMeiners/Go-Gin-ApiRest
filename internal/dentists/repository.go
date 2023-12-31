package dentists

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/ncondezo/final/internal/domain"
)

var (
	ErrPrepareStatement = errors.New("error prepare statement")
	ErrExecStatement    = errors.New("error exec statement")
	ErrLastInsertedId   = errors.New("error last inserted id")
	ErrNotFound         = errors.New("error not found dentist")
	ErrAlreadyExists    = errors.New("error dentist already exists")
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Create is a method that creates a new dentist.
func (r *repository) Create(ctx context.Context, dentist domain.Dentist) (domain.Dentist, error) {
	var mysqlError *mysql.MySQLError

	statement, err := r.db.Prepare(QueryInsertDentist)
	if err != nil {
		return domain.Dentist{}, ErrPrepareStatement
	}
	defer statement.Close()

	result, err := statement.Exec(
		dentist.Name,
		dentist.LastName,
		dentist.Registration,
	)

	ok := errors.As(err, &mysqlError)
	if ok {
		switch mysqlError.Number {
		case 1062:
			return domain.Dentist{}, ErrAlreadyExists
		default:
			return domain.Dentist{}, ErrExecStatement
		}
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return domain.Dentist{}, ErrLastInsertedId
	}

	dentist.Id = int(lastId)

	return dentist, nil
}

// GetByID is a method that returns a dentist by ID.
func (r *repository) GetByID(ctx context.Context, id int) (domain.Dentist, error) {
	row := r.db.QueryRow(QueryGetDentistById, id)

	var dentist domain.Dentist
	err := row.Scan(
		&dentist.Id,
		&dentist.Name,
		&dentist.LastName,
		&dentist.Registration,
	)
	if err == sql.ErrNoRows {
		return domain.Dentist{}, ErrNotFound
	}
	if err != nil {
		return domain.Dentist{}, ErrExecStatement
	}

	return dentist, nil
}

// Update is a method that updates a dentist by ID.
func (r *repository) Update(ctx context.Context, dentist domain.Dentist, id int) (domain.Dentist, error) {
	statement, err := r.db.Prepare(QueryUpdateDentist)
	if err != nil {
		return domain.Dentist{}, ErrPrepareStatement
	}
	defer statement.Close()

	_, err = statement.Exec(
		dentist.Name,
		dentist.LastName,
		dentist.Registration,
		id,
	)

	if err != nil {
		return domain.Dentist{}, ErrExecStatement
	}

	return dentist, nil
}

// Patch is a method that updates a dentist registry by ID.
func (r *repository) Patch(ctx context.Context, dentist domain.Dentist, id int) (domain.Dentist, error) {
	statement, err := r.db.Prepare(QueryPatchDentist)
	if err != nil {
		return domain.Dentist{}, ErrPrepareStatement
	}
	defer statement.Close()

	_, err = statement.Exec(dentist.Registration, id)

	if err != nil {
		return domain.Dentist{}, ErrExecStatement
	}

	return dentist, nil
}

// Delete is a method that deletes a patient by ID.
func (r *repository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(QueryDeleteDentist, id)
	if err != nil {
		return ErrExecStatement
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return ErrNotFound
	}

	return nil
}
