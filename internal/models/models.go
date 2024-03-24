package models

type Employee struct {
	tableName struct{} `pg:"employee"` // for go-pg
	Id        uint32   `db:"empl_id" gorm:"primaryKey;column:empl_id" pg:"empl_id,pk"`
	//Id   uint32 `pg:"empl_id,pk"` // for go-pg
	Name string
}

func (Employee) TableName() string { // for gorm
	return "employee"
}
