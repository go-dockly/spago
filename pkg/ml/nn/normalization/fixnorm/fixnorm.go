// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Reference: "Improving Lexical Choice in Neural Machine Translation" by Toan Q. Nguyen and David Chiang (2018)
// (https://arxiv.org/pdf/1710.01329.pdf)
package fixnorm

import (
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
)

var (
	_ nn.Model     = &Model{}
	_ nn.Processor = &Processor{}
)

// Model is an empty model used to instantiate a new Processor.
type Model struct{}

// New returns a new model.
func New() *Model {
	return &Model{}
}

type Processor struct {
	nn.BaseProcessor
}

// NewProc returns a new processor to execute the forward step.
func (m *Model) NewProc(ctx nn.Context) nn.Processor {
	return &Processor{
		BaseProcessor: nn.BaseProcessor{
			Model:             m,
			Mode:              ctx.Mode,
			Graph:             ctx.Graph,
			FullSeqProcessing: false,
		},
	}
}

// Forward performs the forward step for each input and returns the result.
func (p *Processor) Forward(xs ...ag.Node) []ag.Node {
	g := p.Graph
	ys := make([]ag.Node, len(xs))
	eps := g.NewScalar(1e-10)
	for i, x := range xs {
		norm := g.Sqrt(g.ReduceSum(g.Square(x)))
		ys[i] = g.DivScalar(x, g.AddScalar(norm, eps))
	}
	return ys
}
