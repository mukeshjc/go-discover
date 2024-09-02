// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package linodego

import (
	"context"
	"fmt"
	"time"
)

// Ticket represents a support ticket object
type Ticket struct {
	ID          int           `json:"id"`
	Attachments []string      `json:"attachments"`
	Closed      *time.Time    `json:"-"`
	Description string        `json:"description"`
	Entity      *TicketEntity `json:"entity"`
	GravatarID  string        `json:"gravatar_id"`
	Opened      *time.Time    `json:"-"`
	OpenedBy    string        `json:"opened_by"`
	Status      TicketStatus  `json:"status"`
	Summary     string        `json:"summary"`
	Updated     *time.Time    `json:"-"`
	UpdatedBy   string        `json:"updated_by"`
}

// TicketEntity refers a ticket to a specific entity
type TicketEntity struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// TicketStatus constants start with Ticket and include Linode API Ticket Status values
type TicketStatus string

// TicketStatus constants reflect the current status of a Ticket
const (
	TicketNew    TicketStatus = "new"
	TicketClosed TicketStatus = "closed"
	TicketOpen   TicketStatus = "open"
)

// TicketsPagedResponse represents a paginated ticket API response
type TicketsPagedResponse struct {
	*PageOptions
	Data []Ticket `json:"data"`
}

func (TicketsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Tickets.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *TicketsPagedResponse) appendData(r *TicketsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListTickets returns a collection of Support Tickets on the Account. Support Tickets
// can be both tickets opened with Linode for support, as well as tickets generated by
// Linode regarding the Account. This collection includes all Support Tickets generated
// on the Account, with open tickets returned first.
func (c *Client) ListTickets(ctx context.Context, opts *ListOptions) ([]Ticket, error) {
	response := TicketsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetTicket gets a Support Ticket on the Account with the specified ID
func (c *Client) GetTicket(ctx context.Context, id int) (*Ticket, error) {
	e, err := c.Tickets.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := coupleAPIErrors(c.R(ctx).
		SetResult(&Ticket{}).
		Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*Ticket), nil
}
