package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RafaySystems/rcloud-base/components/adminsrv/pkg/internal/models"
	systemv3 "github.com/RafaySystems/rcloud-base/components/adminsrv/proto/types/systempb/v3"
	"github.com/RafaySystems/rcloud-base/components/common/pkg/persistence/provider/pg"
	v3 "github.com/RafaySystems/rcloud-base/components/common/proto/types/commonpb/v3"
	"github.com/google/uuid"
	bun "github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	organizationKind     = "Organization"
	organizationListKind = "OrganizationList"
)

// OrganizationService is the interface for organization operations
type OrganizationService interface {
	Close() error
	// create organization
	Create(ctx context.Context, organization *systemv3.Organization) (*systemv3.Organization, error)
	// get organization by id
	GetByID(ctx context.Context, id string) (*systemv3.Organization, error)
	// get organization by id
	GetByName(ctx context.Context, name string) (*systemv3.Organization, error)
	// create or update organization
	Update(ctx context.Context, organization *systemv3.Organization) (*systemv3.Organization, error)
	// delete organization
	Delete(ctx context.Context, organization *systemv3.Organization) (*systemv3.Organization, error)
	// list organization
	List(ctx context.Context, organization *systemv3.Organization) (*systemv3.OrganizationList, error)
}

// organizationService implements OrganizationService
type organizationService struct {
	dao pg.EntityDAO
}

// NewOrganizationService return new organization service
func NewOrganizationService(db *bun.DB) OrganizationService {
	return &organizationService{
		dao: pg.NewEntityDAO(db),
	}
}

func (s *organizationService) Create(ctx context.Context, org *systemv3.Organization) (*systemv3.Organization, error) {

	partnerId, _ := uuid.Parse(org.GetMetadata().GetPartner())

	//update default organization setting values
	org.Spec.Settings = &systemv3.OrganizationSettings{
		Lockout: &systemv3.Lockout{
			Enabled:   true,
			PeriodMin: 15,
			Attempts:  5,
		},
		IdleLogoutMin: 60,
	}
	sb, err := json.MarshalIndent(org.GetSpec().GetSettings(), "", "\t")
	if err != nil {
		org.Status = &v3.Status{
			ConditionType:   "Create",
			ConditionStatus: v3.ConditionStatus_StatusFailed,
			LastUpdated:     timestamppb.Now(),
			Reason:          err.Error(),
		}
	}
	//convert v3 spec to internal models
	organization := models.Organization{
		Name:              org.GetMetadata().GetName(),
		Description:       org.GetMetadata().GetDescription(),
		Trash:             false,
		Settings:          json.RawMessage(sb),
		BillingAddress:    org.GetSpec().GetBillingAddress(),
		PartnerId:         partnerId,
		Active:            org.GetSpec().GetActive(),
		Approved:          org.GetSpec().GetApproved(),
		Type:              org.GetSpec().GetType(),
		AddressLine1:      org.GetSpec().GetAddressLine1(),
		AddressLine2:      org.GetSpec().GetAddressLine2(),
		City:              org.GetSpec().GetCity(),
		Country:           org.GetSpec().GetCountry(),
		Phone:             org.GetSpec().GetPhone(),
		State:             org.GetSpec().GetState(),
		Zipcode:           org.GetSpec().GetZipcode(),
		IsPrivate:         org.GetSpec().GetIsPrivate(),
		IsTOTPEnabled:     org.GetSpec().GetIsTotpEnabled(),
		AreClustersShared: org.GetSpec().GetAreClustersShared(),
		CreatedAt:         time.Now(),
		ModifiedAt:        time.Now(),
	}
	entity, err := s.dao.Create(ctx, &organization)
	if err != nil {
		org.Status = &v3.Status{
			ConditionType:   "Create",
			ConditionStatus: v3.ConditionStatus_StatusFailed,
			LastUpdated:     timestamppb.Now(),
			Reason:          err.Error(),
		}
		return org, err
	}

	if createdOrg, ok := entity.(*models.Organization); ok {
		//update v3 spec
		org.Metadata.Id = createdOrg.ID.String()
		org.Status = &v3.Status{
			ConditionType:   "Create",
			ConditionStatus: v3.ConditionStatus_StatusOK,
			LastUpdated:     timestamppb.Now(),
		}
	}

	return org, nil
}

