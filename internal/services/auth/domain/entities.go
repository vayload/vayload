/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import (
	"strconv"
	"strings"
	"time"
)

type UserPolicyChannel int

// DEFAULT PROFILE IDS
// FREE -> 546
// PREMIUM -> 547
// PRO -> 548

const (
	PolicyPhoneChannel UserPolicyChannel = 1
	PolicyEmailChannel UserPolicyChannel = 2
)

type RoleQuotas map[string]int
type RolePermissions map[string]bool

type RawUserPolicy struct {
	ID           string           `json:"id"`
	ClientId     string           `json:"client_id"`
	Profile      UserProfile      `json:"profile"`
	Quotas       RoleQuotas       `json:"quotas"`
	Permissions  RolePermissions  `json:"permissions"`
	Meta         map[string]any   `json:"meta"`
	Subscription UserSubscription `json:"subscription"`
	Signature    string           `json:"signature"`
	CreatedAt    string           `json:"created_at"`
}

type UserProfile struct {
	Name        string `json:"name"`
	Paid        bool   `json:"paid"`
	Active      bool   `json:"active"`
	Country     string `json:"country"`
	Description string `json:"description"`
	ProfileID   int64  `json:"profile_id"`
	TagName     string `json:"tag_name"`
	UserType    string `json:"user_type"`
}

type UserQuotas struct {
	Lmed  int `json:"lmed"`
	Ltok  int `json:"ltok"`
	Izco  int `json:"izco"`
	Izutc int `json:"izutc"`
}

type UserPermissions struct {
	Mmc   bool `json:"mmc"`
	Mmm   bool `json:"mmm"`
	Pmed  bool `json:"pmed"`
	Mmn   bool `json:"mmn"`
	Mmsp  bool `json:"mmsp"`
	Mmsd  bool `json:"mmsd"`
	Mmcia bool `json:"mmcia"`
	Ptok  bool `json:"ptok"`
	Mma   bool `json:"mma"`
	Esms  bool `json:"esms"`
	Eema  bool `json:"eema"`
	Ewha  bool `json:"ewha"`
	Evad  bool `json:"evad"`
	Cstma bool `json:"cstma"`
	Cstme bool `json:"cstme"`
}

type UserSubscription struct {
	ID           *int    `json:"id"`
	Method       *string `json:"method"`
	Status       *string `json:"status"`
	PlanID       *string `json:"plan_id"`
	RecurrenceID *string `json:"recurrence_id"`
}

type LastPayment struct {
	ID          int       `json:"id"`
	Transaction string    `json:"transaction"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type Subscription struct {
	ID              string      `json:"id"`
	RecurrenceID    string      `json:"recurrence_id"`
	PaymentMethod   string      `json:"payment_method"`
	PlanID          string      `json:"plan_id"`
	Status          string      `json:"status"`
	StartDate       string      `json:"start_date"`
	EndDate         *string     `json:"end_date"`
	NextBillingDate *string     `json:"next_billing_date"`
	OriginID        *string     `json:"origin_id"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	LastPayment     LastPayment `json:"last_payment"`
	PaymentsCount   int         `json:"payments_count"`
}

func (s *Subscription) IsActive() bool {
	return strings.ToLower(s.Status) == "active"
}

func (s *Subscription) CreatedAtTime() (time.Time, error) {
	return time.Parse(time.RFC3339, s.CreatedAt)
}

func (s *Subscription) UpdatedAtTime() (time.Time, error) {
	return time.Parse(time.RFC3339, s.UpdatedAt)
}

type UserPolicy struct {
	Id           int64             `json:"id"`
	ProfileId    int64             `json:"profile_id"`
	Permissions  RolePermissions   `json:"permissions"`
	Quotas       RoleQuotas        `json:"quotas"`
	Meta         map[string]any    `json:"meta"`
	Subscription *UserSubscription `json:"subscription"`
	Signature    string            `json:"signature"`
	CreatedAt    string            `json:"created_at"`
	Plan         string            `json:"plan"` // Added Plan field to UserPolicy (e.g., "free", "premium", "pro")

	RawUserProfile *UserProfile `json:"-"` // Added raw user profile for additional information
}

