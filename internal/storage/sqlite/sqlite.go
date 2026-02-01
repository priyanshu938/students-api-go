package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshu938/students-api/internal/config"
	"github.com/priyanshu938/students-api/internal/types"
)

type SQLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	);`)

	if err != nil {
		return nil, err
	}

	return &SQLite{Db: db}, nil
}

func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *SQLite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found", id)
		}

		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *SQLite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}


func (s *SQLite) UpdateStudentById(id int64, updates map[string]interface{}) error {
    if len(updates) == 0 {
        return fmt.Errorf("no fields to update")
    }

    query := "UPDATE students SET "
    args := []interface{}{}
    first := true

    for key, value := range updates {
        if !first {
            query += ", "
        }
        query += key + " = ?"
        args = append(args, value)
        first = false
    }

    query += " WHERE id = ?"
    args = append(args, id)

    stmt, err := s.Db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    result, err := stmt.Exec(args...)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("student with id %d not found", id)
    }

    return nil
}

func (s *SQLite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
