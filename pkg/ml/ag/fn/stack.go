// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import "brillion.io/spago/pkg/mat"

type Stack struct {
	xs []Operand
}

func NewStack(xs []Operand) *Stack {
	return &Stack{xs: xs}
}

// Forward computes the output of the function.
func (r *Stack) Forward() mat.Matrix {
	rows := r.xs[0].Value().Rows()
	cols := len(r.xs)
	ms := mat.NewEmptyDense(rows, cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			ms.Set(r.xs[i].Value().At(j, 0), i, j)
		}
	}
	return ms
}

func (r *Stack) Backward(gy mat.Matrix) {
	sizes := make([]int, len(r.xs))
	for i, x := range r.xs {
		sizes[i] = x.Value().Size()
	}
	for i, gx := range gy.(*mat.Dense).SplitV(sizes...) {
		if r.xs[i].RequiresGrad() {
			r.xs[i].PropagateGrad(gx)
		}
	}
}