// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"brillion.io/spago/pkg/mat"
)

// SoftShrink(x) = ​x − λ if x > λ; x + λ if x < −λ; 0 otherwise ​
type SoftShrink struct {
	x      Operand
	lambda Operand // scalar
}

func NewSoftShrink(x, lambda Operand) *SoftShrink {
	return &SoftShrink{x: x, lambda: lambda}
}

// Forward computes the output of the function.
func (r *SoftShrink) Forward() mat.Matrix {
	y := r.x.Value().ZerosLike()
	y.ApplyWithAlpha(softShrink, r.x.Value(), r.lambda.Value().Scalar())
	return y
}

func (r *SoftShrink) Backward(gy mat.Matrix) {
	if r.x.RequiresGrad() {
		gx := r.x.Value().ZerosLike()
		gx.ApplyWithAlpha(softShrinkDeriv, r.x.Value(), r.lambda.Value().Scalar())
		gx.ProdInPlace(gy)
		r.x.PropagateGrad(gx)
	}
}