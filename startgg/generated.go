// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package startgg

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// Game specific H2H set data such as character, stage, and stock info
type BracketSetGameDataInput struct {
	// Entrant ID of game winner
	WinnerId int `json:"winnerId"`
	// Game number
	GameNum int `json:"gameNum"`
	// Score for entrant 1 (if applicable). For smash, this is stocks remaining.
	Entrant1Score int `json:"entrant1Score"`
	// Score for entrant 2 (if applicable). For smash, this is stocks remaining.
	Entrant2Score int `json:"entrant2Score"`
	// ID of the stage that was played for this game (if applicable)
	StageId int `json:"stageId"`
	// List of selections for the game, typically character selections.
	Selections []BracketSetGameSelectionInput `json:"selections"`
}

// GetWinnerId returns BracketSetGameDataInput.WinnerId, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetWinnerId() int { return v.WinnerId }

// GetGameNum returns BracketSetGameDataInput.GameNum, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetGameNum() int { return v.GameNum }

// GetEntrant1Score returns BracketSetGameDataInput.Entrant1Score, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetEntrant1Score() int { return v.Entrant1Score }

// GetEntrant2Score returns BracketSetGameDataInput.Entrant2Score, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetEntrant2Score() int { return v.Entrant2Score }

// GetStageId returns BracketSetGameDataInput.StageId, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetStageId() int { return v.StageId }

// GetSelections returns BracketSetGameDataInput.Selections, and is useful for accessing the field via an interface.
func (v *BracketSetGameDataInput) GetSelections() []BracketSetGameSelectionInput { return v.Selections }

// Game specific H2H selections made by the entrants, such as character info
type BracketSetGameSelectionInput struct {
	// Entrant ID that made selection
	EntrantId int `json:"entrantId"`
	// Character selected by this entrant for this game.
	CharacterId int `json:"characterId"`
}

// GetEntrantId returns BracketSetGameSelectionInput.EntrantId, and is useful for accessing the field via an interface.
func (v *BracketSetGameSelectionInput) GetEntrantId() int { return v.EntrantId }

// GetCharacterId returns BracketSetGameSelectionInput.CharacterId, and is useful for accessing the field via an interface.
func (v *BracketSetGameSelectionInput) GetCharacterId() int { return v.CharacterId }

// GenerateRegistrationTokenForEventResponse is returned by GenerateRegistrationTokenForEvent on success.
type GenerateRegistrationTokenForEventResponse struct {
	// Generate tournament registration Token on behalf of user
	GenerateRegistrationToken string `json:"generateRegistrationToken"`
}

// GetGenerateRegistrationToken returns GenerateRegistrationTokenForEventResponse.GenerateRegistrationToken, and is useful for accessing the field via an interface.
func (v *GenerateRegistrationTokenForEventResponse) GetGenerateRegistrationToken() string {
	return v.GenerateRegistrationToken
}

// GenerateRegistrationTokenResponse is returned by GenerateRegistrationToken on success.
type GenerateRegistrationTokenResponse struct {
	// Generate tournament registration Token on behalf of user
	GenerateRegistrationToken string `json:"generateRegistrationToken"`
}

// GetGenerateRegistrationToken returns GenerateRegistrationTokenResponse.GenerateRegistrationToken, and is useful for accessing the field via an interface.
func (v *GenerateRegistrationTokenResponse) GetGenerateRegistrationToken() string {
	return v.GenerateRegistrationToken
}

// GetCurrentUserCurrentUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user
type GetCurrentUserCurrentUser struct {
	Id int `json:"id"`
	// Public facing user name that respects user publishing settings
	Name string `json:"name"`
}

// GetId returns GetCurrentUserCurrentUser.Id, and is useful for accessing the field via an interface.
func (v *GetCurrentUserCurrentUser) GetId() int { return v.Id }

// GetName returns GetCurrentUserCurrentUser.Name, and is useful for accessing the field via an interface.
func (v *GetCurrentUserCurrentUser) GetName() string { return v.Name }

