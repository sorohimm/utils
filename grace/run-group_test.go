package grace

import (
	"log"
	"reflect"
	"testing"

	"github.com/oklog/run"
)

func TestNewRunGroup(t *testing.T) {
	tests := []struct {
		name string
		want *RunGroup
	}{
		{
			name: "ok",
			want: &RunGroup{&run.Group{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRunGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRunGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunGroup_Add(t *testing.T) {
	type fields struct {
		g *run.Group
	}
	type args struct {
		exec  onExec
		inter onInterrupt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"ok",
			fields{g: &run.Group{}},
			args{
				exec: func() error {
					return nil
				},
				inter: func(err error) {
					log.Println(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &RunGroup{
				g: tt.fields.g,
			}
			o.Add(tt.args.exec, tt.args.inter)
		})
	}
}

func TestRunGroup_Run(t *testing.T) {
	type fields struct {
		g *run.Group
	}
	type args struct {
		onShutdown func(error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &RunGroup{
				g: tt.fields.g,
			}
			if err := o.Run(tt.args.onShutdown); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
