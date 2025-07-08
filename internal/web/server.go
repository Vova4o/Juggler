package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"juggler/internal/juggler"
)

// StatsResponse represents the JSON response for stats
type StatsResponse struct {
	InHand      int            `json:"in_hand"`
	InAir       int            `json:"in_air"`
	Balls       []juggler.Ball `json:"balls"`
	TimeElapsed int            `json:"time_elapsed"`
	IsFinished  bool           `json:"is_finished"`
	IsRunning   bool           `json:"is_running"`
	TotalBalls  int            `json:"total_balls"`
	TotalTime   int            `json:"total_time"`
}

// StartRequest represents the request to start juggling
type StartRequest struct {
	TotalBalls  int `json:"total_balls"`
	TimeMinutes int `json:"time_minutes"`
}

// Server represents the web server
type Server struct {
	juggler *juggler.Juggler
	port    int
}

// NewServer creates a new web server
func NewServer(j *juggler.Juggler, port int) *Server {
	return &Server{
		juggler: j,
		port:    port,
	}
}

// Start starts the web server
func (s *Server) Start() {
	http.HandleFunc("/", s.HandleHome)
	http.HandleFunc("/api/stats", s.HandleStats)
	http.HandleFunc("/api/start", s.HandleStart)
	http.HandleFunc("/api/stop", s.HandleStop)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Веб-сервер запущен на порту %d", s.port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// HandleHome serves the main HTML page
func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>🤹 Жонглер - Интерактивный контроль</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f0f0f0; }
        .container { max-width: 900px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        
        .controls { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; border: 2px solid #e9ecef; }
        .control-group { margin: 15px 0; display: flex; align-items: center; }
        .control-group label { display: inline-block; width: 180px; font-weight: bold; color: #495057; }
        .control-group input { padding: 10px; border: 2px solid #ced4da; border-radius: 5px; width: 120px; font-size: 16px; }
        .control-group input:focus { border-color: #007bff; outline: none; }
        
        .control-buttons { margin-top: 25px; text-align: center; }
        .btn { padding: 12px 30px; margin: 0 10px; border: none; border-radius: 6px; cursor: pointer; font-size: 16px; font-weight: bold; transition: all 0.3s; }
        .btn-start { background-color: #28a745; color: white; }
        .btn-stop { background-color: #dc3545; color: white; }
        .btn:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.2); }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; transform: none; box-shadow: none; }
        
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }
        .stat-card { background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; border: 2px solid #e9ecef; }
        .stat-number { font-size: 2.5em; font-weight: bold; color: #007bff; margin-bottom: 5px; }
        .stat-label { font-size: 14px; color: #6c757d; font-weight: bold; }
        
        .status { text-align: center; margin: 20px 0; padding: 15px; border-radius: 8px; }
        .status.running { background-color: #d4edda; border: 2px solid #c3e6cb; color: #155724; }
        .status.stopped { background-color: #f8d7da; border: 2px solid #f5c6cb; color: #721c24; }
        .status.finished { background-color: #d1ecf1; border: 2px solid #bee5eb; color: #0c5460; }
        
        .balls-container { margin: 25px 0; }
        .balls-container h3 { color: #495057; margin-bottom: 15px; }
        .ball { 
            display: inline-block; 
            margin: 8px; 
            padding: 12px 18px; 
            border-radius: 25px; 
            color: white; 
            font-weight: bold;
            min-width: 100px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .ball-in-hand { background: linear-gradient(135deg, #28a745, #20c997); }
        .ball-in-flight { background: linear-gradient(135deg, #ffc107, #fd7e14); color: #000; }
        .ball-dropped { background: linear-gradient(135deg, #dc3545, #e83e8c); }
        
        .time { font-size: 1.4em; color: #495057; text-align: center; margin: 20px 0; padding: 15px; background: #e9ecef; border-radius: 8px; }
        .progress-bar { width: 100%; height: 10px; background: #e9ecef; border-radius: 5px; margin: 10px 0; overflow: hidden; }
        .progress-fill { height: 100%; background: linear-gradient(90deg, #28a745, #20c997); transition: width 0.3s; }
        
        .message { text-align: center; margin: 15px 0; padding: 10px; border-radius: 5px; }
        .message.success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .message.error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🤹 Жонглер - Интерактивный контроль</h1>
        
        <div class="controls">
            <h3>⚙️ Настройки жонглирования</h3>
            <div class="control-group">
                <label for="balls-input">Количество мячей:</label>
                <input type="number" id="balls-input" min="1" max="10" value="3">
            </div>
            <div class="control-group">
                <label for="time-input">Время (минуты):</label>
                <input type="number" id="time-input" min="1" max="60" value="2">
            </div>
            <div class="control-buttons">
                <button class="btn btn-start" id="start-btn" onclick="startJuggling()">🚀 Начать жонглирование</button>
                <button class="btn btn-stop" id="stop-btn" onclick="stopJuggling()" disabled>🛑 Остановить</button>
            </div>
        </div>
        
        <div id="message" class="message" style="display: none;"></div>
        
        <div class="time" id="time">Время: 0 секунд</div>
        <div class="progress-bar">
            <div class="progress-fill" id="progress" style="width: 0%;"></div>
        </div>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number" id="in-hand">0</div>
                <div class="stat-label">В руках</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="in-air">0</div>
                <div class="stat-label">В воздухе</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="total-balls">0</div>
                <div class="stat-label">Всего мячей</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="total-time">0</div>
                <div class="stat-label">Время (мин)</div>
            </div>
        </div>
        
        <div class="status stopped" id="status">
            <span>⏹️ Жонглирование остановлено</span>
        </div>
        
        <div class="balls-container">
            <h3>🏀 Состояние мячей:</h3>
            <div id="balls"></div>
        </div>
    </div>

    <script>
        let isRunning = false;
        
        function showMessage(text, type = 'success') {
            const messageEl = document.getElementById('message');
            messageEl.textContent = text;
            messageEl.className = 'message ' + type;
            messageEl.style.display = 'block';
            setTimeout(() => {
                messageEl.style.display = 'none';
            }, 3000);
        }
        
        function startJuggling() {
            const balls = parseInt(document.getElementById('balls-input').value);
            const time = parseInt(document.getElementById('time-input').value);
            
            if (balls < 1 || time < 1) {
                showMessage('Количество мячей и время должны быть больше 0', 'error');
                return;
            }
            
            fetch('/api/start', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    total_balls: balls,
                    time_minutes: time
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'started') {
                    isRunning = true;
                    document.getElementById('start-btn').disabled = true;
                    document.getElementById('stop-btn').disabled = false;
                    document.getElementById('balls-input').disabled = true;
                    document.getElementById('time-input').disabled = true;
                    showMessage(data.message, 'success');
                }
            })
            .catch(error => {
                showMessage('Ошибка при запуске: ' + error, 'error');
            });
        }
        
        function stopJuggling() {
            fetch('/api/stop', {
                method: 'POST'
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'stopped') {
                    isRunning = false;
                    document.getElementById('start-btn').disabled = false;
                    document.getElementById('stop-btn').disabled = true;
                    document.getElementById('balls-input').disabled = false;
                    document.getElementById('time-input').disabled = false;
                    showMessage(data.message, 'success');
                }
            })
            .catch(error => {
                showMessage('Ошибка при остановке: ' + error, 'error');
            });
        }
        
        function updateStats() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('in-hand').textContent = data.in_hand;
                    document.getElementById('in-air').textContent = data.in_air;
                    document.getElementById('total-balls').textContent = data.total_balls;
                    document.getElementById('total-time').textContent = data.total_time;
                    document.getElementById('time').textContent = 'Время: ' + data.time_elapsed + ' секунд';
                    
                    // Update progress bar
                    const progress = data.total_time > 0 ? (data.time_elapsed / (data.total_time * 60)) * 100 : 0;
                    document.getElementById('progress').style.width = Math.min(progress, 100) + '%';
                    
                    // Update status
                    const statusElement = document.getElementById('status');
                    if (data.is_running) {
                        statusElement.className = 'status running';
                        statusElement.innerHTML = '<span>🎯 Жонглирование активно</span>';
                        isRunning = true;
                        document.getElementById('start-btn').disabled = true;
                        document.getElementById('stop-btn').disabled = false;
                        document.getElementById('balls-input').disabled = true;
                        document.getElementById('time-input').disabled = true;
                    } else if (data.is_finished) {
                        statusElement.className = 'status finished';
                        statusElement.innerHTML = '<span>✅ Жонглирование завершено</span>';
                        isRunning = false;
                        document.getElementById('start-btn').disabled = false;
                        document.getElementById('stop-btn').disabled = true;
                        document.getElementById('balls-input').disabled = false;
                        document.getElementById('time-input').disabled = false;
                    } else {
                        statusElement.className = 'status stopped';
                        statusElement.innerHTML = '<span>⏹️ Жонглирование остановлено</span>';
                        isRunning = false;
                        document.getElementById('start-btn').disabled = false;
                        document.getElementById('stop-btn').disabled = true;
                        document.getElementById('balls-input').disabled = false;
                        document.getElementById('time-input').disabled = false;
                    }
                    
                    // Update balls - maintain consistent layout
                    const ballsContainer = document.getElementById('balls');
                    
                    // Create a map of balls by ID for quick lookup
                    const ballsById = {};
                    data.balls.forEach(ball => {
                        ballsById[ball.id] = ball;
                    });
                    
                    // If total balls changed, recreate the container
                    if (ballsContainer.children.length !== data.total_balls) {
                        ballsContainer.innerHTML = '';
                        for (let i = 1; i <= data.total_balls; i++) {
                            const ballElement = document.createElement('div');
                            ballElement.className = 'ball';
                            ballElement.id = 'ball-' + i;
                            ballsContainer.appendChild(ballElement);
                        }
                    }
                    
                    // Update each ball element in place
                    for (let i = 1; i <= data.total_balls; i++) {
                        const ballElement = document.getElementById('ball-' + i);
                        const ball = ballsById[i];
                        
                        if (ball) {
                            if (ball.status === 'in_hand') {
                                ballElement.className = 'ball ball-in-hand';
                                ballElement.textContent = '🏀 Мяч ' + ball.id;
                            } else if (ball.status === 'in_flight') {
                                ballElement.className = 'ball ball-in-flight';
                                ballElement.textContent = '🚀 Мяч ' + ball.id + ' (' + ball.elapsed + '/' + ball.flight_time + 's)';
                            } else {
                                ballElement.className = 'ball ball-dropped';
                                ballElement.textContent = '💥 Мяч ' + ball.id;
                            }
                        } else {
                            // Ball doesn't exist, show as empty
                            ballElement.className = 'ball';
                            ballElement.textContent = '';
                            ballElement.style.display = 'none';
                        }
                    }
                })
                .catch(error => {
                    console.error('Ошибка при получении статистики:', error);
                });
        }
        
        // Update stats every second
        setInterval(updateStats, 1000);
        // Initial update but don't start juggling automatically
        updateStats();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// HandleStats serves the stats API endpoint
func (s *Server) HandleStats(w http.ResponseWriter, r *http.Request) {
	inHand, inAir, balls := s.juggler.GetStats()

	var timeElapsed int
	if s.juggler.IsRunning() {
		timeElapsed = int(time.Since(s.juggler.GetStartTime()).Seconds())
	} else {
		timeElapsed = 0
	}

	stats := StatsResponse{
		InHand:      inHand,
		InAir:       inAir,
		Balls:       balls,
		TimeElapsed: timeElapsed,
		IsFinished:  s.juggler.IsFinished(),
		IsRunning:   s.juggler.IsRunning(),
		TotalBalls:  s.juggler.GetTotalBalls(),
		TotalTime:   int(s.juggler.GetJugglingTime().Minutes()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// HandleStart handles requests to start juggling
func (s *Server) HandleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TotalBalls <= 0 || req.TimeMinutes <= 0 {
		http.Error(w, "Balls and time must be positive", http.StatusBadRequest)
		return
	}

	s.juggler.Reset(req.TotalBalls, req.TimeMinutes)

	s.juggler.Start()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "started",
		"message": fmt.Sprintf("Juggling started with %d balls for %d minutes", req.TotalBalls, req.TimeMinutes),
	})
}

// HandleStop handles requests to stop juggling
func (s *Server) HandleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.juggler.Stop()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "stopped",
		"message": "Juggling stopped",
	})
}