// GetCurrentUserResponse is returned by GetCurrentUser on success.
type GetCurrentUserResponse struct {
	// Returns the authenticated user
	CurrentUser GetCurrentUserCurrentUser `json:"currentUser"`
}

// GetCurrentUser returns GetCurrentUserResponse.CurrentUser, and is useful for accessing the field via an interface.
func (v *GetCurrentUserResponse) GetCurrentUser() GetCurrentUserCurrentUser { return v.CurrentUser }

// GetEntrantByNameEvent includes the requested fields of the GraphQL type Event.
// The GraphQL type's documentation follows.
//
// An event in a tournament
type GetEntrantByNameEvent struct {
	// The entrants that belong to an event, paginated by filter criteria
	Entrants GetEntrantByNameEventEntrantsEntrantConnection `json:"entrants"`
}

// GetEntrants returns GetEntrantByNameEvent.Entrants, and is useful for accessing the field via an interface.
func (v *GetEntrantByNameEvent) GetEntrants() GetEntrantByNameEventEntrantsEntrantConnection {
	return v.Entrants
}

// GetEntrantByNameEventEntrantsEntrantConnection includes the requested fields of the GraphQL type EntrantConnection.
type GetEntrantByNameEventEntrantsEntrantConnection struct {
	Nodes []GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant `json:"nodes"`
}

// GetNodes returns GetEntrantByNameEventEntrantsEntrantConnection.Nodes, and is useful for accessing the field via an interface.
func (v *GetEntrantByNameEventEntrantsEntrantConnection) GetNodes() []GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant {
	return v.Nodes
}

// GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant includes the requested fields of the GraphQL type Entrant.
// The GraphQL type's documentation follows.
//
// An entrant in an event
type GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant struct {
	Id int `json:"id"`
	// The entrant name as it appears in bracket: gamerTag of the participant or team name
	Name string `json:"name"`
}

// GetId returns GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant.Id, and is useful for accessing the field via an interface.
func (v *GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant) GetId() int { return v.Id }

// GetName returns GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant.Name, and is useful for accessing the field via an interface.
func (v *GetEntrantByNameEventEntrantsEntrantConnectionNodesEntrant) GetName() string { return v.Name }

// GetEntrantByNameResponse is returned by GetEntrantByName on success.
type GetEntrantByNameResponse struct {
	// Returns an event given its id or slug
	Event GetEntrantByNameEvent `json:"event"`
}

// GetEvent returns GetEntrantByNameResponse.Event, and is useful for accessing the field via an interface.
func (v *GetEntrantByNameResponse) GetEvent() GetEntrantByNameEvent { return v.Event }

// GetEntrantsEvent includes the requested fields of the GraphQL type Event.
// The GraphQL type's documentation follows.
//
// An event in a tournament
type GetEntrantsEvent struct {
	// The entrants that belong to an event, paginated by filter criteria
	Entrants GetEntrantsEventEntrantsEntrantConnection `json:"entrants"`
}

// GetEntrants returns GetEntrantsEvent.Entrants, and is useful for accessing the field via an interface.
func (v *GetEntrantsEvent) GetEntrants() GetEntrantsEventEntrantsEntrantConnection { return v.Entrants }

// GetEntrantsEventEntrantsEntrantConnection includes the requested fields of the GraphQL type EntrantConnection.
type GetEntrantsEventEntrantsEntrantConnection struct {
	Nodes []GetEntrantsEventEntrantsEntrantConnectionNodesEntrant `json:"nodes"`
}

// GetNodes returns GetEntrantsEventEntrantsEntrantConnection.Nodes, and is useful for accessing the field via an interface.
func (v *GetEntrantsEventEntrantsEntrantConnection) GetNodes() []GetEntrantsEventEntrantsEntrantConnectionNodesEntrant {
	return v.Nodes
}

// GetEntrantsEventEntrantsEntrantConnectionNodesEntrant includes the requested fields of the GraphQL type Entrant.
// The GraphQL type's documentation follows.
//
// An entrant in an event
type GetEntrantsEventEntrantsEntrantConnectionNodesEntrant struct {
	// The entrant name as it appears in bracket: gamerTag of the participant or team name
	Name string `json:"name"`
}

