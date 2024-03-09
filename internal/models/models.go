package models

type Employee struct {
	Id   uint32 `db:"empl_id" gorm:"primaryKey;column:empl_id"`
	Name string
}

func (Employee) TableName() string { // for gorm
	return "employee"
}
