#!/bin/bash

cd web && REACT_APP_API_URL=https://octagon.beer npm run build
cd ..
PORT=6001 go run api/main.go