// GetName returns GetEntrantsEventEntrantsEntrantConnectionNodesEntrant.Name, and is useful for accessing the field via an interface.
func (v *GetEntrantsEventEntrantsEntrantConnectionNodesEntrant) GetName() string { return v.Name }

// GetEntrantsOutEvent includes the requested fields of the GraphQL type Event.
// The GraphQL type's documentation follows.
//
// An event in a tournament
type GetEntrantsOutEvent struct {
	// The entrants that belong to an event, paginated by filter criteria
	Entrants GetEntrantsOutEventEntrantsEntrantConnection `json:"entrants"`
}

// GetEntrants returns GetEntrantsOutEvent.Entrants, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEvent) GetEntrants() GetEntrantsOutEventEntrantsEntrantConnection {
	return v.Entrants
}

// GetEntrantsOutEventEntrantsEntrantConnection includes the requested fields of the GraphQL type EntrantConnection.
type GetEntrantsOutEventEntrantsEntrantConnection struct {
	Nodes []GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant `json:"nodes"`
}

// GetNodes returns GetEntrantsOutEventEntrantsEntrantConnection.Nodes, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEventEntrantsEntrantConnection) GetNodes() []GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant {
	return v.Nodes
}

// GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant includes the requested fields of the GraphQL type Entrant.
// The GraphQL type's documentation follows.
//
// An entrant in an event
type GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant struct {
	// The entrant name as it appears in bracket: gamerTag of the participant or team name
	Name string `json:"name"`
	// Standing for this entrant given an event. All entrants queried must be in the same event (for now).
	Standing GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding `json:"standing"`
}

// GetName returns GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant.Name, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant) GetName() string { return v.Name }

// GetStanding returns GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant.Standing, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrant) GetStanding() GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding {
	return v.Standing
}

// GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding includes the requested fields of the GraphQL type Standing.
// The GraphQL type's documentation follows.
//
// A standing indicates the placement of something within a container.
type GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding struct {
	IsFinal   bool `json:"isFinal"`
	Placement int  `json:"placement"`
}

// GetIsFinal returns GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding.IsFinal, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding) GetIsFinal() bool {
	return v.IsFinal
}

// GetPlacement returns GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding.Placement, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutEventEntrantsEntrantConnectionNodesEntrantStanding) GetPlacement() int {
	return v.Placement
}

// GetEntrantsOutResponse is returned by GetEntrantsOut on success.
type GetEntrantsOutResponse struct {
	// Returns an event given its id or slug
	Event GetEntrantsOutEvent `json:"event"`
}

// GetEvent returns GetEntrantsOutResponse.Event, and is useful for accessing the field via an interface.
func (v *GetEntrantsOutResponse) GetEvent() GetEntrantsOutEvent { return v.Event }

// GetEntrantsResponse is returned by GetEntrants on success.
type GetEntrantsResponse struct {
	// Returns an event given its id or slug
	Event GetEntrantsEvent `json:"event"`
}

// GetEvent returns GetEntrantsResponse.Event, and is useful for accessing the field via an interface.
func (v *GetEntrantsResponse) GetEvent() GetEntrantsEvent { return v.Event }

// GetEventEvent includes the requested fields of the GraphQL type Event.
// The GraphQL type's documentation follows.
//
// An event in a tournament
type GetEventEvent struct {
	Id int `json:"id"`
	// Title of event set by organizer
	Name string `json:"name"`
}

// GetId returns GetEventEvent.Id, and is useful for accessing the field via an interface.
func (v *GetEventEvent) GetId() int { return v.Id }

// GetName returns GetEventEvent.Name, and is useful for accessing the field via an interface.
func (v *GetEventEvent) GetName() string { return v.Name }

// GetEventResponse is returned by GetEvent on success.
type GetEventResponse struct {
	// Returns an event given its id or slug
	Event GetEventEvent `json:"event"`
}

// GetEvent returns GetEventResponse.Event, and is useful for accessing the field via an interface.
func (v *GetEventResponse) GetEvent() GetEventEvent { return v.Event }

