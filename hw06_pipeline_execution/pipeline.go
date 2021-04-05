package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipeOut := in
	out := make(Bi)

	worker := func(pipeOut Out) {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case res, ok := <-pipeOut:
				if !ok {
					return
				}

				select {
				case out <- res:
					// continue
				case <-done:
					return
				}
			}
		}
	}

	for _, stage := range stages {
		pipeOut = stage(pipeOut)
	}

	go worker(pipeOut)

	// Place your code here.
	return out
}
