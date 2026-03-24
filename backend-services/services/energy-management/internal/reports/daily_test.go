package reports

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseClientID(t *testing.T) {
	_, err := ParseClientID("not-a-uuid")
	require.Error(t, err)
	id := uuid.New()
	p, err := ParseClientID(id.String())
	require.NoError(t, err)
	require.Equal(t, id, p)
}
