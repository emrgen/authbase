package service

import (
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testEmail1 = "ab1@gmail.com"
	testEmail2 = "ab2@gmail.com"
)

func createProject(t *testing.T, name, visibleName, email, pass string) (*authbase, *ProjectService, *v1.Project) {
	ab := createAdminProject(t)

	// create a project
	projectService := NewProjectService(ab.perm, ab.provider, ab.redis)
	res, err := projectService.CreateProject(ab.ctx, &v1.CreateProjectRequest{
		Name:        name,
		VisibleName: visibleName,
		Email:       email,
		Password:    &pass,
	})
	assert.NoError(t, err, "failed to create project")

	return ab, projectService, res.Project
}

func TestProjectService_CreateProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	createProject(t, "test-project-1", "Test Project 1", testEmail1, "password")
}

func TestProjectService_GetProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	ab := createAdminProject(t)

	// create a project
	projectService := NewProjectService(ab.perm, ab.provider, ab.redis)
	pass := "password"
	res, err := projectService.CreateProject(ab.ctx, &v1.CreateProjectRequest{
		Name:        "test-project-2",
		VisibleName: "Test Project 2",
		Email:       testEmail2,
		Password:    &pass,
	})
	assert.NoError(t, err, "failed to create project")

	// get the project
	projectRes, err := projectService.GetProject(ab.ctx, &v1.GetProjectRequest{
		Id: res.Project.Id,
	})
	assert.NoError(t, err, "failed to get project")

	assert.Equal(t, projectRes.Project.Name, "test-project-2", "project name is not correct")
}

func TestProjectService_ListProjects(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	ab := createAdminProject(t)

	// create a project
	projectService := NewProjectService(ab.perm, ab.provider, ab.redis)
	pass := "password"
	_, err := projectService.CreateProject(ab.ctx, &v1.CreateProjectRequest{
		Name:        "test-project-3",
		VisibleName: "Test Project 3",
		Email:       testEmail1,
		Password:    &pass,
	})
	assert.NoError(t, err, "failed to create project")

	// list projects
	projectRes, err := projectService.ListProjects(ab.ctx, &v1.ListProjectsRequest{})
	assert.NoError(t, err, "failed to list projects")

	assert.Equal(t, len(projectRes.Projects), 2, "project count is not correct")
}

func TestProjectService_UpdateProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	ab, projectService, project := createProject(t, "test-project-1", "Test Project 1", testEmail1, "password")

	// update the project
	_, err := projectService.UpdateProject(ab.ctx, &v1.UpdateProjectRequest{
		Id:   project.Id,
		Name: "test-project-4",
	})
	assert.NoError(t, err, "failed to update project")

	// get the project
	projectRes, err := projectService.GetProject(ab.ctx, &v1.GetProjectRequest{
		Id: project.Id,
	})
	assert.NoError(t, err, "failed to get project")

	assert.Equal(t, projectRes.Project.Name, "test-project-4", "project name is not correct")
}

func TestProjectService_DeleteProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	ab, projectService, project := createProject(t, "test-project-1", "Test Project 1", testEmail1, "password")

	// delete the project
	_, err := projectService.DeleteProject(ab.ctx, &v1.DeleteProjectRequest{
		Id: project.Id,
	})
	assert.NoError(t, err, "failed to delete project: %v", err)

	// get the project
	_, err = projectService.GetProject(ab.ctx, &v1.GetProjectRequest{
		Id: project.Id,
	})
	assert.Error(t, err, "project should not exist")
}
