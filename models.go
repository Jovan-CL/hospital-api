package main

// Dummy data for feature testing
type Patient struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Condition string `json:"condition"`
}

type Staff struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type Role string

const (
	RoleAdmin        Role = "Admin"
	RoleDoctor       Role = "Doctor"
	RoleNurse        Role = "Nurse"
	RoleReceptionist Role = "Receptionist"
)

func (r Role) isValid() bool {
	switch r {
	case RoleAdmin, RoleDoctor, RoleNurse, RoleReceptionist:
		return true
	default:
		return false
	}
}

/*FOR OFFICIAL TESTING READY FOR USE IN PRODUCTION

type Patient struct {
	ID               string    `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	DateOfBirth      string    `json:"date_of_birth"` // Format: YYYY-MM-DD
	Gender           string    `json:"gender"`
	BloodType        string    `json:"blood_type"`
	AdmissionStatus  string    `json:"admission_status"`
	AcuityLevel      string    `json:"acuity_level"`
	RoomNumber       string    `json:"room_number,omitempty"` // omitempty hides it if blank
	AttendingDoctor  string    `json:"attending_doctor"`
	AdmissionDate    time.Time `json:"admission_date"`
	DischargeDate    time.Time `json:"discharge_date,omitempty"` // omitempty hides it if blank
	MedicalHistory   []string  `json:"medical_history,omitempty"` // omitempty hides it if blank
	CurrentMedications []string `json:"current_medications,omitempty"` // omitempty hides it if blank
}*/

/* FOR OFFICIAL TESTING READY FOR USE IN PRODUCTION STAFF ACCOUNT STRUCTURE

type Staff struct {
	// --- Authentication & System Fields ---
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // "-" completely hides this from all API JSON responses
	Role         string    `json:"role"` // e.g., "Doctor", "Nurse", "Admin", "Receptionist"
	IsActive     bool      `json:"is_active"` // For deactivating accounts without deleting history

	// --- Personal & Demographic Details ---
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth string    `json:"date_of_birth"` // Format: YYYY-MM-DD

	// --- Professional & Clinical Details ---
	Department     string    `json:"department"`      // e.g., "Emergency", "Pediatrics", "ICU"
	LicenseNumber  string    `json:"license_number,omitempty"` // Medical license ID (empty for non-medical staff)
	Specialization string    `json:"specialization,omitempty"` // e.g., "Cardiologist" (if doctor)

	// --- System Audit Logs ---
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
*/
