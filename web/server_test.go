package web

import (
	"context"
	"fmt"
	"time"
)

func ExampleServer_ServeStatic() {
	s := Server{Directory: "/tmp"}
	staticServer, _, err := s.ServeStatic(true)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	if err := staticServer.Shutdown(context.TODO()); err != nil {
		panic(err)
	}

	fmt.Println("Done")

	// Output:
	// Done
}
