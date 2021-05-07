package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/go-feature-flag/internal/model"
	"github.com/thomaspoignant/go-feature-flag/testutils/testconvert"
	"testing"
)

func TestDiffCache_HasDiff(t *testing.T) {
	type fields struct {
		Deleted map[string]model.Flag
		Added   map[string]model.Flag
		Updated map[string]model.DiffUpdated
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "null fields",
			fields: fields{},
			want:   false,
		},
		{
			name: "empty fields",
			fields: fields{
				Deleted: map[string]model.Flag{},
				Added:   map[string]model.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: false,
		},
		{
			name: "only Deleted",
			fields: fields{
				Deleted: map[string]model.Flag{
					"flag": &model.FlagData{
						Percentage: testconvert.Float64(100),
						True:       testconvert.Interface(true),
						False:      testconvert.Interface(true),
						Default:    testconvert.Interface(true),
					},
				},
				Added:   map[string]model.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: true,
		},
		{
			name: "only Added",
			fields: fields{
				Added: map[string]model.Flag{
					"flag": &model.FlagData{
						Percentage: testconvert.Float64(100),
						True:       testconvert.Interface(true),
						False:      testconvert.Interface(true),
						Default:    testconvert.Interface(true),
					},
				},
				Deleted: map[string]model.Flag{},
				Updated: map[string]model.DiffUpdated{},
			},
			want: true,
		},
		{
			name: "only Updated",
			fields: fields{
				Added:   map[string]model.Flag{},
				Deleted: map[string]model.Flag{},
				Updated: map[string]model.DiffUpdated{
					"flag": {
						Before: &model.FlagData{
							Percentage: testconvert.Float64(100),
							True:       testconvert.Interface(true),
							False:      testconvert.Interface(true),
							Default:    testconvert.Interface(true),
						},
						After: &model.FlagData{
							Percentage: testconvert.Float64(100),
							True:       testconvert.Interface(true),
							False:      testconvert.Interface(true),
							Default:    testconvert.Interface(false),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "all fields",
			fields: fields{
				Added: map[string]model.Flag{
					"flag": &model.FlagData{
						Percentage: testconvert.Float64(100),
						True:       testconvert.Interface(true),
						False:      testconvert.Interface(true),
						Default:    testconvert.Interface(true),
					},
				},
				Deleted: map[string]model.Flag{
					"flag": &model.FlagData{
						Percentage: testconvert.Float64(100),
						True:       testconvert.Interface(true),
						False:      testconvert.Interface(true),
						Default:    testconvert.Interface(true),
					},
				},
				Updated: map[string]model.DiffUpdated{
					"flag": {
						Before: &model.FlagData{
							Percentage: testconvert.Float64(100),
							True:       testconvert.Interface(true),
							False:      testconvert.Interface(true),
							Default:    testconvert.Interface(true),
						},
						After: &model.FlagData{
							Percentage: testconvert.Float64(100),
							True:       testconvert.Interface(true),
							False:      testconvert.Interface(true),
							Default:    testconvert.Interface(false),
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := model.DiffCache{
				Deleted: tt.fields.Deleted,
				Added:   tt.fields.Added,
				Updated: tt.fields.Updated,
			}
			assert.Equal(t, tt.want, d.HasDiff())
		})
	}
}