// GetParticipantsResponse is returned by GetParticipants on success.
type GetParticipantsResponse struct {
	// Returns a tournament given its id or slug
	Tournament GetParticipantsTournament `json:"tournament"`
}

// GetTournament returns GetParticipantsResponse.Tournament, and is useful for accessing the field via an interface.
func (v *GetParticipantsResponse) GetTournament() GetParticipantsTournament { return v.Tournament }

// GetParticipantsTournament includes the requested fields of the GraphQL type Tournament.
// The GraphQL type's documentation follows.
//
// A tournament
type GetParticipantsTournament struct {
	Id int `json:"id"`
	// The tournament name
	Name string `json:"name"`
	// Paginated, queryable list of participants
	Participants GetParticipantsTournamentParticipantsParticipantConnection `json:"participants"`
}

// GetId returns GetParticipantsTournament.Id, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournament) GetId() int { return v.Id }

// GetName returns GetParticipantsTournament.Name, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournament) GetName() string { return v.Name }

// GetParticipants returns GetParticipantsTournament.Participants, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournament) GetParticipants() GetParticipantsTournamentParticipantsParticipantConnection {
	return v.Participants
}

// GetParticipantsTournamentParticipantsParticipantConnection includes the requested fields of the GraphQL type ParticipantConnection.
type GetParticipantsTournamentParticipantsParticipantConnection struct {
	PageInfo GetParticipantsTournamentParticipantsParticipantConnectionPageInfo           `json:"pageInfo"`
	Nodes    []GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant `json:"nodes"`
}

// GetPageInfo returns GetParticipantsTournamentParticipantsParticipantConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnection) GetPageInfo() GetParticipantsTournamentParticipantsParticipantConnectionPageInfo {
	return v.PageInfo
}

// GetNodes returns GetParticipantsTournamentParticipantsParticipantConnection.Nodes, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnection) GetNodes() []GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant {
	return v.Nodes
}

// GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant includes the requested fields of the GraphQL type Participant.
// The GraphQL type's documentation follows.
//
// A participant of a tournament; either a spectator or competitor
type GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant struct {
	Id int `json:"id"`
	// The tag that was used when the participant registered, e.g. Mang0
	GamerTag string `json:"gamerTag"`
	// Contact Info selected during registration. Falls back to User.location and/or
	// User.name if necessary. These fields are for admin use only. If you are not a
	// tournament admin or the participant being queried, these fields will be null.
	// Do not display this information publicly.
	ContactInfo GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo `json:"contactInfo"`
}

// GetId returns GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant.Id, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant) GetId() int {
	return v.Id
}

// GetGamerTag returns GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant.GamerTag, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant) GetGamerTag() string {
	return v.GamerTag
}

// GetContactInfo returns GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant.ContactInfo, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipant) GetContactInfo() GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo {
	return v.ContactInfo
}

// GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo includes the requested fields of the GraphQL type ContactInfo.
// The GraphQL type's documentation follows.
//
// Name, address, etc
type GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo struct {
	// First Name
	NameFirst string `json:"nameFirst"`
	// Last Name
	NameLast string `json:"nameLast"`
}

// GetNameFirst returns GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo.NameFirst, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo) GetNameFirst() string {
	return v.NameFirst
}

// GetNameLast returns GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo.NameLast, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionNodesParticipantContactInfo) GetNameLast() string {
	return v.NameLast
}

// GetParticipantsTournamentParticipantsParticipantConnectionPageInfo includes the requested fields of the GraphQL type PageInfo.
type GetParticipantsTournamentParticipantsParticipantConnectionPageInfo struct {
	Total int `json:"total"`
}

// GetTotal returns GetParticipantsTournamentParticipantsParticipantConnectionPageInfo.Total, and is useful for accessing the field via an interface.
func (v *GetParticipantsTournamentParticipantsParticipantConnectionPageInfo) GetTotal() int {
	return v.Total
}