var (
	MappingPermissions = map[string]string{
		"calendar":                               "MMC",
		"medication":                             "MMM",
		"medication-records-limit":               "LMED",
		"purchase-additional-medication-records": "PMED",
		"notifications":                          "MMN",
		"subscription-digital-product":           "MMSP",
		"subscription-dosage-medication-pack":    "MMSD",
		"coach-ia":                               "MMCIA",
		"token-limit-ia":                         "LTOK",
		"purchase-additional-tokens":             "LTOK",
		"medication-adherence":                   "MMA",
		"sms-enable":                             "ESMS",
		"email-enable":                           "EEMA",
		"whatsapp-enable":                        "EWHA",
		"vademecum-enable":                       "EVAD",
		"iso-2-letter-country-code":              "ISOC",
		"cybersalud-telemedicine-availability":   "CSTMA",
		"cybersalud-telemedicine":                "CSTME",
	}

	MappingQuotas = map[string]string{
		"medication-records":                     "LMED",
		"purchase-additional-medication-records": "PMED",
		"notifications":                          "MMN",
		"subscription-digital-product":           "MMSP",
		"subscription-dosage-medication-pack":    "MMSD",
		"coach-ia":                               "MMCIA",
	}
)

func NewUserPolicy(raw *RawUserPolicy) *UserPolicy {
	client_id, _ := strconv.ParseInt(raw.ClientId, 10, 64)

	return &UserPolicy{
		Id:          client_id,
		ProfileId:   raw.Profile.ProfileID,
		Permissions: raw.Permissions,
		Quotas:      raw.Quotas,
		Meta:        raw.Meta,
		Subscription: &UserSubscription{
			ID:           raw.Subscription.ID,
			Method:       raw.Subscription.Method,
			Status:       raw.Subscription.Status,
			PlanID:       raw.Subscription.PlanID,
			RecurrenceID: raw.Subscription.RecurrenceID,
		},
		Signature: raw.Signature,
		CreatedAt: raw.CreatedAt,
		Plan:      raw.Profile.UserType,
		RawUserProfile: &UserProfile{
			Name:        raw.Profile.Name,
			Paid:        raw.Profile.Paid,
			Active:      raw.Profile.Active,
			Country:     raw.Profile.Country,
			Description: raw.Profile.Description,
			ProfileID:   raw.Profile.ProfileID,
			TagName:     raw.Profile.TagName,
			UserType:    raw.Profile.UserType,
		},
	}
}

func (policy *UserPolicy) HasPermissions(rules []string) bool {
	if len(rules) == 0 {
		return false
	}

	for _, scope := range rules {
		if rawRule, ok := MappingPermissions[scope]; ok {
			if hasPermission, exists := policy.Permissions[rawRule]; exists && hasPermission {
				return true
			}
		}
	}

	return false
}

func (policy *UserPolicy) HasQuotaFullyUsed(rule string, used int) bool {
	if rawRule, ok := MappingQuotas[rule]; ok {
		if limit, exists := policy.Quotas[rawRule]; exists {
			return limit <= used
		}
	}

	return false
}

// Gets the profile id asigned to client
func (policy *UserPolicy) GetProfileId() int64 { return policy.ProfileId }

// Gets the client id assigned to user
func (policy *UserPolicy) GetClientId() int64 { return policy.Id }

// Get the signature of permissions + quotas + meta
func (policy *UserPolicy) GetSignature() string { return policy.Signature }

// Get current plan asigned to client
func (policy *UserPolicy) GetPlan() string { return policy.Plan }

// Get current user profile
func (policy *UserPolicy) GetUserProfile() *UserProfile { return policy.RawUserProfile }

// Check if current user is free
func (policy *UserPolicy) IsFreePlan() bool { return strings.ToLower(policy.Plan) == "free" }

// IsPremiumPlan checks if the current user is on a premium plan
func (policy *UserPolicy) IsPremiumPlan() bool { return strings.ToLower(policy.Plan) == "premium" }

// IsTrialPlan checks if the current user is on a trial plan
func (policy *UserPolicy) IsTrialPlan() bool { return strings.ToLower(policy.Plan) == "trial" }

func (policy *UserPolicy) GetSubscription() *UserSubscription { return policy.Subscription }
