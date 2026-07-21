// Command kernel is the DevOS kernel binary. It bootstraps the kernel and runs
// until interrupted. This is the skeleton entry point; no business logic,
// networking, or authentication is performed here (Sprint 0, Component 2).
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/kernel"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	k, err := kernel.New()
	if err != nil {
		log.Fatalf("kernel bootstrap failed: %v", err)
	}

	log.Printf("devos kernel started (env=%s service=%s)", k.Config.Environment, k.Config.ServiceName)
	if err := k.Run(ctx); err != nil {
		log.Fatalf("kernel run failed: %v", err)
	}
	log.Print("devos kernel stopped")
}
