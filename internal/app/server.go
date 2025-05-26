package app

import (
	"fmt"

	"github.com/MingPV/clean-go-template/utils"
)

func Start() {
	app, port := SetupApp("dev")

	// Graceful shutdown
	utils.StartServerWithGracefulShutdown(app, ":"+port)
	fmt.Println("Server started on port", port)
}
