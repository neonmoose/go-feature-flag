package ffexporter_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"

	"github.com/thomaspoignant/go-feature-flag/ffexporter"
	"github.com/thomaspoignant/go-feature-flag/testutils"
)

func TestLog_Export(t *testing.T) {
	type fields struct {
		LogFormat string
	}
	type args struct {
		featureEvents []ffexporter.FeatureEvent
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedLog string
		wantErr     bool
	}{
		{
			name:   "Default format",
			fields: fields{LogFormat: ""},
			args: args{featureEvents: []ffexporter.FeatureEvent{
				{Kind: "feature", ContextKind: "anonymousUser", UserKey: "ABCD", CreationDate: 1617970547, Key: "random-key",
					Variation: "Default", Value: "YO", Default: false},
			}},
			expectedLog: "^\\[" + testutils.RFC3339Regex + "\\] user=\"ABCD\", flag=\"random-key\", value=\"YO\"\n",
		},
		{
			name: "Custom format",
			fields: fields{
				LogFormat: "key=\"{{ .Key}}\" [{{ .FormattedDate}}]",
			},
			args: args{featureEvents: []ffexporter.FeatureEvent{
				{Kind: "feature", ContextKind: "anonymousUser", UserKey: "ABCD", CreationDate: 1617970547, Key: "random-key",
					Variation: "Default", Value: "YO", Default: false},
			}},
			expectedLog: "key=\"random-key\" \\[" + testutils.RFC3339Regex + "\\]\n",
		},
		{
			name: "LogFormat error",
			fields: fields{
				LogFormat: "key=\"{{ .Key}\" [{{ .FormattedDate}}]",
			},
			args: args{featureEvents: []ffexporter.FeatureEvent{
				{Kind: "feature", ContextKind: "anonymousUser", UserKey: "ABCD", CreationDate: 1617970547, Key: "random-key",
					Variation: "Default", Value: "YO", Default: false},
			}},
			expectedLog: "^\\[" + testutils.RFC3339Regex + "\\] user=\"ABCD\", flag=\"random-key\", value=\"YO\"\n",
		},
		{
			name: "Field does not exist",
			fields: fields{
				LogFormat: "key=\"{{ .UnknownKey}}\" [{{ .FormattedDate}}]",
			},
			args: args{featureEvents: []ffexporter.FeatureEvent{
				{Kind: "feature", ContextKind: "anonymousUser", UserKey: "ABCD", CreationDate: 1617970547, Key: "random-key",
					Variation: "Default", Value: "YO", Default: false},
			}},
			expectedLog: "^\\[" + testutils.RFC3339Regex + "\\] user=\"ABCD\", flag=\"random-key\", value=\"YO\"\n",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ffexporter.Log{
				LogFormat: tt.fields.LogFormat,
			}

			logFile, _ := ioutil.TempFile("", "")
			logger := log.New(logFile, "", 0)

			err := f.Export(context.Background(), logger, tt.args.featureEvents)

			if tt.wantErr {
				assert.Error(t, err, "It should return an error")
				return
			}

			assert.NoError(t, err, "Log exporter should not throw errors")

			logContent, _ := ioutil.ReadFile(logFile.Name())
			assert.Regexp(t, tt.expectedLog, string(logContent))
		})
	}
}

func TestLog_IsBulk(t *testing.T) {
	exporter := ffexporter.Log{}
	assert.False(t, exporter.IsBulk(), "File exporter is not a bulk exporter")
}
