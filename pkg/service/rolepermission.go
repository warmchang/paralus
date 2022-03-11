package service

import (
	"context"

	"github.com/RafaySystems/rcloud-base/internal/models"
	"github.com/RafaySystems/rcloud-base/internal/persistence/provider/pg"
	"github.com/RafaySystems/rcloud-base/internal/utils"
	v3 "github.com/RafaySystems/rcloud-base/proto/types/commonpb/v3"
	rolev3 "github.com/RafaySystems/rcloud-base/proto/types/rolepb/v3"
	"github.com/google/uuid"
	bun "github.com/uptrace/bun"
)

const (
	rolepermissionKind     = "RolePermission"
	rolepermissionListKind = "RolePermissionList"
)

// RolepermissionService is the interface for rolepermission operations
type RolepermissionService interface {
	Close() error
	// get rolepermission by name
	GetByName(context.Context, *rolev3.RolePermission) (*rolev3.RolePermission, error)
	// list rolepermissions
	List(context.Context, *rolev3.RolePermission) (*rolev3.RolePermissionList, error)
}

// rolepermissionService implements RolepermissionService
type rolepermissionService struct {
	dao pg.EntityDAO
	l   utils.Lookup
}

// NewRolepermissionService return new rolepermission service
func NewRolepermissionService(db *bun.DB) RolepermissionService {
	return &rolepermissionService{
		dao: pg.NewEntityDAO(db),
		l:   utils.NewLookup(db),
	}
}

func (s *rolepermissionService) toV3Rolepermission(rolepermission *rolev3.RolePermission, rlp *models.ResourcePermission) *rolev3.RolePermission {
	// TODO: should we return resource_urls?
	rolepermission.Metadata = &v3.Metadata{
		Name:        rlp.Name,
		Description: rlp.Description,
	}

	return rolepermission
}

func (s *rolepermissionService) getPartnerOrganization(ctx context.Context, rolepermission *rolev3.RolePermission) (uuid.UUID, uuid.UUID, error) {
	partner := rolepermission.GetMetadata().GetPartner()
	org := rolepermission.GetMetadata().GetOrganization()
	partnerId, err := s.l.GetPartnerId(ctx, partner)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	organizationId, err := s.l.GetOrganizationId(ctx, org)
	if err != nil {
		return partnerId, uuid.Nil, err
	}
	return partnerId, organizationId, nil

}

func (s *rolepermissionService) GetByName(ctx context.Context, rolepermission *rolev3.RolePermission) (*rolev3.RolePermission, error) {
	name := rolepermission.GetMetadata().GetName()
	entity, err := s.dao.GetByName(ctx, name, &models.ResourcePermission{})
	if err != nil {
		return rolepermission, err
	}

	if rle, ok := entity.(*models.ResourcePermission); ok {
		rolepermission = s.toV3Rolepermission(rolepermission, rle)

		return rolepermission, nil
	}
	return rolepermission, nil

}

func (s *rolepermissionService) List(ctx context.Context, rolepermission *rolev3.RolePermission) (*rolev3.RolePermissionList, error) {
	var rolepermissions []*rolev3.RolePermission
	rolepermissionList := &rolev3.RolePermissionList{
		ApiVersion: apiVersion,
		Kind:       rolepermissionListKind,
		Metadata: &v3.ListMetadata{
			Count: 0,
		},
	}
	var rles []models.ResourcePermission
	entities, err := s.dao.List(ctx, uuid.NullUUID{UUID: uuid.Nil, Valid: false}, uuid.NullUUID{UUID: uuid.Nil, Valid: false}, &rles)
	if err != nil {
		return rolepermissionList, err
	}
	if rles, ok := entities.(*[]models.ResourcePermission); ok {
		for _, rle := range *rles {
			entry := &rolev3.RolePermission{Metadata: rolepermission.GetMetadata()}
			entry = s.toV3Rolepermission(entry, &rle)
			rolepermissions = append(rolepermissions, entry)
		}

		//update the list metadata and items response
		rolepermissionList.Metadata = &v3.ListMetadata{
			Count: int64(len(rolepermissions)),
		}
		rolepermissionList.Items = rolepermissions
	}

	return rolepermissionList, nil
}

func (s *rolepermissionService) Close() error {
	return s.dao.Close()
}