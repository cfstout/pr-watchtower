package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type PR struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
	Mergeable         string `json:"mergeable"`
	StatusCheckRollup struct {
		State string `json:"state"`
	} `json:"statusCheckRollup"`
}

type graphQLResponse struct {
	Data struct {
		Search struct {
			Nodes []PR `json:"nodes"`
		} `json:"search"`
	} `json:"data"`
}

func FetchPRs(query string) ([]PR, error) {
	// Construct GraphQL query
	gqlQuery := `
		query($q: String!) {
			search(query: $q, type: ISSUE, first: 30) {
				nodes {
					... on PullRequest {
						number
						title
						url
						state
						updatedAt
						author {
							login
						}
						mergeable
						statusCheckRollup {
							state
						}
					}
				}
			}
		}
	`

	// gh api graphql -f q='<query>' -f query='<gqlQuery>'
	// We need to pass the search query as a variable.
	// Note: The user's config query might contain spaces and special chars.

	args := []string{"api", "graphql", "-f", fmt.Sprintf("q=%s", query), "-f", fmt.Sprintf("query=%s", gqlQuery)}

	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to run gh command: %s", string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to run gh command: %w", err)
	}

	var resp graphQLResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse gh output: %w", err)
	}

	var prs []PR
	for _, node := range resp.Data.Search.Nodes {
		if node.Number != 0 {
			prs = append(prs, node)
		}
	}

	return prs, nil
}

func TriggerWorkflow(prNumber int) error {
	// gh workflow run agent-fix.yml -f pr_number=<prNumber>
	cmd := exec.Command("gh", "workflow", "run", "agent-fix.yml", "-f", fmt.Sprintf("pr_number=%d", prNumber))
	return cmd.Run()
}

func MergePR(prNumber int) error {
	// gh pr merge <number> --merge --auto
	cmd := exec.Command("gh", "pr", "merge", fmt.Sprintf("%d", prNumber), "--merge", "--auto")
	return cmd.Run()
}
