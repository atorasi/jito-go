package jito

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type validatorAPIResponse struct {
	Validators []*jitoValidator `json:"validators"`
}

type jitoValidator struct {
	VoteAccount string `json:"vote_account"`
	RunningJito bool   `json:"running_jito"`
}

func (j *JitoManager) fetchJitoValidators() error {
	j.status("Fetching jito-enabled validators")

	req, err := http.NewRequest("GET", "https://kobe.mainnet.jito.network/api/v1/validators", nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")

	resp, err := j.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch validators: %s", resp.Status)
	}

	var validators validatorAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&validators)
	if err != nil {
		return err
	}

	j.buildJitoValidators(validators.Validators)

	return nil
}

func (j *JitoManager) buildJitoValidators(validators []*jitoValidator) {
	j.Lock.Lock()
	defer j.Lock.Unlock()
	j.JitoValidators = make(map[string]bool)

	for i := range validators {
		if validators[i].RunningJito {
			j.JitoValidators[validators[i].VoteAccount] = true
		}
	}
}
