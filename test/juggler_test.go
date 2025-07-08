package test

import (
	"context"
	"testing"
	"time"

	"juggler/internal/juggler"

	"golang.org/x/sync/errgroup"
)

func TestNewJuggler(t *testing.T) {
	tests := []struct {
		name        string
		totalBalls  int
		timeMinutes int
		expectBalls int
	}{
		{
			name:        "Zero balls",
			totalBalls:  0,
			timeMinutes: 2,
			expectBalls: 0,
		},
		{
			name:        "Three balls",
			totalBalls:  3,
			timeMinutes: 2,
			expectBalls: 3,
		},
		{
			name:        "Five balls",
			totalBalls:  5,
			timeMinutes: 5,
			expectBalls: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := juggler.NewJuggler(tt.totalBalls, tt.timeMinutes)

			if j.GetTotalBalls() != tt.totalBalls {
				t.Errorf("Expected total balls to be %d, got %d", tt.totalBalls, j.GetTotalBalls())
			}

			expectedTime := time.Duration(tt.timeMinutes) * time.Minute
			if j.GetJugglingTime() != expectedTime {
				t.Errorf("Expected juggling time to be %v, got %v", expectedTime, j.GetJugglingTime())
			}

			// Check initial state
			inHand, inAir, balls := j.GetStats()
			if inHand != tt.expectBalls {
				t.Errorf("Expected %d balls in hand, got %d", tt.expectBalls, inHand)
			}

			if inAir != 0 {
				t.Errorf("Expected 0 balls in air, got %d", inAir)
			}

			if len(balls) != tt.expectBalls {
				t.Errorf("Expected %d balls in details, got %d", tt.expectBalls, len(balls))
			}

			if !j.IsFinished() {
				t.Error("Expected new juggler to be finished")
			}

			if j.IsRunning() {
				t.Error("Expected new juggler to not be running")
			}
		})
	}
}

func TestJugglerReset(t *testing.T) {
	j := juggler.NewJuggler(3, 2)

	j.Reset(5, 3)

	if j.GetTotalBalls() != 5 {
		t.Errorf("Expected total balls to be 5 after reset, got %d", j.GetTotalBalls())
	}

	expectedTime := time.Duration(3) * time.Minute
	if j.GetJugglingTime() != expectedTime {
		t.Errorf("Expected juggling time to be %v after reset, got %v", expectedTime, j.GetJugglingTime())
	}

	inHand, inAir, balls := j.GetStats()
	if inHand != 5 {
		t.Errorf("Expected 5 balls in hand after reset, got %d", inHand)
	}

	if inAir != 0 {
		t.Errorf("Expected 0 balls in air after reset, got %d", inAir)
	}

	if len(balls) != 5 {
		t.Errorf("Expected 5 balls in details after reset, got %d", len(balls))
	}

	if j.IsFinished() {
		t.Error("Expected juggler to not be finished after reset")
	}

	if !j.IsRunning() {
		t.Error("Expected juggler to be ready to run after reset")
	}
}

func TestJugglerThrowBall(t *testing.T) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2) 

	ctx := context.Background()
	eg := &errgroup.Group{}

	success := j.ThrowBall(ctx, eg)
	if !success {
		t.Error("Expected ThrowBall to succeed")
	}

	inHand, inAir, balls := j.GetStats()
	if inHand != 2 {
		t.Errorf("Expected 2 balls in hand after throwing, got %d", inHand)
	}

	if inAir != 1 {
		t.Errorf("Expected 1 ball in air after throwing, got %d", inAir)
	}

	if len(balls) != 3 {
		t.Errorf("Expected 3 balls in details after throwing, got %d", len(balls))
	}

	var flightBall *juggler.Ball
	for _, ball := range balls {
		if ball.Status == "in_flight" {
			flightBall = &ball
			break
		}
	}

	if flightBall == nil {
		t.Error("Expected to find a ball in flight")
	} else {
		if flightBall.FlightTime < 5 || flightBall.FlightTime > 10 {
			t.Errorf("Expected flight time to be between 5-10 seconds, got %d", flightBall.FlightTime)
		}
	}
}

func TestJugglerThrowBallNoAvailable(t *testing.T) {
	j := juggler.NewJuggler(0, 2)
	j.Reset(0, 2) 

	ctx := context.Background()
	eg := &errgroup.Group{}

	success := j.ThrowBall(ctx, eg)
	if success {
		t.Error("Expected ThrowBall to fail when no balls available")
	}

	inHand, inAir, balls := j.GetStats()
	if inHand != 0 {
		t.Errorf("Expected 0 balls in hand, got %d", inHand)
	}

	if inAir != 0 {
		t.Errorf("Expected 0 balls in air, got %d", inAir)
	}

	if len(balls) != 0 {
		t.Errorf("Expected 0 balls in details, got %d", len(balls))
	}
}