func (s *organizationService) GetByID(ctx context.Context, id string) (*systemv3.Organization, error) {

	organization := &systemv3.Organization{
		ApiVersion: apiVersion,
		Kind:       organizationKind,
		Metadata: &v3.Metadata{
			Id: id,
		},
		Spec: &systemv3.OrganizationSpec{},
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		organization.Status = &v3.Status{
			ConditionType:   "Describe",
			LastUpdated:     timestamppb.Now(),
			ConditionStatus: v3.ConditionStatus_StatusFailed,
		}
		return organization, err
	}
	entity, err := s.dao.GetByID(ctx, uid, &models.Organization{})
	if err != nil {
		organization.Status = &v3.Status{
			ConditionType:   "Describe",
			LastUpdated:     timestamppb.Now(),
			ConditionStatus: v3.ConditionStatus_StatusFailed,
		}
		return organization, err
	}

	if org, ok := entity.(*models.Organization); ok {

		var settings systemv3.OrganizationSettings
		err := json.Unmarshal(org.Settings, &settings)
		if err != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Describe",
				LastUpdated:     timestamppb.Now(),
				ConditionStatus: v3.ConditionStatus_StatusFailed,
			}
			return organization, err
		}

		organization.Metadata = &v3.Metadata{
			Name:        org.Name,
			Description: org.Description,
			Id:          org.ID.String(),
			Partner:     org.PartnerId.String(),
			ModifiedAt:  timestamppb.New(org.ModifiedAt),
		}
		organization.Spec = &systemv3.OrganizationSpec{
			BillingAddress:    org.BillingAddress,
			Active:            org.Active,
			Approved:          org.Approved,
			Type:              org.Type,
			AddressLine1:      org.AddressLine1,
			AddressLine2:      org.AddressLine2,
			City:              org.City,
			Country:           org.Country,
			Phone:             org.Phone,
			State:             org.State,
			Zipcode:           org.Zipcode,
			IsPrivate:         org.IsPrivate,
			IsTotpEnabled:     org.IsTOTPEnabled,
			AreClustersShared: org.AreClustersShared,
			Settings:          &settings,
		}
		organization.Status = &v3.Status{
			ConditionType:   "Describe",
			LastUpdated:     timestamppb.Now(),
			ConditionStatus: v3.ConditionStatus_StatusOK,
		}

		return organization, nil

	} else {
		organization := &systemv3.Organization{
			ApiVersion: apiVersion,
			Kind:       organizationKind,
			Metadata: &v3.Metadata{
				Id: id,
			},
			Status: &v3.Status{
				ConditionType:   "Describe",
				ConditionStatus: v3.ConditionStatus_StatusNotSet,
				Reason:          "Unable to fetch organization information",
				LastUpdated:     timestamppb.Now(),
			},
		}

		return organization, nil
	}

}

