package model

// The schema defines all the fields that exist within a user record.
type UserSchema struct {
	// Unique identifier for the user.
	// Context: embed, view, edit
	ID int `json:"id"`
	// Login name for the user.
	// Context: edit
	Username string `json:"username"`
	// Display name for the user.
	// Context: embed, view, edit
	Name string `json:"name"`
	// The nickname for the user.
	// Context: edit
	Nickname string `json:"nickname"`
	// Roles assigned to the user.
	// Context: edit
	Roles []string `json:"roles"`
	// All capabilities assigned to the user.
	// Context: edit
	Capabilities map[string]bool `json:"capabilities"`
	// Any extra capabilities assigned to the user.
	// Context: edit
	ExtraCapabilities map[string]bool `json:"extra_capabilities"`
}

type RetrieveUserMeResponse struct {
	UserSchema
}
