# PayVue Backend

API REST en Go para PayVue. Expone autenticación y CRUD de ingresos, deudas y pagos sobre una base de datos SQLite.

- **Router:** [chi v5](https://github.com/go-chi/chi), con CORS abierto (`*`)
- **Base de datos:** SQLite vía `mattn/go-sqlite3` (requiere CGO y un compilador de C)
- **Validación:** go-playground/validator
- **Contraseñas:** bcrypt

> **Nota:** el servidor que se usa hoy es **`cmd/server`**, un único proceso en el puerto 8080 que atiende lecturas y escrituras. Es el que compila el `Dockerfile`. Las carpetas `cmd/reader`, `cmd/writer`, `cmd/app` y `pkg/rest/reader|writer` son de la versión anterior separada en dos servicios (puertos 8080 y 8081) y **no se usan**. Ver [Código heredado](#código-heredado).

---

## Ejecutar

### Con Docker (desde la raíz del repo)

```bash
docker compose up -d --build api
curl http://localhost:8080/health    # → OK - PayVue API Server
```

Volúmenes montados:

| Ruta en el host | Ruta en el contenedor | Contenido |
|-----------------|-----------------------|-----------|
| `backend/data/` | `/app/data` | `payvue.db` |
| `backend/uploads/` | `/app/uploads` | Recibos subidos |

### Sin Docker

Requiere **Go 1.21+** y **GCC** (en Windows: MinGW o WinLibs en el `PATH`).

```bash
cd backend
cp env.example .env        # opcional
CGO_ENABLED=1 go run ./cmd/server
```

Al arrancar, las tablas se crean automáticamente (`users`, `debts`, `incomes`, `payments`). La carpeta `./uploads` también se crea si no existe.

> El `Makefile` solo tiene targets de los servicios heredados (`run-reader` / `run-writer`). Para el servidor actual usa `go run ./cmd/server`.

### Variables de entorno

Se leen de `.env` (si existe) o del entorno del sistema.

| Variable | Descripción | Por defecto |
|----------|-------------|-------------|
| `PORT` | Puerto HTTP | `8080` |
| `DATABASE_PATH` | Ruta del archivo SQLite | `./payvue.db` |
| `ENVIRONMENT` | Nombre del entorno (solo se registra en el log) | `development` |

---

## API

Base URL: `http://localhost:8080`. Las rutas **no** tienen prefijo `/api/v1`.

### Identificación del usuario

No hay tokens. Las rutas de finanzas identifican al usuario con el header **`X-User-ID: <id>`**, o con el query param `?user_id=<id>` como alternativa. El `user_id` lo devuelven el registro y el login.

- **`GET` de listas:** con el header filtra por usuario; **sin él devuelve los registros de todos los usuarios**.
- **`POST`:** asigna el registro al usuario del header.
- **`PUT` / `DELETE` por `{id}`:** no verifican que el registro pertenezca al usuario.

### Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/auth/register` | Crear usuario |
| POST | `/auth/login` | Iniciar sesión |
| POST | `/auth/logout` | Solo devuelve un mensaje (no hay sesión que cerrar) |
| GET | `/finances/income` | Listar ingresos |
| GET | `/finances/income/{id}` | Obtener un ingreso |
| POST | `/finances/income` | Crear ingreso |
| PUT | `/finances/income/{id}` | Actualizar ingreso |
| DELETE | `/finances/income/{id}` | Eliminar ingreso |
| GET | `/finances/debt` | Listar deudas |
| GET | `/finances/debt/{id}` | Obtener una deuda |
| POST | `/finances/debt` | Crear deuda |
| PUT | `/finances/debt/{id}` | Actualizar deuda |
| DELETE | `/finances/debt/{id}` | Eliminar deuda |
| GET | `/finances/payment` | Listar pagos |
| POST | `/finances/payment` | Crear pago (multipart, con recibo opcional) |
| DELETE | `/finances/payment/{id}` | Eliminar pago |
| GET | `/finances/payment/receipt/{filename}` | Descargar un recibo |

Las listas se devuelven como un **array JSON** (no envuelto en un objeto). Los errores tienen la forma `{"error": "...", "message": "..."}`.

### Ejemplos

```bash
# Registro (password mínimo 6 caracteres)
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secreto123"}'
# → 201 {"message":"Usuario registrado exitosamente","user_id":1, ...}

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secreto123"}'
# → 200 {"message":"Inicio de sesión exitoso","user_id":1, ...}

# Crear ingreso
curl -X POST http://localhost:8080/finances/income \
  -H "Content-Type: application/json" -H "X-User-ID: 1" \
  -d '{"amount":3500000,"source":"Salario","date":"2026-11-01"}'

# Crear deuda
curl -X POST http://localhost:8080/finances/debt \
  -H "Content-Type: application/json" -H "X-User-ID: 1" \
  -d '{
    "name": "Tarjeta Visa",
    "total_amount": 2400000,
    "remaining_amount": 1600000,
    "due_date": "2027-05-15",
    "interest_rate": 2.1,
    "num_installments": 12,
    "installment_amount": 200000,
    "payment_day": 15
  }'

# Crear pago con recibo (campos de formulario: amount, debt_id, date, receipt)
curl -X POST http://localhost:8080/finances/payment \
  -H "X-User-ID: 1" \
  -F "amount=200000" -F "debt_id=1" -F "date=2026-11-07" \
  -F "receipt=@recibo.pdf"

# Listar pagos del usuario (cada uno incluye debt_name, remaining_amount y receipt_url)
curl http://localhost:8080/finances/payment -H "X-User-ID: 1"
```

### Campos

**Ingreso**: `amount` (> 0), `source`, `date`. Los tres son obligatorios.

**Deuda**: `name`, `total_amount` (> 0), `remaining_amount` (≥ 0), `due_date`, `num_installments` (> 0), `installment_amount` (> 0) y `payment_day` (1-31) son obligatorios; `interest_rate` (≥ 0) es opcional. En el `PUT` también se envía `paid` (bool). La respuesta añade `remaining_payments`.

**Pago** (multipart): `amount` (> 0), `debt_id` (> 0), `date` y `receipt` (archivo opcional, máximo 10 MB). El recibo se guarda en `./uploads` y la respuesta incluye `receipt_url` (`/finances/payment/receipt/<archivo>`).

---

## Estructura

```
backend/
├── cmd/
│   ├── server/main.go        # ★ Servidor actual: rutas + handlers HTTP
│   ├── app/config/           # Carga de .env / variables de entorno
│   ├── app/container/        # Inyección de dependencias (DB → repos → servicios)
│   ├── app/main.go           # (heredado) proxy por SCOPE
│   ├── reader/  writer/      # (heredado) servicios separados
├── pkg/
│   ├── domain/               # Entidades, servicios y mappers por módulo
│   │   ├── user/  income/  debt/  payment/
│   ├── repository/           # Acceso a SQLite por módulo
│   │   ├── database/         # Conexión y creación de tablas
│   │   └── user/  income/  debt/  payment/
│   ├── rest/
│   │   ├── entities/         # DTOs de request (usados por cmd/server)
│   │   └── reader/  writer/  # (heredado) handlers de los servicios separados
│   └── utils/fileupload/     # Guardado de recibos en ./uploads
├── data/                     # BD usada por Docker (volumen)
├── docs/                     # Documentación académica
├── Dockerfile                # Compila cmd/server
├── Makefile                  # Targets de los servicios heredados
└── env.example
```

Flujo de una petición: `cmd/server` (handler) → `pkg/domain/<módulo>` (servicio) → `pkg/repository/<módulo>` → SQLite.

## Código heredado

La primera versión separaba el backend en un servicio **reader** (GET, puerto 8080) y uno **writer** (POST/PUT/DELETE, puerto 8081), lanzados por `cmd/app` según la variable `SCOPE`. Después se unificó todo en `cmd/server`. El código anterior sigue en el repositorio pero no lo usan ni Docker ni el frontend, y sus rutas (`/api/v1/...`) no coinciden con las actuales.

## Documentación

En [docs/](docs/): [Arquitectura](docs/ARQUITECTURA.md) · [Patrones de diseño](docs/PATRONES_DISENO.md) · [Principios SOLID](docs/PRINCIPIOS_SOLID.md) · [Pruebas](docs/PRUEBAS_SOFTWARE.md) · [Despliegue](docs/DESPLIEGUE.md)
