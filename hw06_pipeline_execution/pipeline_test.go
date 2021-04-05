package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

type mapper func(interface{}) interface{}

type PipelineTestSuite struct {
	suite.Suite
	stages []Stage
}

func TestPipeline(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}

func (s *PipelineTestSuite) SetupSuite() {
	g := func(_ string, f mapper) Stage {
		return s.createMapper(s.createDelayedMapper(f, sleepPerStage))
	}

	s.stages = []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}
}

func (s *PipelineTestSuite) TestSimple() {
	in := make(Bi)
	data := []interface{}{1, 2, 3, 4, 5}

	go s.fillChan(in, data...)

	result := make([]interface{}, 0, len(data))
	out := ExecutePipeline(in, nil, s.stages...)
	expected := []interface{}{"102", "104", "106", "108", "110"}

	// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
	require.Eventually(s.T(), func() bool {
		for {
			select {
			case r, ok := <-out:
				if !ok {
					return true
				}

				result = append(result, r)
			default:
				return false
			}
		}
	}, (time.Duration(len(s.stages)+len(data)-1)*sleepPerStage)+fault, time.Millisecond)

	require.Equal(s.T(), result, expected)
}

func (s *PipelineTestSuite) TestDone() {
	in := make(Bi)
	done := make(Bi)
	data := []interface{}{1, 2, 3, 4, 5}

	// Abort after 200ms
	abortDur := sleepPerStage * 2
	go func() {
		<-time.After(abortDur)
		close(done)
	}()

	go s.fillChan(in, data...)

	result := make([]interface{}, 0, len(data))
	out := ExecutePipeline(in, done, s.stages...)

	require.Eventually(s.T(), func() bool {
		for {
			select {
			case r, ok := <-out:
				if !ok {
					return true
				}

				result = append(result, r)
			default:
				return false
			}
		}
	}, abortDur+fault, time.Millisecond)

	require.Empty(s.T(), result)
}

func (s *PipelineTestSuite) TestEmptyInput() {
	in := make(Bi, 1)
	done := make(Bi)

	go s.fillChan(in)

	result := s.collectChan(ExecutePipeline(in, done, s.stages...))

	require.Empty(s.T(), result)
}

func (s *PipelineTestSuite) TestEmptyStages() {
	in := make(Bi)
	done := make(Bi)
	data := []interface{}{1, 2, 3, 4, 5}

	go s.fillChan(in, data...)

	result := s.collectChan(ExecutePipeline(in, done))

	require.Equal(s.T(), result, data)
}

func (s *PipelineTestSuite) TestLimit() {
	in := make(Bi)
	done := make(Bi)
	tick := time.Tick(sleepPerStage)

	go func() {
		defer close(in)

		for {
			select {
			case <-done:
				return
			case t := <-tick:
				in <- t
			}
		}
	}()

	go func() {
		<-time.After(5 * sleepPerStage)

		close(done)
	}()

	result := s.collectChan(ExecutePipeline(in, done))

	require.LessOrEqual(s.T(), len(result), 5)
}

func (s *PipelineTestSuite) TestFilter() {
	in := make(Bi)
	done := make(Bi)
	data := []interface{}{1, 2, 3, 4, 5}

	go s.fillChan(in, data...)

	stages := []Stage{
		s.createFilter(func(v interface{}) bool { return v.(int)%2 == 0 }),
	}

	result := s.collectChan(ExecutePipeline(in, done, stages...))
	expected := []interface{}{2, 4}

	require.Equal(s.T(), result, expected)
}

func (s *PipelineTestSuite) TestReduce() {
	in := make(Bi)
	done := make(Bi)
	data := []interface{}{1, 2, 3, 4, 5}

	go s.fillChan(in, data...)

	stages := []Stage{
		s.createReducer(func(a interface{}, v interface{}) interface{} { return a.(int) + v.(int) }, 0),
	}

	result := s.collectChan(ExecutePipeline(in, done, stages...))
	expected := []interface{}{15}

	require.Equal(s.T(), result, expected)
}

func (s *PipelineTestSuite) fillChan(ch Bi, data ...interface{}) {
	defer close(ch)

	for _, v := range data {
		ch <- v
	}
}

func (s *PipelineTestSuite) collectChan(ch In) (r []interface{}) {
	for v := range ch {
		r = append(r, v)
	}

	return r
}

func (s *PipelineTestSuite) createDelayedMapper(fn mapper, delay time.Duration) mapper {
	return func(v interface{}) interface{} {
		time.Sleep(delay)

		return fn(v)
	}
}

func (s *PipelineTestSuite) createMapper(fn mapper) Stage {
	return func(in In) Out {
		out := make(Bi)

		go func() {
			defer close(out)

			for v := range in {
				out <- fn(v)
			}
		}()

		return out
	}
}

func (s *PipelineTestSuite) createFilter(p func(v interface{}) bool) Stage {
	return func(in In) Out {
		out := make(Bi)

		go func() {
			defer close(out)

			for v := range in {
				if !p(v) {
					continue
				}

				out <- v
			}
		}()

		return out
	}
}

func (s *PipelineTestSuite) createReducer(r func(a interface{}, v interface{}) interface{}, acc interface{}) Stage {
	return func(in In) Out {
		out := make(Bi)

		go func() {
			defer close(out)

			for v := range in {
				acc = r(acc, v)
			}

			out <- acc
		}()

		return out
	}
}
