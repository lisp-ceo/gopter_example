package coin

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestAdd(t *testing.T) {
	t.Run("With unit tests", func(t *testing.T) {
		type UnitTestCase struct {
			a        Coin
			b        Coin
			overflow bool
			result   uint
		}

		tt := []UnitTestCase{
			// normal
			UnitTestCase{
				a:        Coin(5),
				b:        Coin(10),
				overflow: false,
				result:   15,
			},
			// overflow
			UnitTestCase{
				a:        Coin(1),
				b:        Coin(MaxCoinValue),
				overflow: true,
			},
			// boundary
			UnitTestCase{
				a:        Coin(2),
				b:        Coin(MaxCoinValue - 3),
				overflow: false,
				result:   MaxCoinValue - 1,
			},
		}

		for n, tc := range tt {
			err := tc.a.Add(&tc.b)
			if tc.overflow && err == nil {
				t.Errorf("Expected error, got none for %d", n)
			}
			if !tc.overflow && err != nil {
				t.Errorf("Expected no error, got %#v for %d", err, n)
			}
			if !tc.overflow && err != nil && tc.result != uint(tc.a) {
				t.Errorf("Expected: %d.\nGot: %d for %d", tc.result, uint(tc.a), n)
			}
		}
	})

	t.Run("With property tests", func(t *testing.T) {
		parameters := gopter.DefaultTestParameters()
		parameters.MinSuccessfulTests = 10000
		properties := gopter.NewProperties(parameters)

		properties.Property("For normal additions the coins sum to their components", prop.ForAll(
			func(values []uint) (bool, error) {
				coinA := Coin(values[0])
				coinB := Coin(values[1])
				if err := coinA.Add(&coinB); err != nil {
					return false, err
				}
				if uint(coinA) != (values[0] + values[1]) {
					return false, fmt.Errorf("Expected: %d\nGot: %d", uint(coinA), uint(values[0])+uint(values[1]))
				}
				return true, nil
			},
			gen.SliceOfN(2, gen.UIntRange(0, MaxCoinValue)).SuchThat(func(values []uint) bool {
				return (values[0] + values[1]) < MaxCoinValue
			})))
		properties.Property("Additions of values over the overflow threshold are not allowed", prop.ForAll(
			func(values []uint) (bool, error) {
				coinA := Coin(values[0])
				coinB := Coin(values[1])
				if err := coinA.Add(&coinB); err != nil {
					return true, nil
				}
				return false, fmt.Errorf("Expected error, got nil for %#v + %#v", coinA, coinB)
			},
			gen.SliceOfN(2, gen.UIntRange(0, MaxCoinValue)).SuchThat(func(values []uint) bool {
				return values[0]+values[1] > MaxCoinValue
			})))
		properties.Property("Additions of values under the overflow threshold are allowed", prop.ForAll(
			func(a, b uint) (bool, error) {
				coinA := Coin(a)
				coinB := Coin(b)
				if err := coinA.Add(&coinB); err != nil {
					return false, err
				}
				if uint(coinA) != (a + b) {
					return false, fmt.Errorf("Expected: %d\nGot: %d", uint(coinA), uint(a)+uint(b))
				}
				return true, nil
			},
			gen.UIntRange(0, 3),
			gen.UIntRange(MaxCoinValue-3, MaxCoinValue-3)))

		properties.TestingRun(t)
	})
}
