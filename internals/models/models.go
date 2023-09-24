package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"primary key;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

type Company struct {
	Base
	Name      string `json:"name" gorm:"not null;index;unique"`
	Employees uint   `json:"employees" gorm:"not null;default:1;type:smallint"`
	IsActive  bool   `json:"is_active" gorm:"not null;default:true"`
}

type User struct {
	Base
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Name     string `json:"name" gorm:"index"`

	CompanyId uuid.UUID `json:"companyId" gorm:"not null;type:uuid"`
	Company   Company   `json:"company" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

type Project struct {
	Base
	Name     string `json:"name" gorm:"not null;index;index:uq_projectname,unique"`
	IsActive bool   `json:"isActive" gorm:"not null;default:true"`

	Company   Company   `json:"company" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CompanyId uuid.UUID `json:"companyId" gorm:"not null;type:uuid;index:uq_projectname,unique"`
}

type Supplier struct {
	Base
	SupplierId   string `json:"supplier_id" gorm:"not null;index:uq_supplierid,unique"`
	Name         string `json:"name" gorm:"not null;index:uq_suppliername,unique;index"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`

	Company   Company   `json:"company" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CompanyId uuid.UUID `json:"companyId" gorm:"not null;type:uuid;index:uq_supplierid,unique;index:uq_suppliername,unique"`
}

type BudgetItem struct {
	Base
	Code        string `json:"code" gorm:"not null;index;index:uq_budgetitemcode,unique"`
	Name        string `json:"name" gorm:"not null;index;index:uq_budgetitemname,unique"`
	Accumulates bool   `json:"accumulates" gorm:"not null;default:false"`
	Level       uint   `json:"level" gorm:"not null;type:uint;default:1"`

	BudgetItem   *BudgetItem `json:"budgetItem" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	BudgetItemId uuid.UUID   `json:"budgetItemId" gorm:"type:uuid"`

	Company   Company   `json:"company" gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CompanyId uuid.UUID `json:"companyId" gorm:"not null;type:uuid;index:uq_budgetitemcode,unique;index:uq_budgetitemname,unique"`
}
