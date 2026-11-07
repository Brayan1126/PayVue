# Patrones de Diseño Implementados - PayVue

## 1. Repository Pattern

### Descripción
El patrón Repository actúa como una capa de abstracción entre la lógica de negocio y el acceso a datos, proporcionando una interfaz uniforme para operaciones CRUD.

### Implementación en PayVue

```go
// pkg/domain/debt/container.go - Interface del Repository
type Repository interface {
    CreateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    GetAllDebts(ctx context.Context) ([]Debt, error)
    GetDebtByID(ctx context.Context, id int) (*Debt, error)
    UpdateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    DeleteDebt(ctx context.Context, id int) error
}

// pkg/repository/debt/repository.go - Implementación concreta
type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) debt.Repository {
    return &repository{db: db}
}

func (r *repository) CreateDebt(ctx context.Context, d *debt.Debt) (*debt.Debt, error) {
    query := `INSERT INTO debts (...) VALUES (...)`
    result, err := r.db.ExecContext(ctx, query, ...)
    // ...
}
```

### Beneficios
- ✅ Desacopla la lógica de negocio del acceso a datos
- ✅ Facilita pruebas unitarias con mocks
- ✅ Permite cambiar la base de datos sin afectar la lógica de negocio
- ✅ Centraliza las consultas SQL

---

## 2. Dependency Injection (DI) / Container Pattern

### Descripción
Patrón que permite inyectar dependencias en lugar de crearlas internamente, facilitando la inversión de control y las pruebas.

### Implementación en PayVue

```go
// cmd/app/container/container.go
type Container struct {
    DebtService    debt.Service
    IncomeService  income.Service
    PaymentService payment.Service
    UserService    user.Service
    DB             *sql.DB
}

func New(cfg config.Config) *Container {
    // Inicializar base de datos
    db, err := database.InitDB(cfg.DatabasePath)
    
    // Crear repositorios
    debtRepository := debtRepo.NewRepository(db)
    
    // Inyectar repositorio en el contenedor del dominio
    debtContainer := &debt.Container{
        Repository: debtRepository,
    }
    
    // Crear servicio con dependencias inyectadas
    debtService := debt.New(debtContainer)
    
    return &Container{
        DebtService: debtService,
        // ...
    }
}
```

### Uso en los Handlers

```go
// cmd/server/main.go
func main() {
    cfg := config.Get()
    
    // El container crea todas las dependencias
    globalContainer := container.New(cfg)
    defer globalContainer.Close()
    
    // Inyectar servicios en los handlers
    debtHandler := writerDebt.NewHandler(globalContainer.DebtService)
    // ...
}
```

### Beneficios
- ✅ Facilita las pruebas unitarias (se pueden inyectar mocks)
- ✅ Reduce el acoplamiento entre componentes
- ✅ Centraliza la creación de objetos
- ✅ Facilita el mantenimiento y la extensibilidad

---

## 3. Factory Pattern

### Descripción
El patrón Factory proporciona una interfaz para crear objetos sin especificar sus clases concretas.

### Implementación en PayVue

```go
// pkg/domain/debt/service.go
func New(container *Container) Service {
    return &service{
        Container: container,
    }
}

// pkg/repository/debt/repository.go
func NewRepository(db *sql.DB) debt.Repository {
    return &repository{
        db: db,
    }
}

// pkg/rest/writer/debt/handler.go
func NewHandler(debtService debt.Service) *handler {
    return &handler{
        debtService: debtService,
    }
}
```

### Beneficios
- ✅ Encapsula la lógica de creación de objetos
- ✅ Permite crear objetos sin conocer su implementación concreta
- ✅ Facilita la adición de nuevos tipos sin modificar el código cliente

---

## 4. Strategy Pattern (implícito en Service Layer)

### Descripción
Define una familia de algoritmos encapsulados e intercambiables.

### Implementación en PayVue

