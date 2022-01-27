package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/blang/semver/v4"
)

type githubReleaseResponse struct {
	Name string `json:"name"`
}

func GetGithubReleaseURL(releasesURL string, v *semver.Version) string {
	return fmt.Sprintf("%v/tag/v%v", releasesURL, v)
}

func BuildGithubReleasesRequestFrom(ctx context.Context, releasesURL string) ReleasesGetter {
	return func() ([]string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
		if err != nil {
			return nil, fmt.Errorf("couldn't build request: %w", err)
		}
		req.Header.Add("Accept", "application/vnd.github.v3+json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("couldn't deliver request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read response body: %w", err)
		}

		responses := []githubReleaseResponse{}
		if err = json.Unmarshal(body, &responses); err != nil {
			// try to parse as a general error message which would be useful information
			// to know eg. if we were blocked due to GitHub rate-limiting
			m := struct {
				Message string `json:"message"`
			}{}
			if mErr := json.Unmarshal(body, &m); mErr == nil {
				return nil, fmt.Errorf("couldn't read response message: %s: %w", m.Message, err)
			}

			return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
		}

		releases := make([]string, 0, len(responses))
		for _, response := range responses {
			releases = append(releases, response.Name)
		}

		return releases, nil
	}
}
