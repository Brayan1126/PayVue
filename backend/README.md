# PayVue Backend - Go Microservices

Backend en Go para gestión financiera personal, migrado desde Python Flask siguiendo la arquitectura de microservicios.

## 🚀 Inicio Rápido con Docker (Recomendado)

### Prerequisitos
- Docker
- Docker Compose

### Levantar todos los servicios

```bash
docker-compose up -d
```

¡Eso es todo! Los servicios estarán disponibles en:
- **Reader (GET)**: http://localhost:8080
- **Writer (POST/PUT/DELETE)**: http://localhost:8081

### Comandos útiles

```bash
# Ver logs
docker compose logs -f

# Ver logs de un servicio específico
docker compose logs -f reader
docker compose logs -f writer

# Detener servicios
docker compose down

# Reconstruir imágenes
docker compose up -d --build

# Ver estado de servicios
docker compose ps
```

---

## 🏗️ Arquitectura

### Servicios

```
┌─────────────────────────────────────────────┐
│           PayVue Backend                     │
├─────────────────────────────────────────────┤
│                                              │
│  ┌──────────────┐      ┌──────────────┐    │
│  │   Reader     │      │   Writer     │    │
│  │  (Port 8080) │      │ (Port 8081)  │    │
│  │              │      │              │    │
│  │  GET /api/*  │      │ POST /api/*  │    │
│  │              │      │ PUT /api/*   │    │
│  │              │      │ DELETE /api/*│    │
│  └──────┬───────┘      └──────┬───────┘    │
│         │                     │             │
│         └──────────┬──────────┘             │
│                    │                        │
│              ┌─────▼─────┐                  │
│              │  SQLite   │                  │
│              │ Database  │                  │
│              └───────────┘                  │
└─────────────────────────────────────────────┘
```

### Módulos Implementados

- **Auth**: Registro y login de usuarios
- **Debts**: Gestión de deudas
- **Incomes**: Gestión de ingresos
- **Payments**: Gestión de pagos con subida de recibos

---

## 📡 APIs Disponibles

### Reader (GET - Puerto 8080)

```bash
# Health check
curl http://localhost:8080/health

# Debts
curl http://localhost:8080/api/v1/debts
curl http://localhost:8080/api/v1/debts/{id}

# Incomes
curl http://localhost:8080/api/v1/incomes
curl http://localhost:8080/api/v1/incomes/{id}

# Payments
curl http://localhost:8080/api/v1/payments
curl http://localhost:8080/api/v1/payments/{id}/receipt
```

### Writer (POST/PUT/DELETE - Puerto 8081)

```bash
# Registro
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass"}'

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass"}'

# Crear Deuda
curl -X POST http://localhost:8081/api/v1/debts \
  -H "Content-Type: application/json" \
  -d '{
    "creditor": "Banco Nacional",
    "amount": 5000.00,
    "interest_rate": 12.5,
    "due_date": "2025-12-31",
    "description": "Préstamo personal"
  }'

# Crear Ingreso
curl -X POST http://localhost:8081/api/v1/incomes \
  -H "Content-Type: application/json" \
  -d '{
    "source": "Salario",
    "amount": 3000.00,
    "date": "2025-10-01",
    "is_recurring": true
  }'

# Crear Pago (con archivo)
curl -X POST http://localhost:8081/api/v1/payments \
  -F "debt_id=1" \
  -F "amount=500.00" \
  -F "payment_date=2025-10-20" \
  -F "receipt=@recibo.pdf"
```

---

## 🛠️ Desarrollo Sin Docker

### Prerequisitos
- Go 1.21+
- GCC (para SQLite)
- Make (opcional)

### Configuración

1. Copiar archivo de configuración:
```bash
cp env.example .env
```

2. Editar `.env` según necesites:
```env
PORT=8080
DATABASE_PATH=./payvue.db
ENV=development
CGO_ENABLED=1
```

### Ejecutar con Make

```bash
# Reader
make run-reader

# Writer
make run-writer
```

### Ejecutar con go run

```bash
# Reader
CGO_ENABLED=1 go run cmd/reader/main.go

# Writer
CGO_ENABLED=1 PORT=8081 go run cmd/writer/main.go
```

---

## 🔧 Configuración Avanzada

### Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `SCOPE` | Servicio a ejecutar (reader/writer) | reader |
| `PORT` | Puerto del servidor | 8080 |
| `DATABASE_PATH` | Ruta a la base de datos SQLite | ./payvue.db |
| `ENV` | Entorno (development/production) | development |
| `CGO_ENABLED` | Habilitar CGO para SQLite | 1 |

### Volúmenes Docker

- `./data`: Base de datos SQLite persistente
- `./uploads`: Archivos de recibos de pagos

---

## 📂 Estructura del Proyecto

```
payvue_proyecto_software/
├── cmd/
│   ├── app/          # Configuración y container principal
│   ├── reader/       # Servicio de lectura (GET)
│   └── writer/       # Servicio de escritura (POST/PUT/DELETE)
├── pkg/
│   ├── domain/       # Lógica de negocio
│   │   ├── debt/
│   │   ├── income/
│   │   ├── payment/
│   │   └── user/
│   ├── repository/   # Capa de datos
│   │   ├── debt/
│   │   ├── income/
│   │   ├── payment/
│   │   ├── user/
│   │   └── database/
│   ├── rest/         # Capa HTTP
│   │   ├── entities/ # DTOs
│   │   ├── reader/   # Handlers GET
│   │   └── writer/   # Handlers POST/PUT/DELETE
│   └── utils/        # Utilidades
├── docker-compose.yml
├── Dockerfile
└── .env
```

---

## 🧪 Testing

```bash
# Probar Reader
curl http://localhost:8080/health

# Probar Writer
curl http://localhost:8081/health
```

