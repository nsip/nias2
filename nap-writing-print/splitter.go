// splitter.go

package nap_writing_print

import "context"

//
// simple splitter that allows printing
// of audit & script files to be done
// in parallel
//
func createSplitter(ctx context.Context, in <-chan map[string]string) (
	<-chan map[string]string, <-chan map[string]string, <-chan error, error) {

	out1 := make(chan map[string]string)
	out2 := make(chan map[string]string)
	errc := make(chan error, 1)

	go func() {
		defer close(out1)
		defer close(out2)
		defer close(errc)
		for n := range in {
			// Send the data to the output channel 1 but return early
			// if the context has been cancelled.
			select {
			case out1 <- n:
			case <-ctx.Done():
				return
			}
			// Send the data to the output channel 2 but return early
			// if the context has been cancelled.
			select {
			case out2 <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out1, out2, errc, nil
}
