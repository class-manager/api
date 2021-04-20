package model

// Account represents a Classman account
type Account struct {
	BaseModel
	Email    string `gorm:"index:,unique; not null"`
	Password []byte `gorm:"index:; not null"`
	Name     string `gorm:"not null"`
}