// GetSetsForEntrantEvent includes the requested fields of the GraphQL type Event.
// The GraphQL type's documentation follows.
//
// An event in a tournament
type GetSetsForEntrantEvent struct {
	// Paginated sets for this Event
	Sets GetSetsForEntrantEventSetsSetConnection `json:"sets"`
}

// GetSets returns GetSetsForEntrantEvent.Sets, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEvent) GetSets() GetSetsForEntrantEventSetsSetConnection { return v.Sets }

// GetSetsForEntrantEventSetsSetConnection includes the requested fields of the GraphQL type SetConnection.
type GetSetsForEntrantEventSetsSetConnection struct {
	PageInfo GetSetsForEntrantEventSetsSetConnectionPageInfo   `json:"pageInfo"`
	Nodes    []GetSetsForEntrantEventSetsSetConnectionNodesSet `json:"nodes"`
}

// GetPageInfo returns GetSetsForEntrantEventSetsSetConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnection) GetPageInfo() GetSetsForEntrantEventSetsSetConnectionPageInfo {
	return v.PageInfo
}

// GetNodes returns GetSetsForEntrantEventSetsSetConnection.Nodes, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnection) GetNodes() []GetSetsForEntrantEventSetsSetConnectionNodesSet {
	return v.Nodes
}

// GetSetsForEntrantEventSetsSetConnectionNodesSet includes the requested fields of the GraphQL type Set.
// The GraphQL type's documentation follows.
//
// A set
type GetSetsForEntrantEventSetsSetConnectionNodesSet struct {
	Id int `json:"id"`
	// A possible spot in a set. Use this to get all entrants in a set. Use this for all bracket types (FFA, elimination, etc)
	Slots []GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot `json:"slots"`
}

// GetId returns GetSetsForEntrantEventSetsSetConnectionNodesSet.Id, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnectionNodesSet) GetId() int { return v.Id }

// GetSlots returns GetSetsForEntrantEventSetsSetConnectionNodesSet.Slots, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnectionNodesSet) GetSlots() []GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot {
	return v.Slots
}

// GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot includes the requested fields of the GraphQL type SetSlot.
// The GraphQL type's documentation follows.
//
// A slot in a set where a seed currently or will eventually exist in order to participate in the set.
type GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot struct {
	Entrant GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant `json:"entrant"`
}

// GetEntrant returns GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot.Entrant, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlot) GetEntrant() GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant {
	return v.Entrant
}

// GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant includes the requested fields of the GraphQL type Entrant.
// The GraphQL type's documentation follows.
//
// An entrant in an event
type GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant struct {
	// The entrant name as it appears in bracket: gamerTag of the participant or team name
	Name string `json:"name"`
}

// GetName returns GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant.Name, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnectionNodesSetSlotsSetSlotEntrant) GetName() string {
	return v.Name
}

// GetSetsForEntrantEventSetsSetConnectionPageInfo includes the requested fields of the GraphQL type PageInfo.
type GetSetsForEntrantEventSetsSetConnectionPageInfo struct {
	Total int `json:"total"`
}

// GetTotal returns GetSetsForEntrantEventSetsSetConnectionPageInfo.Total, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantEventSetsSetConnectionPageInfo) GetTotal() int { return v.Total }

// GetSetsForEntrantResponse is returned by GetSetsForEntrant on success.
type GetSetsForEntrantResponse struct {
	// Returns an event given its id or slug
	Event GetSetsForEntrantEvent `json:"event"`
}

// GetEvent returns GetSetsForEntrantResponse.Event, and is useful for accessing the field via an interface.
func (v *GetSetsForEntrantResponse) GetEvent() GetSetsForEntrantEvent { return v.Event }

// GetTournamentResponse is returned by GetTournament on success.
type GetTournamentResponse struct {
	// Returns a tournament given its id or slug
	Tournament GetTournamentTournament `json:"tournament"`
}

// GetTournament returns GetTournamentResponse.Tournament, and is useful for accessing the field via an interface.
func (v *GetTournamentResponse) GetTournament() GetTournamentTournament { return v.Tournament }

