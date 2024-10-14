package auth

// Identity denotes an identity that has been authenticated, which may contain some resolved authorizations.
type Identity struct {
	DID               string
	AuthorizedActions []string
}

// Can check if the identity is authorized to perform a specific action.
func (i Identity) Can(action string) bool {
	for _, a := range i.AuthorizedActions {
		if a == action {
			return true
		}
	}
	return false
}
