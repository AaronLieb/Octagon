// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package main

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// __generateRegistrationTokenInput is used internally by genqlient
type __generateRegistrationTokenInput struct {
	EventId int `json:"eventId"`
	UserId  int `json:"userId"`
}

// GetEventId returns __generateRegistrationTokenInput.EventId, and is useful for accessing the field via an interface.
func (v *__generateRegistrationTokenInput) GetEventId() int { return v.EventId }

// GetUserId returns __generateRegistrationTokenInput.UserId, and is useful for accessing the field via an interface.
func (v *__generateRegistrationTokenInput) GetUserId() int { return v.UserId }

// __getEventInput is used internally by genqlient
type __getEventInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __getEventInput.Slug, and is useful for accessing the field via an interface.
func (v *__getEventInput) GetSlug() string { return v.Slug }

// __getTournamentInput is used internally by genqlient
type __getTournamentInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __getTournamentInput.Slug, and is useful for accessing the field via an interface.
func (v *__getTournamentInput) GetSlug() string { return v.Slug }

// generateRegistrationTokenResponse is returned by generateRegistrationToken on success.
type generateRegistrationTokenResponse struct {
	GenerateRegistrationToken string `json:"generateRegistrationToken"`
}

// GetGenerateRegistrationToken returns generateRegistrationTokenResponse.GenerateRegistrationToken, and is useful for accessing the field via an interface.
func (v *generateRegistrationTokenResponse) GetGenerateRegistrationToken() string {
	return v.GenerateRegistrationToken
}

// getCurrentUserCurrentUser includes the requested fields of the GraphQL type User.
type getCurrentUserCurrentUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// GetId returns getCurrentUserCurrentUser.Id, and is useful for accessing the field via an interface.
func (v *getCurrentUserCurrentUser) GetId() int { return v.Id }

// GetName returns getCurrentUserCurrentUser.Name, and is useful for accessing the field via an interface.
func (v *getCurrentUserCurrentUser) GetName() string { return v.Name }

// getCurrentUserResponse is returned by getCurrentUser on success.
type getCurrentUserResponse struct {
	CurrentUser getCurrentUserCurrentUser `json:"currentUser"`
}

// GetCurrentUser returns getCurrentUserResponse.CurrentUser, and is useful for accessing the field via an interface.
func (v *getCurrentUserResponse) GetCurrentUser() getCurrentUserCurrentUser { return v.CurrentUser }

// getEventEvent includes the requested fields of the GraphQL type Event.
type getEventEvent struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// GetId returns getEventEvent.Id, and is useful for accessing the field via an interface.
func (v *getEventEvent) GetId() int { return v.Id }

// GetName returns getEventEvent.Name, and is useful for accessing the field via an interface.
func (v *getEventEvent) GetName() string { return v.Name }

// getEventResponse is returned by getEvent on success.
type getEventResponse struct {
	Event getEventEvent `json:"event"`
}

// GetEvent returns getEventResponse.Event, and is useful for accessing the field via an interface.
func (v *getEventResponse) GetEvent() getEventEvent { return v.Event }

// getTournamentResponse is returned by getTournament on success.
type getTournamentResponse struct {
	Tournament getTournamentTournament `json:"tournament"`
}

// GetTournament returns getTournamentResponse.Tournament, and is useful for accessing the field via an interface.
func (v *getTournamentResponse) GetTournament() getTournamentTournament { return v.Tournament }

// getTournamentTournament includes the requested fields of the GraphQL type Tournament.
type getTournamentTournament struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// GetId returns getTournamentTournament.Id, and is useful for accessing the field via an interface.
func (v *getTournamentTournament) GetId() int { return v.Id }

// GetName returns getTournamentTournament.Name, and is useful for accessing the field via an interface.
func (v *getTournamentTournament) GetName() string { return v.Name }

// The mutation executed by generateRegistrationToken.
const generateRegistrationToken_Operation = `
mutation generateRegistrationToken ($eventId: ID!, $userId: ID!) {
	generateRegistrationToken(registration: {eventIds:[$eventId]}, userId: $userId)
}
`

func generateRegistrationToken(
	ctx_ context.Context,
	client_ graphql.Client,
	eventId int,
	userId int,
) (data_ *generateRegistrationTokenResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "generateRegistrationToken",
		Query:  generateRegistrationToken_Operation,
		Variables: &__generateRegistrationTokenInput{
			EventId: eventId,
			UserId:  userId,
		},
	}

	data_ = &generateRegistrationTokenResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by getCurrentUser.
const getCurrentUser_Operation = `
query getCurrentUser {
	currentUser {
		id
		name
	}
}
`

func getCurrentUser(
	ctx_ context.Context,
	client_ graphql.Client,
) (data_ *getCurrentUserResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "getCurrentUser",
		Query:  getCurrentUser_Operation,
	}

	data_ = &getCurrentUserResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by getEvent.
const getEvent_Operation = `
query getEvent ($slug: String) {
	event(slug: $slug) {
		id
		name
	}
}
`

func getEvent(
	ctx_ context.Context,
	client_ graphql.Client,
	slug string,
) (data_ *getEventResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "getEvent",
		Query:  getEvent_Operation,
		Variables: &__getEventInput{
			Slug: slug,
		},
	}

	data_ = &getEventResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by getTournament.
const getTournament_Operation = `
query getTournament ($slug: String) {
	tournament(slug: $slug) {
		id
		name
	}
}
`

func getTournament(
	ctx_ context.Context,
	client_ graphql.Client,
	slug string,
) (data_ *getTournamentResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "getTournament",
		Query:  getTournament_Operation,
		Variables: &__getTournamentInput{
			Slug: slug,
		},
	}

	data_ = &getTournamentResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}