func (s *organizationService) GetByName(ctx context.Context, name string) (*systemv3.Organization, error) {

	organization := &systemv3.Organization{
		ApiVersion: apiVersion,
		Kind:       organizationKind,
		Metadata: &v3.Metadata{
			Name: name,
		},
	}
	entity, err := s.dao.GetByName(ctx, name, &models.Organization{})
	if err != nil {
		organization.Metadata = &v3.Metadata{
			Name: name,
		}
		organization.Status = &v3.Status{
			ConditionType:   "Describe",
			ConditionStatus: v3.ConditionStatus_StatusFailed,
			Reason:          err.Error(),
			LastUpdated:     timestamppb.Now(),
		}
		return organization, err
	}

	if org, ok := entity.(*models.Organization); ok {

		var settings systemv3.OrganizationSettings
		err := json.Unmarshal(org.Settings, &settings)
		if err != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Describe",
				ConditionStatus: v3.ConditionStatus_StatusFailed,
			}
			return organization, err
		}

		organization.Metadata = &v3.Metadata{
			Name:        org.Name,
			Description: org.Description,
			Id:          org.ID.String(),
			Partner:     org.PartnerId.String(),
			ModifiedAt:  timestamppb.New(org.ModifiedAt),
		}
		organization.Spec = &systemv3.OrganizationSpec{
			BillingAddress:    org.BillingAddress,
			Active:            org.Active,
			Approved:          org.Approved,
			Type:              org.Type,
			AddressLine1:      org.AddressLine1,
			AddressLine2:      org.AddressLine2,
			City:              org.City,
			Country:           org.Country,
			Phone:             org.Phone,
			State:             org.State,
			Zipcode:           org.Zipcode,
			IsPrivate:         org.IsPrivate,
			IsTotpEnabled:     org.IsTOTPEnabled,
			AreClustersShared: org.AreClustersShared,
			Settings:          &settings,
		}
		organization.Status = &v3.Status{
			ConditionType:   "Describe",
			LastUpdated:     timestamppb.Now(),
			ConditionStatus: v3.ConditionStatus_StatusOK,
		}
	}

	return organization, nil
}

func (s *organizationService) Update(ctx context.Context, organization *systemv3.Organization) (*systemv3.Organization, error) {

	entity, err := s.dao.GetByName(ctx, organization.Metadata.Name, &models.Organization{})
	if err != nil {
		organization.Status = &v3.Status{
			ConditionType:   "Update",
			ConditionStatus: v3.ConditionStatus_StatusFailed,
			Reason:          err.Error(),
			LastUpdated:     timestamppb.Now(),
		}
		return organization, err
	}

	if org, ok := entity.(*models.Organization); ok {

		sb, err := json.MarshalIndent(organization.GetSpec().GetSettings(), "", "\t")
		if err != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Update",
				ConditionStatus: v3.ConditionStatus_StatusFailed,
				Reason:          err.Error(),
				LastUpdated:     timestamppb.Now(),
			}
		}

		//update organization details
		org.Name = organization.GetMetadata().GetName()
		org.Description = organization.GetMetadata().GetDescription()
		org.ModifiedAt = time.Now()
		org.Trash = false
		org.Settings = json.RawMessage(sb)
		org.BillingAddress = organization.GetSpec().GetBillingAddress()
		org.Active = organization.GetSpec().GetActive()
		org.Approved = organization.GetSpec().GetApproved()
		org.Type = organization.GetSpec().GetType()
		org.AddressLine1 = organization.GetSpec().GetAddressLine1()
		org.AddressLine2 = organization.GetSpec().GetAddressLine2()
		org.City = organization.GetSpec().GetCity()
		org.Country = organization.GetSpec().GetCountry()
		org.Phone = organization.GetSpec().GetPhone()
		org.State = organization.GetSpec().GetState()
		org.Zipcode = organization.GetSpec().GetZipcode()
		org.IsPrivate = organization.GetSpec().GetIsPrivate()
		org.IsTOTPEnabled = organization.GetSpec().GetIsTotpEnabled()
		org.AreClustersShared = organization.GetSpec().GetAreClustersShared()

		_, err = s.dao.Update(ctx, org.ID, org)
		if err != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Update",
				ConditionStatus: v3.ConditionStatus_StatusFailed,
				LastUpdated:     timestamppb.Now(),
				Reason:          err.Error(),
			}
			return organization, err
		}

		//update status
		if organization.Status != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Update",
				ConditionStatus: v3.ConditionStatus_StatusOK,
				LastUpdated:     timestamppb.Now(),
			}
		}
	}

	return organization, nil
}

