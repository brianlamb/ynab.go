// Copyright (c) 2018, Bruno M V Souza <github@b.bmvs.io>. All rights reserved.
// Use of this source code is governed by a BSD-2-Clause license that can be
// found in the LICENSE file.

package account

import (
	"encoding/json"
	"fmt"

	"go.bmvs.io/ynab/api"
)

// NewService facilitates the creation of a new account service instance
func NewService(c api.ClientReaderWriter) *Service {
	return &Service{c}
}

// Service wraps YNAB account API endpoints
type Service struct {
	c api.ClientReaderWriter
}

// GetAccounts fetches the list of accounts from a budget
// https://api.youneedabudget.com/v1#/Accounts/getAccounts
func (s *Service) GetAccounts(budgetID string, f *api.Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Accounts        []*Account `json:"accounts"`
			ServerKnowledge uint64     `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/accounts", budgetID)
	if f != nil {
		url = fmt.Sprintf("%s?%s", url, f.ToQuery())
	}
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		Accounts:        resModel.Data.Accounts,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetAccount fetches a specific account from a budget
// https://api.youneedabudget.com/v1#/Accounts/getAccountById
func (s *Service) GetAccount(budgetID, accountID string) (*Account, error) {
	resModel := struct {
		Data struct {
			Account *Account `json:"account"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/accounts/%s", budgetID, accountID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Account, nil
}

type PayloadAccount struct {
	Name string `json:"name"`
	Type Type `json:"type"`
	Balance int64 `json:"balance"`
}


// CreateAccount creates an account for a budget
// https://api.youneedabudget.com/v1#/Accounts/createAccount
func (s *Service) CreateAccount(budgetID string, accountPayload PayloadAccount) (*Account, error) {
	resModel := struct {
		Data struct {
			Account *Account `json:"account"`
		} `json:"data"`
	}{}

	payload := struct {
		Account PayloadAccount `json:"account"`
	}{
		accountPayload,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/budgets/%s/accounts/", budgetID)
	if err := s.c.POST(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Account, nil
}
