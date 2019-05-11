package math

import "math/big"

func Div(a float64, b float64) (result float64) {
	result, _ = big.NewFloat(0).Quo(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}
func DivI(a int64, b int64) (result float64) {
	result, _ = big.NewFloat(0).Quo(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}

func Mul(a float64, b float64) (result float64) {
	result, _ = big.NewFloat(0).Mul(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}
func MulI(a int64, b int64) (result float64) {
	result, _ = big.NewFloat(0).Mul(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}

func Add(a float64, b float64) (result float64) {
	result, _ = big.NewFloat(0).Add(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}

func AddI(a int64, b int64) (result float64) {
	result, _ = big.NewFloat(0).Add(big.NewFloat(float64(a)), big.NewFloat(float64(b))).Float64()
	return
}
