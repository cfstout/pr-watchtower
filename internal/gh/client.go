package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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
}

func FetchPRs(query string) ([]PR, error) {
	fmt.Printf("DEBUG: Executing gh search prs with query: %q\n", query)
	// gh search prs --json number,title,url,state,updatedAt,author --limit 30 <query parts...>
	args := []string{"search", "prs", "--json", "number,title,url,state,updatedAt,author", "--limit", "30"}
	// Split query by spaces to pass as separate arguments
	// This is a simple split; for complex queries with quoted spaces, we'd need a proper parser.
	// But for our simple config, this should suffice.
	// Actually, gh search prs takes the query as separate arguments.
	// If we pass "review-requested:@me state:open" as one arg, gh might be trying to parse it as one token.
	// Let's try splitting it.
	queryParts := strings.Fields(query)
	args = append(args, queryParts...)

	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to run gh command: %s", string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to run gh command: %w", err)
	}

	var prs []PR
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse gh output: %w", err)
	}

	return prs, nil
}

func TriggerWorkflow(prNumber int) error {
	// gh workflow run agent-fix.yml -f pr_number=<prNumber>
	cmd := exec.Command("gh", "workflow", "run", "agent-fix.yml", "-f", fmt.Sprintf("pr_number=%d", prNumber))
	return cmd.Run()
}
