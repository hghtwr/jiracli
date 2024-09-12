package jiraApi

import (
	"fmt"
	"net/url"
	"os"

	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hghtwr/jiracli/notifications"
)

type MyIssuesResponse struct {
	Issues []jira.Issue
}

type IssueDetailResponse struct {
	Issue *jira.Issue
}

type IssueParentResponse struct {
	Issue *jira.Issue
}

type ChildIssuesResponse struct {
	Issues []jira.Issue
}

type SearchAssigneesResponse struct {
	Assignees []jira.User
}

const (
	NoResponse = "No response returned"
)

var client, _ = CreateClient() // TODO: handle error

func FetchChildIssues(issueId string) tea.Cmd {
	return func() tea.Msg {
		issues, err := fetchIssuesWithJql(fmt.Sprintf("parent = %s", issueId))

		if err != nil {
			return notifications.NotificationCmd{
				Message: fmt.Sprintf("%s - %s", "Error fetching child Issues", err.Error()),
				Type: notifications.Error,
				Mode: notifications.Tray,
			}
		}

		return ChildIssuesResponse{
			Issues: issues,
		}
	}
}

func FetchAssignedIssues() tea.Cmd {
	return func() tea.Msg {
		issues, err := fetchIssuesWithJql("assignee = currentUser() AND resolution = Unresolved ORDER BY priority DESC, updated DESC")

		if err != nil {
			return notifications.NotificationCmd{
				Message: fmt.Sprintf("%s - %s", "Error fetching assigned Issues", err.Error()),
				Type: notifications.Error,
				Mode: notifications.Tray,
			}
		}

		return MyIssuesResponse{
			Issues: issues,
		}
	}
}
func FetchParentIssueDetails(issueId string) tea.Cmd {
	return func() tea.Msg {

		issue, _, err := client.Issue.Get(issueId, nil)

		if err != nil {
			return notifications.NotificationCmd{
				Message: fmt.Sprintf("%s - %s", "Error fetching issue details", err.Error()),
				Type: notifications.Error,
				Mode: notifications.Tray,
			}
		}

		return IssueParentResponse{
			Issue: issue,
		}
	}
}

func FetchIssueDetails(issueId string) tea.Cmd {
	return func() tea.Msg {

		issue, _, err := client.Issue.Get(issueId, nil)

		if err != nil {
			return notifications.NotificationCmd{
				Message: fmt.Sprintf("%s - %s", "Error fetching issue details", err.Error()),
				Type: notifications.Error,
				Mode: notifications.Tray,
			}
		}

		return IssueDetailResponse{
			Issue: issue,
		}
	}
}



func fetchIssuesWithJql(query string) ([]jira.Issue, error) {
	var err error

	var fetchedIssues []jira.Issue

	storeIssues := func(issue jira.Issue) error {
		fetchedIssues = append(fetchedIssues, issue)
		return nil
	}

	err = client.Issue.SearchPages(query,
		&jira.SearchOptions{MaxResults: 100}, storeIssues)

	if err != nil {
		return nil, err
	}
	return fetchedIssues, nil
}

func CreateClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("API_TOKEN"),
	}

	jiraClient, clientErr := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))

	if clientErr != nil {
		return nil, clientErr
	}
	return jiraClient, nil
}


func SearchAssignees(username string) tea.Cmd {

	return func() tea.Msg {

		encodedUsername := url.QueryEscape(username)
		users, _, respError := client.User.Find(encodedUsername)

		if respError != nil {
			if respError.Error() != NoResponse {

				return notifications.NotificationCmd{
					Message: fmt.Sprintf("%s - %s", "Error fetching issue details", respError.Error()),
					Type: notifications.Error,
					Mode: notifications.Tray,
				}
			}
		}
		return SearchAssigneesResponse{
			Assignees: users,
		}
	}
}
