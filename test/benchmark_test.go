package test

import (
	"context"
	"fmt"
	"testing"

	"juggler/internal/juggler"
	"golang.org/x/sync/errgroup"
)

func BenchmarkJugglerCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = juggler.NewJuggler(3, 2)
	}
}

func BenchmarkJugglerReset(b *testing.B) {
	j := juggler.NewJuggler(3, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.Reset(3, 2)
	}
}

func BenchmarkJugglerGetStats(b *testing.B) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = j.GetStats()
	}
}

func BenchmarkJugglerThrowBall(b *testing.B) {
	j := juggler.NewJuggler(3, 2)
	ctx := context.Background()
	eg := &errgroup.Group{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.Reset(3, 2) 
		j.ThrowBall(ctx, eg)
	}
}

func BenchmarkJugglerConcurrentAccess(b *testing.B) {
	j := juggler.NewJuggler(10, 5)
	j.Reset(10, 5)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = j.GetStats()
		}
	})
}

func BenchmarkJugglerWithManyBalls(b *testing.B) {
	ballCounts := []int{1, 5, 10, 20, 50}
	
	for _, count := range ballCounts {
		b.Run(fmt.Sprintf("Balls_%d", count), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				j := juggler.NewJuggler(count, 2)
				j.Reset(count, 2)
				_, _, _ = j.GetStats()
			}
		})
	}
}

func BenchmarkJugglerStateOperations(b *testing.B) {
	j := juggler.NewJuggler(5, 3)
	j.Reset(5, 3)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j.IsRunning()
		j.IsFinished()
		j.AllBallsInHand()
		j.GetTotalBalls()
		j.GetJugglingTime()
	}
}

func BenchmarkMemoryAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := juggler.NewJuggler(3, 2)
		j.Reset(3, 2)
		
		ctx := context.Background()
		eg := &errgroup.Group{}
		j.ThrowBall(ctx, eg)
		j.GetStats()
		j.Stop()
	}
}
