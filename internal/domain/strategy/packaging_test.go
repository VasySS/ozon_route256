package strategy

import (
	"testing"

	"workshop-1/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewPackaging(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Packaging
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание упаковки",
			args: args{
				name: "wrap",
			},
			want:    Wrap{},
			wantErr: assert.NoError,
		},
		{
			name: "успешное создание коробки",
			args: args{
				name: "box",
			},
			want:    Box{},
			wantErr: assert.NoError,
		},
		{
			name: "успешное создание пакета",
			args: args{
				name: "bag",
			},
			want:    Bag{},
			wantErr: assert.NoError,
		},
		{
			name: "неизвестный тип упаковки",
			args: args{
				name: "foobar",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPackaging(tt.args.name)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPackagingWithWrap_Apply(t *testing.T) {
	type args struct {
		order *domain.Order
	}

	tests := []struct {
		name    string
		p       PackagingWithWrap
		args    args
		want    *domain.Order
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание упаковки с пленкой",
			p:    PackagingWithWrap{MainPackaging: Bag{}},
			args: args{
				order: &domain.Order{Weight: 9, Price: 10},
			},
			want:    &domain.Order{Weight: 9, Price: 16},
			wantErr: assert.NoError,
		},
		{
			name: "попытка использовать плёнку дважды",
			p:    PackagingWithWrap{MainPackaging: Wrap{}},
			args: args{
				order: &domain.Order{Weight: 9, Price: 10},
			},
			want:    &domain.Order{Weight: 9, Price: 10},
			wantErr: assert.Error,
		},
		{
			name: "попытка использовать доп. плёнку без основной упаковки",
			p:    PackagingWithWrap{MainPackaging: nil},
			args: args{
				order: &domain.Order{Weight: 9, Price: 10},
			},
			want:    &domain.Order{Weight: 9, Price: 10},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Apply(tt.args.order)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, tt.args.order)
		})
	}
}

func TestBag_Apply(t *testing.T) {
	type args struct {
		order *domain.Order
	}

	tests := []struct {
		name    string
		b       Bag
		args    args
		want    *domain.Order
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание и применение пакета",
			b:    Bag{},
			args: args{
				order: &domain.Order{Weight: 9, Price: 10},
			},
			want:    &domain.Order{Weight: 9, Price: 15},
			wantErr: assert.NoError,
		},
		{
			name: "попытка использовать пакет при превышенном весе",
			b:    Bag{},
			args: args{
				order: &domain.Order{Weight: 10, Price: 10},
			},
			want:    &domain.Order{Weight: 10, Price: 10},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Apply(tt.args.order)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, tt.args.order)
		})
	}
}

func TestBox_Apply(t *testing.T) {
	type args struct {
		order *domain.Order
	}

	tests := []struct {
		name    string
		b       Box
		args    args
		want    *domain.Order
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание и применение коробки",
			b:    Box{},
			args: args{
				order: &domain.Order{Weight: 29, Price: 10},
			},
			want:    &domain.Order{Weight: 29, Price: 30},
			wantErr: assert.NoError,
		},
		{
			name: "попытка использовать коробку при превышенном весе",
			b:    Box{},
			args: args{
				order: &domain.Order{Weight: 30, Price: 10},
			},
			want:    &domain.Order{Weight: 30, Price: 10},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Apply(tt.args.order)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, tt.args.order)
		})
	}
}
