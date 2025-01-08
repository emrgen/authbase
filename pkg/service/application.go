package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

func NewApplicationService(store store.Provider) *ApplicationService {
	return &ApplicationService{store: store}
}

var _ v1.ApplicationServiceServer = (*ApplicationService)(nil)

type ApplicationService struct {
	store store.Provider
	v1.UnimplementedApplicationServiceServer
}

func (a *ApplicationService) CreateApplication(ctx context.Context, request *v1.CreateApplicationRequest) (*v1.CreateApplicationResponse, error) {
	var err error
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.PoolId)
	if err != nil {
		return nil, err
	}

	app := &model.Application{
		ID:     uuid.New().String(),
		Name:   request.Name,
		PoolID: poolID.String(),
	}

	err = as.CreateApplication(ctx, app)
	if err != nil {
		return nil, err
	}

	return &v1.CreateApplicationResponse{
		Application: &v1.Application{
			Id:     app.ID,
			Name:   app.Name,
			PoolId: app.PoolID,
		},
	}, nil
}

func (a *ApplicationService) GetApplication(ctx context.Context, request *v1.GetApplicationRequest) (*v1.GetApplicationResponse, error) {
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	app, err := as.GetApplication(ctx, appID)
	if err != nil {
		return nil, err
	}

	return &v1.GetApplicationResponse{
		Application: &v1.Application{
			Id:     app.ID,
			Name:   app.Name,
			PoolId: app.PoolID,
		},
	}, nil
}

func (a *ApplicationService) ListApplications(ctx context.Context, request *v1.ListApplicationsRequest) (*v1.ListApplicationsResponse, error) {
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.PoolId)
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)

	apps, total, err := as.ListApplications(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var applications []*v1.Application
	for _, app := range apps {
		applications = append(applications, &v1.Application{
			Id:     app.ID,
			Name:   app.Name,
			PoolId: app.PoolID,
		})
	}

	return &v1.ListApplicationsResponse{
		Applications: applications,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

func (a *ApplicationService) UpdateApplication(ctx context.Context, request *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		app, err := as.GetApplication(ctx, appID)
		if err != nil {
			return err
		}

		app.Name = request.Name

		err = as.UpdateApplication(ctx, app)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateApplicationResponse{}, nil
}

func (a *ApplicationService) DeleteApplication(ctx context.Context, request *v1.DeleteApplicationRequest) (*v1.DeleteApplicationResponse, error) {
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	err = as.DeleteApplication(ctx, appID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteApplicationResponse{}, nil
}
