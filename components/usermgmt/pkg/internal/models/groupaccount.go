package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GroupAccount struct {
	bun.BaseModel `bun:"table:authsrv_groupaccount,alias:groupaccount"`

	ID          uuid.UUID `bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	ModifiedAt  time.Time `bun:"modified_at,notnull,default:current_timestamp"`
	Trash       bool      `bun:"trash,notnull,default:false"`
	// OrganizationId uuid.UUID `bun:"organization_id,type:uuid"`
	// PartnerId      uuid.UUID `bun:"partner_id,type:uuid"`
	AccountId uuid.UUID `bun:"account_id,type:uuid"`
	GroupId   uuid.UUID `bun:"group_id,type:uuid"`
	Active    bool      `bun:"active,notnull"`
}
