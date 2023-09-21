package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Represents an IPLD link. Links are generally but not necessarily implemented as CIDs
type Link string

func (l Link) String() string {
	return string(l)
}

// UnmarshalGQLContext implements the graphql.ContextUnmarshaler interface
func (l *Link) UnmarshalGQLContext(_ context.Context, v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Link must be a string")
	}
	*l = Link(s)
	return nil
}

// MarshalGQLContext implements the graphql.ContextMarshaler interface
func (l Link) MarshalGQLContext(_ context.Context, w io.Writer) error {
	encodable := map[string]string{
		"/": l.String(),
	}
	return json.NewEncoder(w).Encode(encodable)
}