func (s *organizationService) Delete(ctx context.Context, organization *systemv3.Organization) (*systemv3.Organization, error) {

	entity, err := s.dao.GetByName(ctx, organization.Metadata.Name, &models.Organization{})
	if err != nil {
		organization.Status = &v3.Status{
			ConditionType:   "Delete",
			ConditionStatus: v3.ConditionStatus_StatusFailed,
			Reason:          err.Error(),
			LastUpdated:     timestamppb.Now(),
		}
		return organization, err
	}

	if org, ok := entity.(*models.Organization); ok {
		err = s.dao.Delete(ctx, org.ID, org)
		if err != nil {
			organization.Status = &v3.Status{
				ConditionType:   "Delete",
				ConditionStatus: v3.ConditionStatus_StatusFailed,
				Reason:          err.Error(),
				LastUpdated:     timestamppb.Now(),
			}
			return organization, err
		}
		//update v3 status
		organization.Metadata.Id = org.ID.String()
		organization.Metadata.Name = org.Name
		organization.Metadata.ModifiedAt = timestamppb.New(org.ModifiedAt)
		organization.Status = &v3.Status{
			ConditionType:   "Deleted",
			ConditionStatus: v3.ConditionStatus_StatusOK,
			LastUpdated:     timestamppb.Now(),
		}
	}
	return organization, nil

}

func (s *organizationService) List(ctx context.Context, organization *systemv3.Organization) (*systemv3.OrganizationList, error) {

	var organizations []*systemv3.Organization
	organinzationList := &systemv3.OrganizationList{
		ApiVersion: apiVersion,
		Kind:       organizationListKind,
		Metadata: &v3.ListMetadata{
			Count: 0,
		},
	}
	if len(organization.Metadata.Partner) > 0 {
		var partner models.Partner
		_, err := s.dao.GetByName(ctx, organization.Metadata.Partner, &partner)
		if err != nil {
			return organinzationList, err
		}

		var orgs []models.Organization
		entities, err := s.dao.List(ctx, uuid.NullUUID{UUID: partner.ID, Valid: true}, uuid.NullUUID{UUID: uuid.Nil}, &orgs)
		if err != nil {
			return organinzationList, err
		}
		if orgs, ok := entities.(*[]models.Organization); ok {
			for _, org := range *orgs {
				var settings systemv3.OrganizationSettings
				err := json.Unmarshal(org.Settings, &settings)
				if err != nil {
					fmt.Print(err)
					return nil, err
				}
				organization.Metadata = &v3.Metadata{
					Name:        org.Name,
					Description: org.Description,
					Id:          org.ID.String(),
					Partner:     org.PartnerId.String(),
					ModifiedAt:  timestamppb.New(org.ModifiedAt),
				}
				organization.Spec = &systemv3.OrganizationSpec{
					BillingAddress:    org.BillingAddress,
					Active:            org.Active,
					Approved:          org.Approved,
					Type:              org.Type,
					AddressLine1:      org.AddressLine1,
					AddressLine2:      org.AddressLine2,
					City:              org.City,
					Country:           org.Country,
					Phone:             org.Phone,
					State:             org.State,
					Zipcode:           org.Zipcode,
					IsPrivate:         org.IsPrivate,
					IsTotpEnabled:     org.IsTOTPEnabled,
					AreClustersShared: org.AreClustersShared,
					Settings:          &settings,
				}
				organizations = append(organizations, organization)
			}

			//update the list metadata and items response
			organinzationList.Metadata = &v3.ListMetadata{
				Count: int64(len(organizations)),
			}
			organinzationList.Items = organizations
		}

	} else {
		return organinzationList, fmt.Errorf("missing partner id in metadata")
	}
	return organinzationList, nil
}

func (s *organizationService) Close() error {
	return s.dao.Close()
}
