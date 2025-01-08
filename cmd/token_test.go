package cmd

import (
	"github.com/emrgen/authbase/pkg/server"
	"sync"
)

func startServer() {
	svr := server.NewServerFromEnv()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := svr.Start("5000", "5001")
		if err != nil {
			panic(err)
		}
	}()

	<-svr.Ready()
}

//func TestCreateAccessKey(t *testing.T) {
//	startServer()
//}
