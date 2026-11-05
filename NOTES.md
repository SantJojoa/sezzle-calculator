# Repository structure

calculator/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/calculator/   
│   ├── internal/api/          
│   ├── go.mod
│   └── Dockerfile
├── frontend/                  
│   └── Dockerfile
├── docker-compose.yml
├── README.md
└── PROMPTS.md

# API Design
## Separate Endpoints to can tested independently and without server

POST /api/v1/add
POST /api/v1/subtract
POST /api/v1/multiply
POST /api/v1/divide
POST /api/v1/power
POST /api/v1/sqrt
POST /api/v1/percentage
GET /health

