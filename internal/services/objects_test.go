package services_test

import (
	"testing"

	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestBaseObject_GetCornerCoordinates(t *testing.T) {
	type fields struct {
		X float64
		Y float64
	}
	tests := []struct {
		name      string
		fields    fields
		expectedX int
		expectedY int
	}{
		{"common", fields{4.3, 7.2}, 4, 7},
		{"zeros", fields{0, 0}, 0, 0},
		{"negative", fields{-4.3, -7.8}, -4, -8},
		{"integer", fields{1.0, 5.0}, 1, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseObject := &services.BaseObject{
				IsBlocked: false,
				IsDrawn:   false,
				Active:    true,
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Style:     tcell.StyleDefault.Background(tcell.ColorReset),
				Speed:     1,
				View:      "A",
			}
			x, y := baseObject.GetCornerCoordinates()
			require.Equal(t, x, tt.expectedX)
			require.Equal(t, y, tt.expectedY)
		})
	}
}
