package juggler

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Ball represents a juggling ball
type Ball struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`      // "in_hand", "in_flight", "dropped"
	FlightTime int       `json:"flight_time"` // seconds
	Elapsed    int       `json:"elapsed"`     // seconds elapsed in flight
	StartTime  time.Time `json:"start_time"`
}

// Juggler manages the juggling process
type Juggler struct {
	balls        map[int]*Ball
	ballsInHand  []int
	ballsInAir   []int
	nextBallID   int
	mu           sync.RWMutex
	totalBalls   int
	jugglingTime time.Duration
	startTime    time.Time
	finished     bool
}

// NewJuggler creates a new juggler
func NewJuggler(totalBalls int, jugglingTimeMinutes int) *Juggler {
	j := &Juggler{
		balls:        make(map[int]*Ball),
		ballsInHand:  make([]int, 0),
		ballsInAir:   make([]int, 0),
		nextBallID:   1,
		totalBalls:   totalBalls,
		jugglingTime: time.Duration(jugglingTimeMinutes) * time.Minute,
		startTime:    time.Now(),
		finished:     true, // Start as finished/not running
	}

	if totalBalls > 0 {
		for i := 0; i < totalBalls; i++ {
			ball := &Ball{
				ID:     j.nextBallID,
				Status: "in_hand",
			}
			j.balls[j.nextBallID] = ball
			j.ballsInHand = append(j.ballsInHand, j.nextBallID)
			j.nextBallID++
		}
	}

	return j
}

// GetStats returns current juggling statistics
func (j *Juggler) GetStats() (inHand, inAir int, ballDetails []Ball) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	inHand = len(j.ballsInHand)
	inAir = len(j.ballsInAir)

	ballDetails = make([]Ball, 0, len(j.balls))
	for _, ball := range j.balls {
		ballDetails = append(ballDetails, *ball)
	}

	return inHand, inAir, ballDetails
}

// ThrowBall throws a ball into the air
func (j *Juggler) ThrowBall(ctx context.Context, eg *errgroup.Group) bool {
	j.mu.Lock()
	defer j.mu.Unlock()

	if len(j.ballsInHand) == 0 {
		return false
	}

	ballID := j.ballsInHand[0]
	j.ballsInHand = j.ballsInHand[1:]

	j.ballsInAir = append(j.ballsInAir, ballID)

	ball := j.balls[ballID]
	ball.Status = "in_flight"
	ball.FlightTime = rand.Intn(6) + 5
	ball.Elapsed = 0
	ball.StartTime = time.Now()

	eg.Go(func() error {
		return j.flyBall(ctx, ballID)
	})

	return true
}

// flyBall simulates a ball flying in the air
func (j *Juggler) flyBall(ctx context.Context, ballID int) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			j.mu.Lock()
			ball := j.balls[ballID]
			ball.Elapsed++

			fmt.Printf("Ball %d: %d/%d seconds\n", ballID, ball.Elapsed, ball.FlightTime)

			if ball.Elapsed >= ball.FlightTime {
				j.catchBall(ballID)
				j.mu.Unlock()
				return nil
			}
			j.mu.Unlock()
		}
	}
}

// catchBall catches a ball and puts it back in hand
func (j *Juggler) catchBall(ballID int) {
	for i, id := range j.ballsInAir {
		if id == ballID {
			j.ballsInAir = append(j.ballsInAir[:i], j.ballsInAir[i+1:]...)
			break
		}
	}

	j.ballsInHand = append(j.ballsInHand, ballID)

	ball := j.balls[ballID]
	ball.Status = "in_hand"
	ball.Elapsed = 0
}

// IsJugglingTimeOver checks if juggling time is over
func (j *Juggler) IsJugglingTimeOver() bool {
	return time.Since(j.startTime) >= j.jugglingTime
}

// AllBallsInHand checks if all balls are in hand (none in air)
func (j *Juggler) AllBallsInHand() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return len(j.ballsInAir) == 0
}

// PrintStats prints current juggling statistics
func (j *Juggler) PrintStats() {
	j.mu.RLock()
	defer j.mu.RUnlock()

	fmt.Printf("\n=== Juggling State ===\n")
	fmt.Printf("Elapsed Time: %.0f seconds\n", time.Since(j.startTime).Seconds())
	fmt.Printf("Balls in Hand: %d\n", len(j.ballsInHand))
	fmt.Printf("Balls in Air: %d\n", len(j.ballsInAir))
	fmt.Printf("Ball Details:\n")

	for _, ball := range j.balls {
		status := ball.Status
		switch status {
		case "in_flight":
			status = fmt.Sprintf("in flight (%d/%d sec)", ball.Elapsed, ball.FlightTime)
		case "in_hand":
			status = "in hand"
		}
		fmt.Printf("  Ball %d: %s\n", ball.ID, status)
	}
	fmt.Printf("========================\n\n")
}

// SetFinished marks juggling as finished
func (j *Juggler) SetFinished() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.finished = true
}

// IsFinished checks if juggling is finished
func (j *Juggler) IsFinished() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.finished
}

// GetStartTime returns the start time of juggling
func (j *Juggler) GetStartTime() time.Time {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.startTime
}

// GetTotalBalls returns the total number of balls
func (j *Juggler) GetTotalBalls() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.totalBalls
}

// GetJugglingTime returns the juggling time duration
func (j *Juggler) GetJugglingTime() time.Duration {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.jugglingTime
}

// Reset resets the juggler to initial state with new configuration
func (j *Juggler) Reset(totalBalls int, jugglingTimeMinutes int) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.balls = make(map[int]*Ball)
	j.ballsInHand = make([]int, 0)
	j.ballsInAir = make([]int, 0)
	j.nextBallID = 1
	j.totalBalls = totalBalls
	j.jugglingTime = time.Duration(jugglingTimeMinutes) * time.Minute
	j.startTime = time.Now()
	j.finished = false

	for i := 0; i < totalBalls; i++ {
		ball := &Ball{
			ID:     j.nextBallID,
			Status: "in_hand",
		}
		j.balls[j.nextBallID] = ball
		j.ballsInHand = append(j.ballsInHand, j.nextBallID)
		j.nextBallID++
	}
}

// Start starts the juggling simulation
func (j *Juggler) Start() {
	go func() {
		ctx := context.Background()
		eg := &errgroup.Group{}

		throwTicker := time.NewTicker(time.Millisecond * 500) 
		defer throwTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-throwTicker.C:
				if j.IsFinished() || j.IsJugglingTimeOver() {
					j.SetFinished()
					return
				}

				thrownCount := 0
				for {
					if j.ThrowBall(ctx, eg) {
						thrownCount++
					} else {
						break 
					}
				}

				if thrownCount > 0 {
					fmt.Printf("Threw %d ball(s)! Time: %.0f seconds\n", thrownCount, time.Since(j.GetStartTime()).Seconds())
				}
			}
		}
	}()
}

// IsRunning checks if the juggler is currently running
func (j *Juggler) IsRunning() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return !j.finished && time.Since(j.startTime) < j.jugglingTime
}

// Stop stops the juggling process
func (j *Juggler) Stop() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.finished = true
}