// GetTournamentTournament includes the requested fields of the GraphQL type Tournament.
// The GraphQL type's documentation follows.
//
// A tournament
type GetTournamentTournament struct {
	Id int `json:"id"`
	// The tournament name
	Name string `json:"name"`
	// The slug used to form the url
	Slug string `json:"slug"`
	// The short slug used to form the url
	ShortSlug string `json:"shortSlug"`
}

// GetId returns GetTournamentTournament.Id, and is useful for accessing the field via an interface.
func (v *GetTournamentTournament) GetId() int { return v.Id }

// GetName returns GetTournamentTournament.Name, and is useful for accessing the field via an interface.
func (v *GetTournamentTournament) GetName() string { return v.Name }

// GetSlug returns GetTournamentTournament.Slug, and is useful for accessing the field via an interface.
func (v *GetTournamentTournament) GetSlug() string { return v.Slug }

// GetShortSlug returns GetTournamentTournament.ShortSlug, and is useful for accessing the field via an interface.
func (v *GetTournamentTournament) GetShortSlug() string { return v.ShortSlug }

// ReportSetReportBracketSet includes the requested fields of the GraphQL type Set.
// The GraphQL type's documentation follows.
//
// A set
type ReportSetReportBracketSet struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

// GetId returns ReportSetReportBracketSet.Id, and is useful for accessing the field via an interface.
func (v *ReportSetReportBracketSet) GetId() int { return v.Id }

// GetState returns ReportSetReportBracketSet.State, and is useful for accessing the field via an interface.
func (v *ReportSetReportBracketSet) GetState() int { return v.State }

// ReportSetResponse is returned by ReportSet on success.
type ReportSetResponse struct {
	// Report set winner or game stats for a H2H bracket set. If winnerId is
	// supplied, mark set as complete. gameData parameter will overwrite any existing
	// reported game data.
	ReportBracketSet []ReportSetReportBracketSet `json:"reportBracketSet"`
}

// GetReportBracketSet returns ReportSetResponse.ReportBracketSet, and is useful for accessing the field via an interface.
func (v *ReportSetResponse) GetReportBracketSet() []ReportSetReportBracketSet {
	return v.ReportBracketSet
}

// __GenerateRegistrationTokenForEventInput is used internally by genqlient
type __GenerateRegistrationTokenForEventInput struct {
	EventId int `json:"eventId"`
	UserId  int `json:"userId"`
}

// GetEventId returns __GenerateRegistrationTokenForEventInput.EventId, and is useful for accessing the field via an interface.
func (v *__GenerateRegistrationTokenForEventInput) GetEventId() int { return v.EventId }

// GetUserId returns __GenerateRegistrationTokenForEventInput.UserId, and is useful for accessing the field via an interface.
func (v *__GenerateRegistrationTokenForEventInput) GetUserId() int { return v.UserId }

// __GenerateRegistrationTokenInput is used internally by genqlient
type __GenerateRegistrationTokenInput struct {
	UserId int `json:"userId"`
}

// GetUserId returns __GenerateRegistrationTokenInput.UserId, and is useful for accessing the field via an interface.
func (v *__GenerateRegistrationTokenInput) GetUserId() int { return v.UserId }

// __GetEntrantByNameInput is used internally by genqlient
type __GetEntrantByNameInput struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// GetSlug returns __GetEntrantByNameInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetEntrantByNameInput) GetSlug() string { return v.Slug }

// GetName returns __GetEntrantByNameInput.Name, and is useful for accessing the field via an interface.
func (v *__GetEntrantByNameInput) GetName() string { return v.Name }

// __GetEntrantsInput is used internally by genqlient
type __GetEntrantsInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __GetEntrantsInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetEntrantsInput) GetSlug() string { return v.Slug }

// __GetEntrantsOutInput is used internally by genqlient
type __GetEntrantsOutInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __GetEntrantsOutInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetEntrantsOutInput) GetSlug() string { return v.Slug }

// __GetEventInput is used internally by genqlient
type __GetEventInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __GetEventInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetEventInput) GetSlug() string { return v.Slug }

// __GetParticipantsInput is used internally by genqlient
type __GetParticipantsInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __GetParticipantsInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetParticipantsInput) GetSlug() string { return v.Slug }

