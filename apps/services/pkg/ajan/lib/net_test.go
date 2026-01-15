package lib_test

import (
	"testing"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectLocalNetwork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		addr      string
		wantLocal bool
		wantErr   bool
	}{
		{ //nolint:exhaustruct
			name:      "loopback_ipv4",
			addr:      "127.0.0.1",
			wantLocal: true,
		},
		{ //nolint:exhaustruct
			name:      "loopback_ipv4_with_port",
			addr:      "127.0.0.1:8080",
			wantLocal: true,
		},
		{ //nolint:exhaustruct
			name:      "remote_ipv4",
			addr:      "203.0.113.1",
			wantLocal: false,
		},
		{ //nolint:exhaustruct
			name:      "remote_ipv4_with_port",
			addr:      "203.0.113.1:8080",
			wantLocal: false,
		},
		{ //nolint:exhaustruct
			name:    "invalid_addr",
			addr:    "not-an-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests { //nolint:varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isLocal, err := lib.DetectLocalNetwork(tt.addr)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, lib.ErrInvalidIPAddress)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantLocal, isLocal)
		})
	}
}
