package service

import (
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProjectService_CreateProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	ab := createAdminProject(t)

	// create a project
	projectService := NewProjectService(ab.perm, ab.provider, ab.redis)
	pass := "password"
	_, err := projectService.CreateProject(ab.ctx, &v1.CreateProjectRequest{
		Name:        "test-project-1",
		VisibleName: "Test Project 1",
		Email:       "ab@gmail.com",
		Password:    &pass,
	})
	assert.NoError(t, err, "failed to create project")
}
