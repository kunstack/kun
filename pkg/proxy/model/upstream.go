package model

import (
	"fmt"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"time"
)

// Upstream used to describe upstream information
type Upstream struct {
	// ID in the ID database
	ID string `json:"id,omitempty" gorm:"type:varchar(36);primary_key;"`
	// CreatedAt database record creation time
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt the update time of the database record
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// DeletedAt the deletion time of the database record
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"`
	// ExpiredAt records the expiration time of the upstream
	ExpiredAt time.Time `json:"expiredAt,omitempty"`
	// Endpoints the address of the apiServer where the upstream is located
	Endpoint string `json:"endpoint,omitempty" gorm:"type:varchar(512)"`
	// DomainName the domain name for upstream
	DomainName string `json:"domainName,omitempty" gorm:"type:varchar(256),index:"`
}

// TableName database table name
func (*Upstream) TableName() string {
	return "upstreams"
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *Upstream) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = fmt.Sprintf("upstream-%s", xid.New())
	}
	return nil
}
