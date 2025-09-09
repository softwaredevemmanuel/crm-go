package models

import (
	"time"
	"github.com/google/uuid"


)

type TutorProfile struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID             uuid.UUID      `gorm:"type:uuid;unique;not null" json:"user_id"` // Foreign key to users table
    
    // Professional Identity
    FirstName          string         `gorm:"type:varchar(100);not null" json:"first_name"`
    LastName           string         `gorm:"type:varchar(100);not null" json:"last_name"`
    ProfilePicture     string         `gorm:"type:varchar(500)" json:"profile_picture,omitempty"`
    ProfessionalTitle  string         `gorm:"type:varchar(255)" json:"professional_title,omitempty"` // e.g., "PhD in Computer Science"
    Bio                string         `gorm:"type:text" json:"bio,omitempty"`
    
    // Contact & Professional Links
    PhoneNumber        string         `gorm:"type:varchar(50)" json:"phone_number,omitempty"`
    Email  string         `gorm:"type:varchar(255)" json:"professional_email,omitempty"`
    Website            string         `gorm:"type:varchar(500)" json:"website,omitempty"`
    PortfolioURL       string         `gorm:"type:varchar(500)" json:"portfolio_url,omitempty"`
    GitHubURL          string         `gorm:"type:varchar(500)" json:"github_url,omitempty"`
    LinkedInURL        string         `gorm:"type:varchar(500)" json:"linkedin_url,omitempty"`
    TwitterURL         string         `gorm:"type:varchar(500)" json:"twitter_url,omitempty"`
    YouTubeURL         string         `gorm:"type:varchar(500)" json:"youtube_url,omitempty"`
    
 
    // Professional Background
    YearsExperience   int            `gorm:"default:0" json:"years_experience"`
    YearsTeaching     int            `gorm:"default:0" json:"years_teaching"`
    ExpertiseAreas    []string       `gorm:"type:text[]" json:"expertise_areas,omitempty"` // e.g., ["web development", "data science"]
    Specializations   []string       `gorm:"type:text[]" json:"specializations,omitempty"` // e.g., ["React", "Python", "Machine Learning"]
    Certifications    []string       `gorm:"type:text[]" json:"certifications,omitempty"`
    Awards            []string       `gorm:"type:text[]" json:"awards,omitempty"`
    Publications      []string       `gorm:"type:text[]" json:"publications,omitempty"`
    
  
    // Performance Metrics
    TotalStudentsTaught int          `gorm:"default:0" json:"total_students_taught"`
    TotalCoursesCreated int          `gorm:"default:0" json:"total_courses_created"`
    TotalReviews        int          `gorm:"default:0" json:"total_reviews"`
    AverageRating       float64      `gorm:"type:decimal(3,2);default:0" json:"average_rating"`
    
	  // Metadata
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
    
    
    // Relationships
    User               User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Courses            []Course     `gorm:"foreignKey:TutorID" json:"courses,omitempty"`
}