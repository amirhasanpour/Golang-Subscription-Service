package data

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// User is the structure which holds one user from the database.
type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
    	id, 
       	email, 
       	first_name, 
       	last_name, 
       	password, 
       	user_active, 
       	is_admin, 
       	created_at, 
       	updated_at
	from 
	    users 
	order by 
	    last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
			select 
			    id, 
			    email, 
			    first_name, 
			    last_name, 
			    password, 
			    user_active, 
			    is_admin, 
			    created_at, 
			    updated_at 
			from 
			    users 
			where 
			    email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get plan, if any
	query = `select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at from 
			plans p
			left join user_plans up on (p.id = up.plan_id)
			where up.user_id = $1`

	var plan Plan
	row = db.QueryRowContext(ctx, query, user.ID)

	err = row.Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.PlanAmount,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err == nil {
		user.Plan = &plan
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, is_admin, created_at, updated_at 
				from users 
				where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get plan, if any
	query = `select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at from 
			plans p
			left join user_plans up on (p.id = up.plan_id)
			where up.user_id = $1`

	var plan Plan
	row = db.QueryRowContext(ctx, query, user.ID)

	err = row.Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.PlanAmount,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err == nil {
		user.Plan = &plan
	} else {
		log.Println("Error getting plan", err)
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6`

	_, err := db.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string, userPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
