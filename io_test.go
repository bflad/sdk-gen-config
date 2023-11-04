package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "gen/test"

func TestLoad_Success(t *testing.T) {
	type args struct {
		langs   []string
		dir     string
		genYaml string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "creates config file if it doesn't exist",
			args: args{
				langs: []string{"go"},
				dir:   testDir,
			},
			want: &Config{
				ConfigVersion: Version,
				Languages: map[string]LanguageConfig{
					"go": {
						Version: "0.0.1",
					},
				},
				Generation: Generation{
					Comments: &Comments{
						DisableComments:                 false,
						OmitDescriptionIfSummaryPresent: false,
					},
					SDKClassName:         "SDK",
					SingleTagPerOp:       false,
					MaintainOpenAPIOrder: true,
				},
				New: map[string]bool{
					"go": true,
				},
				Features: map[string]map[string]string{},
			},
		},
		{
			name: "loads and upgrades pre v1.0.0 config file",
			args: args{
				langs:   []string{"go"},
				dir:     testDir,
				genYaml: readTestFile(t, "pre-v100-gen.yaml"),
			},
			want: &Config{
				ConfigVersion: Version,
				Management: &Management{
					DocChecksum:      "2bba3b8f9d211b02569b3f9aff0d34b4",
					DocVersion:       "0.3.0",
					SpeakeasyVersion: "1.3.1",
				},
				Languages: map[string]LanguageConfig{
					"go": {
						Version: "1.3.0",
						Cfg: map[string]any{
							"packageName": "github.com/speakeasy-api/speakeasy-client-sdk-go",
						},
					},
				},
				Generation: Generation{
					BaseServerURL:          "https://api.prod.speakeasyapi.dev",
					SDKClassName:           "speakeasy",
					SingleTagPerOp:         false,
					TagNamespacingDisabled: false,
					Comments: &Comments{
						OmitDescriptionIfSummaryPresent: true,
						DisableComments:                 false,
					},
				},
				Features: map[string]map[string]string{},
				New:      map[string]bool{},
			},
		},
		{
			name: "loads current version config file",
			args: args{
				langs:   []string{"go"},
				dir:     testDir,
				genYaml: readTestFile(t, "current-gen.yaml"),
			},
			want: &Config{
				ConfigVersion: Version,
				Management: &Management{
					DocChecksum:      "2bba3b8f9d211b02569b3f9aff0d34b4",
					DocVersion:       "0.3.0",
					SpeakeasyVersion: "1.3.1",
				},
				Languages: map[string]LanguageConfig{
					"go": {
						Version: "1.3.0",
						Cfg: map[string]any{
							"packageName": "github.com/speakeasy-api/speakeasy-client-sdk-go",
						},
					},
				},
				Generation: Generation{
					BaseServerURL:          "https://api.prod.speakeasyapi.dev",
					SDKClassName:           "speakeasy",
					SingleTagPerOp:         false,
					TagNamespacingDisabled: false,
					Comments: &Comments{
						DisableComments:                 false,
						OmitDescriptionIfSummaryPresent: true,
					},
				},
				Features: map[string]map[string]string{
					"go": {
						"core": "2.90.0",
					},
				},
				New: map[string]bool{},
			},
		},
		{
			name: "loads current version config file from higher level directory",
			args: args{
				langs:   []string{"go"},
				dir:     filepath.Dir(testDir),
				genYaml: readTestFile(t, "current-gen.yaml"),
			},
			want: &Config{
				ConfigVersion: Version,
				Management: &Management{
					DocChecksum:      "2bba3b8f9d211b02569b3f9aff0d34b4",
					DocVersion:       "0.3.0",
					SpeakeasyVersion: "1.3.1",
				},
				Languages: map[string]LanguageConfig{
					"go": {
						Version: "1.3.0",
						Cfg: map[string]any{
							"packageName": "github.com/speakeasy-api/speakeasy-client-sdk-go",
						},
					},
				},
				Generation: Generation{
					BaseServerURL:          "https://api.prod.speakeasyapi.dev",
					SDKClassName:           "speakeasy",
					SingleTagPerOp:         false,
					TagNamespacingDisabled: false,
					Comments: &Comments{
						DisableComments:                 false,
						OmitDescriptionIfSummaryPresent: true,
					},
				},
				Features: map[string]map[string]string{
					"go": {
						"core": "2.90.0",
					},
				},
				New: map[string]bool{},
			},
		},
		{
			name: "loads current version config file and detects new config for language",
			args: args{
				langs:   []string{"go", "typescript"},
				dir:     testDir,
				genYaml: readTestFile(t, "current-gen.yaml"),
			},
			want: &Config{
				ConfigVersion: Version,
				Management: &Management{
					DocChecksum:      "2bba3b8f9d211b02569b3f9aff0d34b4",
					DocVersion:       "0.3.0",
					SpeakeasyVersion: "1.3.1",
				},
				Languages: map[string]LanguageConfig{
					"go": {
						Version: "1.3.0",
						Cfg: map[string]any{
							"packageName": "github.com/speakeasy-api/speakeasy-client-sdk-go",
						},
					},
					"typescript": {
						Version: "0.0.1",
					},
				},
				Generation: Generation{
					BaseServerURL:          "https://api.prod.speakeasyapi.dev",
					SDKClassName:           "speakeasy",
					SingleTagPerOp:         false,
					TagNamespacingDisabled: false,
					Comments: &Comments{
						DisableComments:                 false,
						OmitDescriptionIfSummaryPresent: true,
					},
				},
				Features: map[string]map[string]string{
					"go": {
						"core": "2.90.0",
					},
				},
				New: map[string]bool{
					"typescript": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(os.TempDir(), tt.args.dir)

			err := createTempFile(tt.args.dir, tt.args.genYaml)
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			opts := []Option{
				WithUpgradeFunc(testUpdateLang),
			}

			for _, lang := range tt.args.langs {
				opts = append(opts, WithLanguages(lang))
			}

			cfg, err := Load(filepath.Join(os.TempDir(), testDir), opts...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, cfg)
			_, err = os.Stat(filepath.Join(dir, "gen.yaml"))
			assert.NoError(t, err)
		})
	}
}

func createTempFile(dir string, contents string) error {
	tmpDir := filepath.Join(os.TempDir(), dir)

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return err
	}

	if contents != "" {
		tmpFile := filepath.Join(tmpDir, "gen.yaml")
		if err := os.WriteFile(tmpFile, []byte(contents), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func readTestFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}
