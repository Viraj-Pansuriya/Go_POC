package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/database"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/handler"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/repository"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/service"
)

// =============================================================================
// MAIN - Application Entry Point
// =============================================================================
// In Java/Spring: @SpringBootApplication + main()
// Go: Manual wiring of dependencies (no magic annotations)

func main() {
	fmt.Println("🚀 Starting Database Integration Exercise...")
	fmt.Println("=" + repeatChar('=', 59))

	// ==========================================================================
	// STEP 1: Connect to Database
	// ==========================================================================
	// Java equivalent: Spring Boot auto-configures DataSource from application.yml
	log.Println("📦 Connecting to database...")
	if err := database.Connect(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// ==========================================================================
	// STEP 2: Run Migrations (Create Tables)
	// ==========================================================================
	// Java equivalent: Hibernate ddl-auto or Flyway migrations
	log.Println("📋 Running migrations...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	// ==========================================================================
	// STEP 3: Seed Sample Data
	// ==========================================================================
	log.Println("🌱 Seeding sample data...")
	seedData()

	// ==========================================================================
	// STEP 4: Demonstrate Various Operations
	// ==========================================================================
	log.Println("\n" + repeatChar('=', 60))
	log.Println("📚 DEMONSTRATING DATABASE OPERATIONS")
	log.Println(repeatChar('=', 60))

	demonstrateOperations()

	// ==========================================================================
	// STEP 5: Start HTTP Server
	// ==========================================================================
	log.Println("\n" + repeatChar('=', 60))
	log.Println("🌐 STARTING HTTP SERVER")
	log.Println(repeatChar('=', 60))

	startServer()
}

// =============================================================================
// SEED DATA - Create initial test data
// =============================================================================
func seedData() {
	db := database.GetDB()

	// Check if data already exists
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("📝 Data already seeded, skipping...")
		return
	}

	// --------------------------------------------------------------------------
	// Create Users with Profiles (One-to-One)
	// --------------------------------------------------------------------------
	users := []model.User{
		{
			Name:  "Viraj Pansuriya",
			Email: "viraj@example.com",
			Age:   25,
			Profile: model.Profile{
				Bio:       "Go learner from Java/C++ background",
				AvatarURL: "https://example.com/viraj.jpg",
				Website:   "https://viraj.dev",
			},
		},
		{
			Name:  "Alice Developer",
			Email: "alice@example.com",
			Age:   28,
			Profile: model.Profile{
				Bio:     "Full-stack developer",
				Website: "https://alice.dev",
			},
		},
		{
			Name:  "Bob Engineer",
			Email: "bob@example.com",
			Age:   32,
			Profile: model.Profile{
				Bio: "Backend specialist",
			},
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("⚠️ Failed to create user %s: %v", user.Name, err)
		}
	}
	log.Println("✅ Created users with profiles")

	// --------------------------------------------------------------------------
	// Create Posts for Users (One-to-Many)
	// --------------------------------------------------------------------------
	posts := []model.Post{
		{Title: "My Go Journey", Content: "Learning Go from Java...", UserID: 1},
		{Title: "GORM vs Hibernate", Content: "Comparing ORMs...", UserID: 1},
		{Title: "Concurrency in Go", Content: "Goroutines are amazing!", UserID: 2},
		{Title: "REST APIs with Gin", Content: "Building APIs...", UserID: 2},
	}

	for _, post := range posts {
		db.Create(&post)
	}
	log.Println("✅ Created posts")

	// --------------------------------------------------------------------------
	// Create Authors and Books with Tags (One-to-Many + Many-to-Many)
	// --------------------------------------------------------------------------
	// First, create tags
	tags := []model.Tag{
		{Name: "Programming", Color: "#3498db"},
		{Name: "Fiction", Color: "#e74c3c"},
		{Name: "Science", Color: "#2ecc71"},
		{Name: "History", Color: "#f39c12"},
		{Name: "Go", Color: "#00ADD8"},
	}
	for _, tag := range tags {
		db.Create(&tag)
	}
	log.Println("✅ Created tags")

	// Create authors with books
	birthDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	authors := []model.Author{
		{
			Name:      "Robert C. Martin",
			Country:   "USA",
			BirthDate: &birthDate,
			Books: []model.Book{
				{Title: "Clean Code", ISBN: "978-0132350884", Price: 39.99},
				{Title: "Clean Architecture", ISBN: "978-0134494166", Price: 34.99},
			},
		},
		{
			Name:    "Donovan & Kernighan",
			Country: "USA",
			Books: []model.Book{
				{Title: "The Go Programming Language", ISBN: "978-0134190440", Price: 44.99},
			},
		},
	}

	for _, author := range authors {
		db.Create(&author)
	}
	log.Println("✅ Created authors with books")

	// Associate tags with books (Many-to-Many)
	var cleanCode model.Book
	db.Where("title = ?", "Clean Code").First(&cleanCode)
	var programmingTag, goTag model.Tag
	db.Where("name = ?", "Programming").First(&programmingTag)
	db.Where("name = ?", "Go").First(&goTag)
	db.Model(&cleanCode).Association("Tags").Append(&programmingTag)

	var goBook model.Book
	db.Where("title = ?", "The Go Programming Language").First(&goBook)
	db.Model(&goBook).Association("Tags").Append(&programmingTag, &goTag)
	log.Println("✅ Associated tags with books")

	// --------------------------------------------------------------------------
	// Create Students, Courses, and Enrollments (Many-to-Many with attributes)
	// --------------------------------------------------------------------------
	students := []model.Student{
		{Name: "John Doe", StudentCode: "STU001", Email: "john@university.edu"},
		{Name: "Jane Smith", StudentCode: "STU002", Email: "jane@university.edu"},
		{Name: "Mike Johnson", StudentCode: "STU003", Email: "mike@university.edu"},
	}
	for _, student := range students {
		db.Create(&student)
	}
	log.Println("✅ Created students")

	courses := []model.Course{
		{Name: "Go Programming", Code: "CS101", Credits: 3, MaxStudents: 30, Price: 299.99},
		{Name: "Database Systems", Code: "CS201", Credits: 4, MaxStudents: 25, Price: 349.99},
		{Name: "Web Development", Code: "CS301", Credits: 3, MaxStudents: 35, Price: 279.99},
	}
	for _, course := range courses {
		db.Create(&course)
	}
	log.Println("✅ Created courses")

	// Create enrollments with grades
	gradeA := "A"
	gradeB := "B"
	enrollments := []model.Enrollment{
		{StudentID: 1, CourseID: 1, EnrolledAt: time.Now().AddDate(0, -2, 0), Grade: &gradeA, Completed: true},
		{StudentID: 1, CourseID: 2, EnrolledAt: time.Now().AddDate(0, -1, 0), Grade: &gradeB, Completed: false},
		{StudentID: 2, CourseID: 1, EnrolledAt: time.Now().AddDate(0, -2, 0), Grade: &gradeA, Completed: true},
		{StudentID: 2, CourseID: 3, EnrolledAt: time.Now(), Completed: false},
		{StudentID: 3, CourseID: 2, EnrolledAt: time.Now().AddDate(0, -1, 0), Completed: false},
	}
	for _, enrollment := range enrollments {
		db.Create(&enrollment)
	}
	log.Println("✅ Created enrollments")

	// --------------------------------------------------------------------------
	// Create Category Tree (Self-referential)
	// --------------------------------------------------------------------------
	techCategory := model.Category{Name: "Technology", Description: "Tech topics"}
	db.Create(&techCategory)

	subCategories := []model.Category{
		{Name: "Programming", Description: "Programming languages", ParentID: &techCategory.ID},
		{Name: "Databases", Description: "Database systems", ParentID: &techCategory.ID},
	}
	for _, cat := range subCategories {
		db.Create(&cat)
	}
	log.Println("✅ Created category tree")

	fmt.Println("\n✅ Sample data seeded successfully!")
}

// =============================================================================
// DEMONSTRATE OPERATIONS
// =============================================================================
func demonstrateOperations() {
	db := database.GetDB()

	// --------------------------------------------------------------------------
	// 1. One-to-One: User with Profile
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 ONE-TO-ONE: User with Profile")
	fmt.Println(repeatChar('-', 40))

	var userWithProfile model.User
	db.Preload("Profile").First(&userWithProfile, 1)
	fmt.Printf("User: %s\n", userWithProfile.Name)
	fmt.Printf("Bio: %s\n", userWithProfile.Profile.Bio)
	fmt.Printf("Website: %s\n", userWithProfile.Profile.Website)

	// --------------------------------------------------------------------------
	// 2. One-to-Many: User with Posts
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 ONE-TO-MANY: User with Posts")
	fmt.Println(repeatChar('-', 40))

	var userWithPosts model.User
	db.Preload("Posts").First(&userWithPosts, 1)
	fmt.Printf("User: %s has %d posts:\n", userWithPosts.Name, len(userWithPosts.Posts))
	for _, post := range userWithPosts.Posts {
		fmt.Printf("  - %s\n", post.Title)
	}

	// --------------------------------------------------------------------------
	// 3. One-to-Many: Author with Books
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 ONE-TO-MANY: Author with Books")
	fmt.Println(repeatChar('-', 40))

	var author model.Author
	db.Preload("Books").First(&author, 1)
	fmt.Printf("Author: %s has %d books:\n", author.Name, len(author.Books))
	for _, book := range author.Books {
		fmt.Printf("  - %s ($%.2f)\n", book.Title, book.Price)
	}

	// --------------------------------------------------------------------------
	// 4. Many-to-Many: Book with Tags
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 MANY-TO-MANY: Book with Tags")
	fmt.Println(repeatChar('-', 40))

	var bookWithTags model.Book
	db.Preload("Tags").Where("title = ?", "The Go Programming Language").First(&bookWithTags)
	fmt.Printf("Book: %s\n", bookWithTags.Title)
	fmt.Printf("Tags: ")
	for i, tag := range bookWithTags.Tags {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s", tag.Name)
	}
	fmt.Println()

	// --------------------------------------------------------------------------
	// 5. Many-to-Many with Join Attributes: Student Enrollments
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 MANY-TO-MANY WITH ATTRIBUTES: Student Enrollments")
	fmt.Println(repeatChar('-', 40))

	var studentEnrollments []model.Enrollment
	db.Preload("Student").Preload("Course").Where("student_id = ?", 1).Find(&studentEnrollments)
	if len(studentEnrollments) > 0 {
		fmt.Printf("Student: %s is enrolled in:\n", studentEnrollments[0].Student.Name)
		for _, e := range studentEnrollments {
			grade := "Not graded"
			if e.Grade != nil {
				grade = *e.Grade
			}
			status := "In Progress"
			if e.Completed {
				status = "Completed"
			}
			fmt.Printf("  - %s (Grade: %s, Status: %s)\n", e.Course.Name, grade, status)
		}
	}

	// --------------------------------------------------------------------------
	// 6. Self-referential: Category Tree
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 SELF-REFERENTIAL: Category Tree")
	fmt.Println(repeatChar('-', 40))

	var parentCategory model.Category
	db.Preload("Children").First(&parentCategory, 1)
	fmt.Printf("Category: %s\n", parentCategory.Name)
	fmt.Printf("Children:\n")
	for _, child := range parentCategory.Children {
		fmt.Printf("  - %s\n", child.Name)
	}

	// --------------------------------------------------------------------------
	// 7. Complex Query: Books by Tag
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 COMPLEX QUERY: Find books by tag name")
	fmt.Println(repeatChar('-', 40))

	bookRepo := repository.NewBookRepository(db)
	books, _ := bookRepo.FindByTagName("Programming")
	fmt.Printf("Books tagged 'Programming': %d\n", len(books))
	for _, book := range books {
		fmt.Printf("  - %s\n", book.Title)
	}

	// --------------------------------------------------------------------------
	// 8. Aggregation: Course Statistics
	// --------------------------------------------------------------------------
	fmt.Println("\n📌 AGGREGATION: Course Statistics")
	fmt.Println(repeatChar('-', 40))

	enrollmentRepo := repository.NewEnrollmentRepository(db)
	count, _ := enrollmentRepo.CountStudentsInCourse(1)
	avgGrade, _ := enrollmentRepo.GetAverageGradeForCourse(1)
	fmt.Printf("Course 'Go Programming':\n")
	fmt.Printf("  - Enrolled students: %d\n", count)
	fmt.Printf("  - Average GPA: %.2f\n", avgGrade)
}

// =============================================================================
// START HTTP SERVER
// =============================================================================
func startServer() {
	// Get database instance
	db := database.GetDB()

	// --------------------------------------------------------------------------
	// Wire up dependencies (Manual Dependency Injection)
	// --------------------------------------------------------------------------
	// In Spring: @Autowired handles this automatically
	// Go: We wire dependencies explicitly (clearer, no magic)

	// Create repositories
	userRepo := repository.NewUserRepository(db)

	// Create services with injected repositories
	userService := service.NewUserService(userRepo)

	// Create handlers with injected services
	userHandler := handler.NewUserHandler(userService)

	// --------------------------------------------------------------------------
	// Setup Gin Router
	// --------------------------------------------------------------------------
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Register routes
	userHandler.RegisterRoutes(r)

	// --------------------------------------------------------------------------
	// Print API documentation
	// --------------------------------------------------------------------------
	fmt.Println("\n📍 API Endpoints:")
	fmt.Println("  POST   /api/users           - Register new user")
	fmt.Println("  GET    /api/users           - Get all users (supports pagination)")
	fmt.Println("  GET    /api/users/:id       - Get user by ID")
	fmt.Println("  PUT    /api/users/:id       - Update user")
	fmt.Println("  DELETE /api/users/:id       - Delete user (soft delete)")
	fmt.Println("  GET    /api/users/search?q= - Search users")
	fmt.Println("  PUT    /api/users/:id/profile - Update profile")

	fmt.Println("\n🚀 Server starting on http://localhost:8080")
	fmt.Println("📝 Try: curl http://localhost:8080/api/users")

	// --------------------------------------------------------------------------
	// Start Server
	// --------------------------------------------------------------------------
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Helper function to repeat a character
func repeatChar(char rune, count int) string {
	result := make([]rune, count)
	for i := range result {
		result[i] = char
	}
	return string(result)
}

