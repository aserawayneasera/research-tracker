research-tracker

Go + chi + SQLite REST API to track papers and submission status.

Run
  export DB_PATH=research_tracker.db
  export PORT=8080
  go run ./cmd/server

Health
  GET /health

Papers
  POST   /papers
  GET    /papers?status=submitted&year=2026&q=retina&limit=50&offset=0
  GET    /papers/{id}
  PUT    /papers/{id}
  DELETE /papers/{id}

Example create
  curl -X POST http://localhost:8080/papers \
    -H "Content-Type: application/json" \
    -d '{"title":"LGF for Small Object Detection","venue":"PRL","year":2026,"status":"submitted","tags":"retinanet,small","notes":"seed=42"}'

Build docker
  docker build -t research-tracker .
  docker run -p 8080:8080 -e DB_PATH=/app/research_tracker.db research-tracker
