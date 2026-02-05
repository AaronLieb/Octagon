#!/bin/bash

# Start the Go API server in the background
echo "Starting Go API server..."
PORT=6001 go run api/main.go &
API_PID=$!

# Start the React development server
echo "Starting React development server..."
cd web && REACT_APP_API_URL=https://octagon.beer PORT=6002 npm start &
WEB_PID=$!

# Function to cleanup processes on exit
cleanup() {
  echo "Stopping servers..."
  kill $API_PID 2>/dev/null
  kill $WEB_PID 2>/dev/null
  exit
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM

echo "Both servers started. Press Ctrl+C to stop."

# Wait for both processes
wait