// __GetSetsForEntrantInput is used internally by genqlient
type __GetSetsForEntrantInput struct {
	Slug      string `json:"slug"`
	EntrantId int    `json:"entrantId"`
}

// GetSlug returns __GetSetsForEntrantInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetSetsForEntrantInput) GetSlug() string { return v.Slug }

// GetEntrantId returns __GetSetsForEntrantInput.EntrantId, and is useful for accessing the field via an interface.
func (v *__GetSetsForEntrantInput) GetEntrantId() int { return v.EntrantId }

// __GetTournamentInput is used internally by genqlient
type __GetTournamentInput struct {
	Slug string `json:"slug"`
}

// GetSlug returns __GetTournamentInput.Slug, and is useful for accessing the field via an interface.
func (v *__GetTournamentInput) GetSlug() string { return v.Slug }

// __ReportSetInput is used internally by genqlient
type __ReportSetInput struct {
	SetId    int                       `json:"setId"`
	WinnerId int                       `json:"winnerId"`
	GameData []BracketSetGameDataInput `json:"gameData"`
}

// GetSetId returns __ReportSetInput.SetId, and is useful for accessing the field via an interface.
func (v *__ReportSetInput) GetSetId() int { return v.SetId }

// GetWinnerId returns __ReportSetInput.WinnerId, and is useful for accessing the field via an interface.
func (v *__ReportSetInput) GetWinnerId() int { return v.WinnerId }

// GetGameData returns __ReportSetInput.GameData, and is useful for accessing the field via an interface.
func (v *__ReportSetInput) GetGameData() []BracketSetGameDataInput { return v.GameData }

// The mutation executed by GenerateRegistrationToken.
const GenerateRegistrationToken_Operation = `
mutation GenerateRegistrationToken ($userId: ID!) {
	generateRegistrationToken(registration: {eventIds:[]}, userId: $userId)
}
`

