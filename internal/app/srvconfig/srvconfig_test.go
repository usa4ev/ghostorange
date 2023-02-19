package srvconfig

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {

	osArgs := []string{
		"-a", "localhost:5555",
		"-d", "db",
		"-s", "100ns",}

	envVars := map[string]string{
		"SERVER_ADDRESS":    "localhost:5555",
		"SESSION_LIFETIME":   "100ns",
		"DATABASE_DSN":      "db",
	}

	filePath := "./testdata/1.json"

	tests := []struct {
		name string
		opts []configOption
		want Config
	}{
		{
			name: "flags only",
			opts: []configOption{WithEnvVars(map[string]string{}), withOsArgs(osArgs)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "db",
				sessionLifeTime: 100,
			},
		},
		{
			name: "envs only",
			opts: []configOption{IgnoreOsArgs(), withOsArgs([]string{}), WithEnvVars(envVars)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "db",
				sessionLifeTime: 100,
			},
		},
		{
			name: "file only",
			opts: []configOption{WithFile(filePath)},
			want: Config{
				srvAddr:       "111",
				dbDSN:         "111",
				sessionLifeTime: 111,
			},
		},
		{
			name: "flags over file",
			opts: []configOption{WithEnvVars(map[string]string{}), withOsArgs(osArgs), WithFile(filePath)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "db",
				sessionLifeTime: 100,
			},
		},
		{
			name: "envs over file",
			opts: []configOption{IgnoreOsArgs(), WithFile(filePath), withOsArgs([]string{}), WithEnvVars(envVars)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "db",
				sessionLifeTime: 100,
			},
		},
		{
			name: "flags over vars",
			opts: []configOption{withOsArgs(osArgs),
				WithEnvVars(map[string]string{
					"SERVER_ADDRESS":    "0:0",
					"SESSION_LIFETIME":   "0",
					"DATABASE_DSN":      "0",
				})},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "db",
				sessionLifeTime: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.opts...); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("New().DBDSN() = %v, want %v", got.DBDSN(), tt.want.dbDSN)
				t.Errorf("New().SrvAddr() = %v, want %v", got.SrvAddr(), tt.want.srvAddr)
				t.Errorf("New().SessionLifetime() = %v, want %v", got.SessionLifetime(), tt.want.sessionLifeTime)
			}		
		})
	}
}