func TestJugglerStopAndIsRunning(t *testing.T) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2)

	if !j.IsRunning() {
		t.Error("Expected juggler to be ready to run after reset")
	}

	j.Stop()
	if !j.IsFinished() {
		t.Error("Expected juggler to be finished after stop")
	}

	if j.IsRunning() {
		t.Error("Expected juggler to not be running after stop")
	}
}

func TestJugglerAllBallsInHand(t *testing.T) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2)

	if !j.AllBallsInHand() {
		t.Error("Expected all balls to be in hand initially")
	}

	ctx := context.Background()
	eg := &errgroup.Group{}
	j.ThrowBall(ctx, eg)

	if j.AllBallsInHand() {
		t.Error("Expected not all balls to be in hand after throwing")
	}
}

func TestJugglerTimeChecks(t *testing.T) {
	j := juggler.NewJuggler(3, 0) 
	j.Reset(3, 0)

	if !j.IsJugglingTimeOver() {
		t.Error("Expected juggling time to be over for 0 minute duration")
	}

	j.Reset(3, 10) 
	if j.IsJugglingTimeOver() {
		t.Error("Expected juggling time to not be over for 10 minute duration")
	}
}

func TestBallStates(t *testing.T) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2)

	_, _, balls := j.GetStats()

	for _, ball := range balls {
		if ball.Status != "in_hand" {
			t.Errorf("Expected ball %d to be in hand, got status %s", ball.ID, ball.Status)
		}
		if ball.ID < 1 {
			t.Errorf("Expected ball ID to be positive, got %d", ball.ID)
		}
	}

	ctx := context.Background()
	eg := &errgroup.Group{}
	j.ThrowBall(ctx, eg)

	_, _, balls = j.GetStats()
	flightBallCount := 0
	handBallCount := 0

	for _, ball := range balls {
		switch ball.Status {
		case "in_hand":
			handBallCount++
		case "in_flight":
			flightBallCount++
			if ball.Elapsed < 0 {
				t.Errorf("Expected ball elapsed time to be non-negative, got %d", ball.Elapsed)
			}
		}
	}

	if handBallCount != 2 {
		t.Errorf("Expected 2 balls in hand, got %d", handBallCount)
	}

	if flightBallCount != 1 {
		t.Errorf("Expected 1 ball in flight, got %d", flightBallCount)
	}
}

func TestJugglerBallFlightAndCatch(t *testing.T) {
	j := juggler.NewJuggler(1, 2)
	j.Reset(1, 2) 

	ctx := context.Background()
	eg := &errgroup.Group{}

	success := j.ThrowBall(ctx, eg)
	if !success {
		t.Error("Expected ThrowBall to succeed")
	}

	inHand, inAir, balls := j.GetStats()
	if inHand != 0 || inAir != 1 {
		t.Errorf("Expected 0 in hand, 1 in air, got %d in hand, %d in air", inHand, inAir)
	}

	var flightBall *juggler.Ball
	for _, ball := range balls {
		if ball.Status == "in_flight" {
			flightBall = &ball
			break
		}
	}

	if flightBall == nil {
		t.Fatal("Expected to find a ball in flight")
	}

	maxWaitTime := 12 * time.Second
	checkInterval := 100 * time.Millisecond
	startTime := time.Now()

	for time.Since(startTime) < maxWaitTime {
		time.Sleep(checkInterval)

		inHand, inAir, currentBalls := j.GetStats()

		if inHand == 1 && inAir == 0 {
			for _, ball := range currentBalls {
				if ball.ID == flightBall.ID && ball.Status == "in_hand" {
					if ball.Elapsed != 0 {
						t.Errorf("Expected caught ball elapsed time to be reset to 0, got %d", ball.Elapsed)
					}
					return
				}
			}
		}
	}

	eg.Wait()

	inHand, inAir, finalBalls := j.GetStats()
	if inHand != 1 || inAir != 0 {
		t.Errorf("Expected ball to be caught after flight time, but got %d in hand, %d in air", inHand, inAir)
	}

	for _, ball := range finalBalls {
		if ball.ID == flightBall.ID {
			if ball.Status != "in_hand" {
				t.Errorf("Expected ball to be in hand after being caught, got status %s", ball.Status)
			}
			if ball.Elapsed != 0 {
				t.Errorf("Expected caught ball elapsed time to be 0, got %d", ball.Elapsed)
			}
		}
	}
}