func GenerateRegistrationToken(
	ctx_ context.Context,
	userId int,
) (data_ *GenerateRegistrationTokenResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GenerateRegistrationToken",
		Query:  GenerateRegistrationToken_Operation,
		Variables: &__GenerateRegistrationTokenInput{
			UserId: userId,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GenerateRegistrationTokenResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The mutation executed by GenerateRegistrationTokenForEvent.
const GenerateRegistrationTokenForEvent_Operation = `
mutation GenerateRegistrationTokenForEvent ($eventId: ID!, $userId: ID!) {
	generateRegistrationToken(registration: {eventIds:[$eventId]}, userId: $userId)
}
`

func GenerateRegistrationTokenForEvent(
	ctx_ context.Context,
	eventId int,
	userId int,
) (data_ *GenerateRegistrationTokenForEventResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GenerateRegistrationTokenForEvent",
		Query:  GenerateRegistrationTokenForEvent_Operation,
		Variables: &__GenerateRegistrationTokenForEventInput{
			EventId: eventId,
			UserId:  userId,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GenerateRegistrationTokenForEventResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetCurrentUser.
const GetCurrentUser_Operation = `
query GetCurrentUser {
	currentUser {
		id
		name
	}
}
`

func GetCurrentUser(
	ctx_ context.Context,
) (data_ *GetCurrentUserResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetCurrentUser",
		Query:  GetCurrentUser_Operation,
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetCurrentUserResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetEntrantByName.
const GetEntrantByName_Operation = `
query GetEntrantByName ($slug: String, $name: String) {
	event(slug: $slug) {
		entrants(query: {page:0,perPage:1,filter:{name:$name}}) {
			nodes {
				id
				name
			}
		}
	}
}
`

func GetEntrantByName(
	ctx_ context.Context,
	slug string,
	name string,
) (data_ *GetEntrantByNameResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetEntrantByName",
		Query:  GetEntrantByName_Operation,
		Variables: &__GetEntrantByNameInput{
			Slug: slug,
			Name: name,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetEntrantByNameResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetEntrants.
const GetEntrants_Operation = `
query GetEntrants ($slug: String) {
	event(slug: $slug) {
		entrants(query: {}) {
			nodes {
				name
			}
		}
	}
}
`

func GetEntrants(
	ctx_ context.Context,
	slug string,
) (data_ *GetEntrantsResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetEntrants",
		Query:  GetEntrants_Operation,
		Variables: &__GetEntrantsInput{
			Slug: slug,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetEntrantsResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetEntrantsOut.
const GetEntrantsOut_Operation = `
query GetEntrantsOut ($slug: String) {
	event(slug: $slug) {
		entrants(query: {page:0,perPage:100}) {
			nodes {
				name
				standing {
					isFinal
					placement
				}
			}
		}
	}
}
`

func GetEntrantsOut(
	ctx_ context.Context,
	slug string,
) (data_ *GetEntrantsOutResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetEntrantsOut",
		Query:  GetEntrantsOut_Operation,
		Variables: &__GetEntrantsOutInput{
			Slug: slug,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetEntrantsOutResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetEvent.
const GetEvent_Operation = `
query GetEvent ($slug: String) {
	event(slug: $slug) {
		id
		name
	}
}
`

func GetEvent(
	ctx_ context.Context,
	slug string,
) (data_ *GetEventResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetEvent",
		Query:  GetEvent_Operation,
		Variables: &__GetEventInput{
			Slug: slug,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetEventResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetParticipants.
const GetParticipants_Operation = `
query GetParticipants ($slug: String) {
	tournament(slug: $slug) {
		id
		name
		participants(query: {}) {
			pageInfo {
				total
			}
			nodes {
				id
				gamerTag
				contactInfo {
					nameFirst
					nameLast
				}
			}
		}
	}
}
`

func GetParticipants(
	ctx_ context.Context,
	slug string,
) (data_ *GetParticipantsResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetParticipants",
		Query:  GetParticipants_Operation,
		Variables: &__GetParticipantsInput{
			Slug: slug,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetParticipantsResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetSetsForEntrant.
const GetSetsForEntrant_Operation = `
query GetSetsForEntrant ($slug: String, $entrantId: ID) {
	event(slug: $slug) {
		sets(page: 0, perPage: 100, filters: {entrantIds:[$entrantId]}, sortType: CALL_ORDER) {
			pageInfo {
				total
			}
			nodes {
				id
				slots {
					entrant {
						name
					}
				}
			}
		}
	}
}
`

func GetSetsForEntrant(
	ctx_ context.Context,
	slug string,
	entrantId int,
) (data_ *GetSetsForEntrantResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetSetsForEntrant",
		Query:  GetSetsForEntrant_Operation,
		Variables: &__GetSetsForEntrantInput{
			Slug:      slug,
			EntrantId: entrantId,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetSetsForEntrantResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The query executed by GetTournament.
const GetTournament_Operation = `
query GetTournament ($slug: String) {
	tournament(slug: $slug) {
		id
		name
		slug
		shortSlug
	}
}
`

func GetTournament(
	ctx_ context.Context,
	slug string,
) (data_ *GetTournamentResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "GetTournament",
		Query:  GetTournament_Operation,
		Variables: &__GetTournamentInput{
			Slug: slug,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &GetTournamentResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}

// The mutation executed by ReportSet.
const ReportSet_Operation = `
mutation ReportSet ($setId: ID!, $winnerId: ID!, $gameData: [BracketSetGameDataInput]) {
	reportBracketSet(setId: $setId, winnerId: $winnerId, gameData: $gameData) {
		id
		state
	}
}
`

func ReportSet(
	ctx_ context.Context,
	setId int,
	winnerId int,
	gameData []BracketSetGameDataInput,
) (data_ *ReportSetResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "ReportSet",
		Query:  ReportSet_Operation,
		Variables: &__ReportSetInput{
			SetId:    setId,
			WinnerId: winnerId,
			GameData: gameData,
		},
	}
	var client_ graphql.Client

	client_, err_ = GetClient(ctx_)
	if err_ != nil {
		return nil, err_
	}

	data_ = &ReportSetResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return data_, err_
}
