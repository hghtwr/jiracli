package changeAssignee

import (
	"reflect"
	testing "testing"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/huh"
	"github.com/hghtwr/jiracli/jiraApi"
	"github.com/hghtwr/jiracli/navigation"
)

func TestCreateInitModel(t *testing.T) {
	model := CreateInitModel()
	if model.NavTo != navigation.ChangeAssigneeView {
		t.Errorf("Expected %v, got %v", navigation.ChangeAssigneeView, model.NavTo)
	}
	if model.form == nil {
		t.Errorf("Expected form to be initialized, got nil")
	}
	if model.form.State  != huh.StateNormal{
		t.Errorf("Expected form to be in state 'normal', got %v", model.form.State)
	}
}

func TestSetContext(t *testing.T) {
	model := CreateInitModel()
	context := navigation.Context{
		IssueId: "TEST-1",
	}
	m := model.SetContext(context)
	if m.GetContext().IssueId != "TEST-1" {
		t.Errorf("Expected %v, got %v", "TEST-1", model.Context.IssueId)
	}
}

func TestGetContext(t *testing.T) {
	model := CreateInitModel()
	context := navigation.Context{
		IssueId: "TEST-1",
	}
	model.Context = context
	if model.GetContext().IssueId != "TEST-1" {
		t.Errorf("Expected %v, got %v", "TEST-1", model.Context.IssueId)
	}
}


func TestUpdateSearchAssignees(t *testing.T) {
	model := CreateInitModel()
	users := []jira.User{
		{
			DisplayName: "Test User",
			AccountID: "1234",
			EmailAddress: "test-user@test.com",
		},
		{
			DisplayName: "Test User 2",
			AccountID: "12345",
			EmailAddress: "test-user2@test.com",
		},
	}

	assigneeResponse := jiraApi.SearchAssigneesResponse{
		Assignees: users,
	}
	m, _ := model.Update(assigneeResponse)
	assigneeModel, ok := m.(ChangeAssigneeModel)
	if !ok {
		t.Errorf("Expected ChangeAssigneeModel, got %v", reflect.TypeOf(m))
	}
	if !reflect.DeepEqual(users, assigneeModel.assigneeSearchResult) {
		t.Errorf("Expected %v, got %v", users, assigneeModel.assigneeSearchResult)
	}

}