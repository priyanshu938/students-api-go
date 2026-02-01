package storage

import "github.com/priyanshu938/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudentById(id int64, updates map[string]interface{}) error
	DeleteStudentById(id int64) error
}
