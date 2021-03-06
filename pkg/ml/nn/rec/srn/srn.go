// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"log"
)

var (
	_ nn.Model     = &Model{}
	_ nn.Processor = &Processor{}
)

// Model contains the serializable parameters.
type Model struct {
	W    *nn.Param `type:"weights"`
	WRec *nn.Param `type:"weights"`
	B    *nn.Param `type:"biases"`
}

// New returns a new model with parameters initialized to zeros.
func New(in, out int) *Model {
	return &Model{
		W:    nn.NewParam(mat.NewEmptyDense(out, in)),
		WRec: nn.NewParam(mat.NewEmptyDense(out, out)),
		B:    nn.NewParam(mat.NewEmptyVecDense(out)),
	}
}

type State struct {
	Y ag.Node
}

type Processor struct {
	nn.BaseProcessor
	w      ag.Node
	wRec   ag.Node
	b      ag.Node
	States []*State
}

// NewProc returns a new processor to execute the forward step.
func (m *Model) NewProc(ctx nn.Context) nn.Processor {
	g := ctx.Graph
	return &Processor{
		BaseProcessor: nn.BaseProcessor{
			Model:             m,
			Mode:              ctx.Mode,
			Graph:             ctx.Graph,
			FullSeqProcessing: false,
		},
		States: nil,
		w:      g.NewWrap(m.W),
		wRec:   g.NewWrap(m.WRec),
		b:      g.NewWrap(m.B),
	}
}

func (p *Processor) SetInitialState(state *State) {
	if len(p.States) > 0 {
		log.Fatal("srn: the initial state must be set before any input")
	}
	p.States = append(p.States, state)
}

// Forward performs the forward step for each input and returns the result.
func (p *Processor) Forward(xs ...ag.Node) []ag.Node {
	ys := make([]ag.Node, len(xs))
	for i, x := range xs {
		s := p.forward(x)
		p.States = append(p.States, s)
		ys[i] = s.Y
	}
	return ys
}

func (p *Processor) LastState() *State {
	n := len(p.States)
	if n == 0 {
		return nil
	}
	return p.States[n-1]
}

// y = tanh(w (dot) x + b + wRec (dot) yPrev)
func (p *Processor) forward(x ag.Node) (s *State) {
	s = new(State)
	yPrev := p.prev()
	s.Y = p.Graph.Tanh(nn.Affine(p.Graph, p.b, p.w, x, p.wRec, yPrev))
	return
}

func (p *Processor) prev() (yPrev ag.Node) {
	s := p.LastState()
	if s != nil {
		yPrev = s.Y
	}
	return
}
