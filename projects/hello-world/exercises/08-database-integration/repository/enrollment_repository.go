package repository

import (
	"errors"
	"time"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"gorm.io/gorm"
)

// =============================================================================
// ENROLLMENT REPOSITORY - Many-to-Many with Join Table Attributes
// =============================================================================
// When your Many-to-Many relationship needs extra attributes (like enrollment date, grade),
// you need an explicit join entity. This is common in real-world applications.
//
// Java equivalent:
// @Entity class Enrollment {
//     @ManyToOne Student student;
//     @ManyToOne Course course;
//     LocalDate enrolledAt;
//     String grade;
// }

// EnrollmentRepository defines operations for Student-Course enrollment
type EnrollmentRepository interface {
	// Enrollment operations
	Enroll(studentID, courseID uint) (*model.Enrollment, error)
	Unenroll(studentID, courseID uint) error
	FindByStudentAndCourse(studentID, courseID uint) (*model.Enrollment, error)
	UpdateGrade(studentID, courseID uint, grade string) error
	MarkCompleted(studentID, courseID uint) error

	// Query operations
	FindByStudent(studentID uint) ([]model.Enrollment, error)
	FindByCourse(courseID uint) ([]model.Enrollment, error)
	FindCompletedByStudent(studentID uint) ([]model.Enrollment, error)
	CountStudentsInCourse(courseID uint) (int64, error)
	GetAverageGradeForCourse(courseID uint) (float64, error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new EnrollmentRepository
func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

// =============================================================================
// ENROLLMENT OPERATIONS
// =============================================================================

// Enroll creates a new enrollment (student joins a course)
// Java: enrollmentRepository.save(new Enrollment(student, course, LocalDate.now()));
func (r *enrollmentRepository) Enroll(studentID, courseID uint) (*model.Enrollment, error) {
	// First, check if already enrolled (unique constraint on student_id + course_id)
	var existing model.Enrollment
	err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&existing).Error
	if err == nil {
		return nil, errors.New("student is already enrolled in this course")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if course has capacity (demonstrate transaction with business logic)
	var course model.Course
	if err := r.db.First(&course, courseID).Error; err != nil {
		return nil, errors.New("course not found")
	}

	var enrolledCount int64
	r.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&enrolledCount)
	if int(enrolledCount) >= course.MaxStudents {
		return nil, errors.New("course is full")
	}

	// Create enrollment
	enrollment := &model.Enrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		EnrolledAt: time.Now(),
		Completed:  false,
	}

	if err := r.db.Create(enrollment).Error; err != nil {
		return nil, err
	}

	return enrollment, nil
}

// Unenroll removes a student from a course
// Java: enrollmentRepository.deleteByStudentIdAndCourseId(studentId, courseId);
func (r *enrollmentRepository) Unenroll(studentID, courseID uint) error {
	result := r.db.
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Delete(&model.Enrollment{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("enrollment not found")
	}
	return nil
}

// FindByStudentAndCourse finds a specific enrollment
// Java: Optional<Enrollment> findByStudentIdAndCourseId(Long studentId, Long courseId);
func (r *enrollmentRepository) FindByStudentAndCourse(studentID, courseID uint) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.db.
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Preload("Student").
		Preload("Course").
		First(&enrollment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &enrollment, err
}

// UpdateGrade sets the grade for an enrollment
// Java: @Modifying @Query("UPDATE Enrollment e SET e.grade = :grade WHERE e.student.id = :sid AND e.course.id = :cid")
func (r *enrollmentRepository) UpdateGrade(studentID, courseID uint, grade string) error {
	// Model() + Where() + Update() for targeted updates
	// More efficient than fetching, modifying, saving
	result := r.db.
		Model(&model.Enrollment{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Update("grade", grade)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("enrollment not found")
	}
	return nil
}

// MarkCompleted marks an enrollment as completed
func (r *enrollmentRepository) MarkCompleted(studentID, courseID uint) error {
	result := r.db.
		Model(&model.Enrollment{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Update("completed", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("enrollment not found")
	}
	return nil
}

// =============================================================================
// QUERY OPERATIONS
// =============================================================================

// FindByStudent gets all enrollments for a student
// Java: List<Enrollment> findByStudentId(Long studentId);
func (r *enrollmentRepository) FindByStudent(studentID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	err := r.db.
		Where("student_id = ?", studentID).
		Preload("Course"). // Include course details
		Find(&enrollments).Error
	return enrollments, err
}

// FindByCourse gets all enrollments for a course
// Java: List<Enrollment> findByCourseId(Long courseId);
func (r *enrollmentRepository) FindByCourse(courseID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	err := r.db.
		Where("course_id = ?", courseID).
		Preload("Student"). // Include student details
		Find(&enrollments).Error
	return enrollments, err
}

// FindCompletedByStudent gets completed courses for a student
// Java: List<Enrollment> findByStudentIdAndCompletedTrue(Long studentId);
func (r *enrollmentRepository) FindCompletedByStudent(studentID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	err := r.db.
		Where("student_id = ? AND completed = ?", studentID, true).
		Preload("Course").
		Find(&enrollments).Error
	return enrollments, err
}

// CountStudentsInCourse counts how many students are enrolled in a course
// Java: @Query("SELECT COUNT(e) FROM Enrollment e WHERE e.course.id = :courseId")
func (r *enrollmentRepository) CountStudentsInCourse(courseID uint) (int64, error) {
	var count int64
	err := r.db.
		Model(&model.Enrollment{}).
		Where("course_id = ?", courseID).
		Count(&count).Error
	return count, err
}

// GetAverageGradeForCourse calculates the average grade for a course
// Demonstrates raw SQL for complex calculations
// Java: @Query("SELECT AVG(CASE WHEN e.grade = 'A' THEN 4.0 ... END) FROM Enrollment e WHERE e.course.id = :cid")
func (r *enrollmentRepository) GetAverageGradeForCourse(courseID uint) (float64, error) {
	var result struct {
		Average float64
	}
	// Raw SQL for complex calculations
	// GORM supports raw queries when the query builder isn't sufficient
	err := r.db.Raw(`
		SELECT AVG(
			CASE grade
				WHEN 'A' THEN 4.0
				WHEN 'B' THEN 3.0
				WHEN 'C' THEN 2.0
				WHEN 'D' THEN 1.0
				WHEN 'F' THEN 0.0
				ELSE NULL
			END
		) as average
		FROM enrollments
		WHERE course_id = ? AND grade IS NOT NULL AND deleted_at IS NULL
	`, courseID).Scan(&result).Error

	return result.Average, err
}
