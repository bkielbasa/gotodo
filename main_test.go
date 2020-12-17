package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	todo "github.com/bkielbasa/gotodo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type projectDetail struct {
	ID       string
	Name     string
	Archived bool
}

type projectsListResponse struct {
	Projects []projectDetail
}

const port = 8090

func TestAddingNewProject(t *testing.T) {
	ctx := context.Background()
	run, shutdown := todo.App(ctx, port)
	defer shutdown()

	go run()
	name := uuid.New().String()
	reqBody := fmt.Sprintf("{\"name\": \"%s\"}", name)
	url := "http://localhost:8090/project/create"
	client := http.Client{
		Timeout: time.Second,
	}

	// create a new project
	resp, err := client.Post(url, "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	r := struct{ ID string }{}
	json.Unmarshal(respBody, &r)
	require.NotEmpty(t, r.ID)

	// list existing projects
	url = "http://localhost:8090/projects"
	resp, err = client.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listResp := projectsListResponse{}
	respBody, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	json.Unmarshal(respBody, &listResp)

	// check if the projects is on the list
	found := false

	for _, proj := range listResp.Projects {
		if proj.ID == r.ID {
			require.Equal(t, name, proj.Name)
			require.False(t, proj.Archived)
			found = true
		}
	}

	require.True(t, found)
}
