package snippets

func (u *UpdateSnippet) UnmarshalJSON(data []byte) error {
	type Alias UpdateSnippet
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set pointer fields to nil if they are empty
	if u.Language != nil && *u.Language == "" {
		u.Language = nil
	}
	if u.Title != nil && *u.Title == "" {
		u.Title = nil
	}
	if u.Code != nil && *u.Code == "" {
		u.Code = nil
	}
	if u.Description != nil && *u.Description == "" {
		u.Description = nil
	}
	if u.Tags != nil && len(*u.Tags) == 0 {
		u.Tags = nil
	}
	if u.Favorite != nil && !*u.Favorite {
		u.Favorite = nil
	}
	if u.Private != nil && !*u.Private {
		u.Private = nil
	}

	return nil
}