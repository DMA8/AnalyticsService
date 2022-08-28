package analytics

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	mock_ports "gitlab.com/g6834/team31/analytics/internal/mocks"
	"gitlab.com/g6834/team31/analytics/internal/ports"
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestService_ApprovedTasks(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := mock_ports.NewMockDbInterface(ctrl)
	type fields struct {
		db ports.DbInterface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Counter
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{mockDB},
			args:    args{ctx: ctx},
			want:    models.Counter{Count: 5},
			wantErr: false,
		},
	}
	mockDB.EXPECT().ApprovedTasks(gomock.Any()).Times(1).Return(models.Counter{Count: 5}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.fields.db,
			}
			got, err := s.ApprovedTasks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ApprovedTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ApprovedTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_DeclinedTasks(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := mock_ports.NewMockDbInterface(ctrl)
	type fields struct {
		db ports.DbInterface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Counter
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{mockDB},
			args:    args{ctx: ctx},
			want:    models.Counter{Count: 5},
			wantErr: false,
		},
	}
	mockDB.EXPECT().DeclinedTasks(gomock.Any()).Times(1).Return(models.Counter{Count: 5}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.fields.db,
			}
			got, err := s.DeclinedTasks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Declined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Declined() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SummaryTime(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := mock_ports.NewMockDbInterface(ctrl)
	type fields struct {
		db ports.DbInterface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.SummaryTime
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{mockDB},
			args: args{ctx: ctx},
			want: []models.SummaryTime{{TaskId: "test", Duration: 1024.0}},
			wantErr: false,
		},
	}
	mockDB.EXPECT().SummaryTime(gomock.Any()).Times(1).Return([]models.SummaryTime{{TaskId: "test", Duration: 1024.0}}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.fields.db,
			}
			got, err := s.SummaryTime(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SummaryTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.SummaryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
