package srvconfig

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {

	osArgs := []string{
		"-a", "localhost:5555",
		"-b", "http://localhost:5555",
		"-f", "/storageTest.csv",
		"-p", "./ssl",
		"-t", "0.0.0.0",
		"-s", "false",
		"-d", "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb"}

	envVars := map[string]string{
		"BASE_URL":          "http://localhost:5555",
		"SERVER_ADDRESS":    "localhost:5555",
		"FILE_STORAGE_PATH": "/storageTest.csv",
		"SSL_PATH":          "./ssl",
		"TRUSTED_SUBNET":    "0.0.0.0",
		"ENABLE_HTTPS":      "false",
		"DATABASE_DSN":      "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
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
				dbDSN:         "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
			},
		},
		{
			name: "envs only",
			opts: []configOption{IgnoreOsArgs(), withOsArgs([]string{}), WithEnvVars(envVars)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
			},
		},
		{
			name: "file only",
			opts: []configOption{WithFile(filePath)},
			want: Config{
				srvAddr:       "111",
				dbDSN:         "111",
			},
		},
		{
			name: "flags over file",
			opts: []configOption{WithEnvVars(map[string]string{}), withOsArgs(osArgs), WithFile(filePath)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
			},
		},
		{
			name: "envs over file",
			opts: []configOption{IgnoreOsArgs(), WithFile(filePath), withOsArgs([]string{}), WithEnvVars(envVars)},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
			},
		},
		{
			name: "flags over vars",
			opts: []configOption{withOsArgs(osArgs),
				WithEnvVars(map[string]string{
					"BASE_URL":          "111",
					"SERVER_ADDRESS":    "111",
					"FILE_STORAGE_PATH": "111",
					"SSL_PATH":          "111",
					"DATABASE_DSN":      "111",
					"ENABLE_HTTPS":      "111",
				})},
			want: Config{
				srvAddr:       "localhost:5555",
				dbDSN:         "user=ubuntu password=test101825 host=localhost port=5432 dbname=testdb",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.opts...); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("New().DBDSN() = %v, want %v", got.DBDSN(), tt.want.dbDSN)
				t.Errorf("New().SrvAddr() = %v, want %v", got.SrvAddr(), tt.want.srvAddr)
			}
		})
	}
}