```go
// pkg/domain/debt/service.go - Interface define el contrato
type Service interface {
    CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error)
    GetAllDebts(ctx context.Context) ([]Debt, error)
    GetDebtByID(ctx context.Context, id int) (*Debt, error)
    UpdateDebt(ctx context.Context, id int, request UpdateDebtRequest) (*Debt, error)
    DeleteDebt(ctx context.Context, id int) error
}

// Implementación concreta puede variar
type service struct {
    *Container
}

func (s *service) CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error) {
    // Estrategia específica de creación
    dueDate, err := time.Parse("2006-01-02", request.DueDate)
    // Validaciones y lógica de negocio
    debt := &Debt{
        Name: request.Name,
        // ...
    }
    return s.Repository.CreateDebt(ctx, debt)
}
```

---

## 5. Middleware Pattern (Chain of Responsibility)

### Descripción
Permite pasar peticiones a través de una cadena de handlers, donde cada uno puede procesar o pasar la petición al siguiente.

### Implementación en PayVue

```go
// cmd/server/main.go
router := chi.NewRouter()

// Cadena de middlewares
router.Use(middleware.Logger)      // 1. Logging
router.Use(middleware.Recoverer)   // 2. Recuperación de panics
router.Use(middleware.RequestID)   // 3. ID único por request
router.Use(middleware.Timeout(60 * time.Second)) // 4. Timeout

// CORS Middleware
router.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           300,
}))
```

### Flujo de la Cadena

```
Request → Logger → Recoverer → RequestID → Timeout → CORS → Handler → Response
```

---

## 6. DTO Pattern (Data Transfer Object)

### Descripción
Objetos que transportan datos entre capas de la aplicación.

### Implementación en PayVue

```go
// pkg/rest/entities/debt.go - Request DTOs
type CreateDebtRequest struct {
    Name              string  `json:"name" validate:"required"`
    TotalAmount       float64 `json:"total_amount" validate:"required,gt=0"`
    // ...
}

// pkg/domain/debt/entities.go - Response DTOs
type DebtResponse struct {
    ID                int     `json:"id"`
    Name              string  `json:"name"`
    TotalAmount       float64 `json:"total_amount"`
    // ...
}

// Mapper: convierte entre entidades y DTOs
func ToDebtResponse(debt *Debt) DebtResponse {
    return DebtResponse{
        ID:          debt.ID,
        Name:        debt.Name,
        TotalAmount: debt.TotalAmount,
        // ...
    }
}
```

---

## 7. Singleton Pattern (Configuration)

### Descripción
Garantiza una única instancia de configuración en toda la aplicación.

### Implementación en PayVue

```go
// cmd/app/config/config.go
func init() {
    // Se ejecuta una sola vez al importar el paquete
    envPath := findEnvFile()
    err := godotenv.Load(envPath)
    // ...
}

func Get() Config {
    // Siempre retorna la misma configuración basada en variables de entorno
    return Config{
        Port:         getEnv("PORT", "8080"),
        DatabasePath: getEnv("DATABASE_PATH", "./payvue.db"),
        Environment:  getEnv("ENVIRONMENT", "development"),
        // ...
    }
}
```

---

## 8. Builder Pattern (Request Validation)

### Descripción
Construye objetos complejos paso a paso con validación.

### Implementación en PayVue

```go
// Uso de validator para construir y validar requests
var validate = validator.New()

func (h *handler) CreateDebt(w http.ResponseWriter, r *http.Request) {
    var request entities.CreateDebtRequest
    
    // Paso 1: Decodificar JSON
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
        return
    }
    
    // Paso 2: Validar estructura
    if err := validate.Struct(request); err != nil {
        respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
        return
    }
    
    // Paso 3: Convertir a dominio y crear
    debt, err := h.debtService.CreateDebt(ctx, request.ToDomain())
    // ...
}
```

---

## Resumen de Patrones

| Patrón | Ubicación | Propósito |
|--------|-----------|-----------|
| Repository | `pkg/repository/*` | Abstracción de acceso a datos |
| Dependency Injection | `cmd/app/container/` | Gestión de dependencias |
| Factory | `New*()` functions | Creación de objetos |
| Strategy | `pkg/domain/*/service.go` | Lógica de negocio intercambiable |
| Middleware | `cmd/*/main.go` | Procesamiento en cadena |
| DTO | `pkg/rest/entities/` | Transferencia de datos |
| Singleton | `cmd/app/config/` | Configuración única |
| Builder | Handlers con validator | Construcción y validación |

