# Principios SOLID Aplicados - PayVue

## Introducción

Los principios SOLID son un conjunto de cinco principios de diseño orientado a objetos que promueven código mantenible, escalable y testeable. A continuación se detalla cómo cada principio se aplica en PayVue.

---

## S - Single Responsibility Principle (SRP)
### "Una clase debe tener una única razón para cambiar"

### Implementación en PayVue

Cada componente tiene una única responsabilidad:

```go
// ✅ CORRECTO: Handler solo maneja HTTP
// pkg/rest/writer/debt/debt_handlers.go
type handler struct {
    debtService debt.Service  // Delega lógica al servicio
}

func (h *handler) CreateDebt(w http.ResponseWriter, r *http.Request) {
    // Solo se encarga de:
    // 1. Parsear request
    // 2. Validar input
    // 3. Llamar al servicio
    // 4. Formatear response
}

// ✅ CORRECTO: Service solo contiene lógica de negocio
// pkg/domain/debt/service.go
type service struct {
    *Container
}

func (s *service) CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error) {
    // Solo se encarga de:
    // 1. Validaciones de negocio
    // 2. Transformaciones de datos
    // 3. Orquestación de operaciones
}

// ✅ CORRECTO: Repository solo accede a datos
// pkg/repository/debt/repository.go
type repository struct {
    db *sql.DB
}

func (r *repository) CreateDebt(ctx context.Context, d *debt.Debt) (*Debt, error) {
    // Solo se encarga de:
    // 1. Ejecutar queries SQL
    // 2. Mapear resultados
}
```

### Estructura de Responsabilidades

| Capa | Responsabilidad Única |
|------|----------------------|
| **Handler** | Manejar peticiones HTTP y respuestas |
| **Service** | Lógica de negocio y validaciones |
| **Repository** | Acceso y persistencia de datos |
| **Entity** | Representar datos del dominio |
| **Config** | Cargar configuración |
| **Container** | Inyectar dependencias |

---

## O - Open/Closed Principle (OCP)
### "Abierto para extensión, cerrado para modificación"

### Implementación en PayVue

Uso de interfaces para permitir extensión sin modificar código existente:

```go
// pkg/domain/debt/container.go
// La interfaz está CERRADA para modificación
type Repository interface {
    CreateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    GetAllDebts(ctx context.Context) ([]Debt, error)
    GetDebtByID(ctx context.Context, id int) (*Debt, error)
    UpdateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    DeleteDebt(ctx context.Context, id int) error
}

// pkg/repository/debt/repository.go
// Implementación SQLite - ABIERTA para extensión
type sqliteRepository struct {
    db *sql.DB
}

// Se puede agregar una nueva implementación sin modificar la existente
// Por ejemplo, para PostgreSQL:
type postgresRepository struct {
    db *sql.DB
}

func (r *postgresRepository) CreateDebt(ctx context.Context, d *debt.Debt) (*Debt, error) {
    // Implementación para PostgreSQL
}
```

### Ejemplo de Extensión: Agregar Cache

```go
// Nueva implementación con cache - NO modifica el código existente
type cachedRepository struct {
    inner debt.Repository  // Decorador
    cache map[int]*debt.Debt
}

func (r *cachedRepository) GetDebtByID(ctx context.Context, id int) (*debt.Debt, error) {
    // Primero buscar en cache
    if cached, ok := r.cache[id]; ok {
        return cached, nil
    }
    // Si no está, buscar en el repositorio real
    return r.inner.GetDebtByID(ctx, id)
}
```

---

## L - Liskov Substitution Principle (LSP)
### "Los objetos de una superclase deben poder reemplazarse por objetos de sus subclases"

### Implementación en PayVue

Todas las implementaciones de interfaces son intercambiables:

```go
// pkg/domain/debt/service.go
type Service interface {
    CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error)
    GetAllDebts(ctx context.Context) ([]Debt, error)
    // ...
}

// Implementación estándar
type standardService struct {
    *Container
}

// El handler acepta cualquier implementación de Service
// pkg/rest/writer/debt/handler.go
type handler struct {
    debtService debt.Service  // Interface, no implementación concreta
}

func NewHandler(debtService debt.Service) *handler {
    return &handler{
        debtService: debtService,
    }
}

// En pruebas, se puede usar un mock que implementa la misma interfaz
type mockDebtService struct{}

func (m *mockDebtService) CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error) {
    // Mock implementation for testing
    return &Debt{ID: 1, Name: "Test"}, nil
}
```

### Demostración de Sustitución

```go
// Ambas implementaciones son intercambiables
var service debt.Service

// Producción
service = debt.New(productionContainer)

// Testing
service = &mockDebtService{}

// El handler funciona igual con ambas
handler := NewHandler(service)
```

---

## I - Interface Segregation Principle (ISP)
### "Los clientes no deben verse obligados a depender de interfaces que no utilizan"

### Implementación en PayVue

Interfaces pequeñas y específicas:

```go
// ✅ CORRECTO: Interfaces segregadas por responsabilidad

// pkg/domain/debt/container.go - Solo operaciones de Debt
type Repository interface {
    CreateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    GetAllDebts(ctx context.Context) ([]Debt, error)
    GetDebtByID(ctx context.Context, id int) (*Debt, error)
    UpdateDebt(ctx context.Context, debt *Debt) (*Debt, error)
    DeleteDebt(ctx context.Context, id int) error
}

// pkg/domain/income/container.go - Solo operaciones de Income
type Repository interface {
    CreateIncome(ctx context.Context, income *Income) (*Income, error)
    GetAllIncomes(ctx context.Context) ([]Income, error)
    GetIncomeByID(ctx context.Context, id int) (*Income, error)
    UpdateIncome(ctx context.Context, income *Income) (*Income, error)
    DeleteIncome(ctx context.Context, id int) error
}

// pkg/domain/user/container.go - Solo operaciones de User
type Repository interface {
    CreateUser(ctx context.Context, user *User) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByID(ctx context.Context, id int) (*User, error)
}
```

### Comparación

```go
// ❌ INCORRECTO: Interfaz "gorda" que obliga a implementar todo
type MegaRepository interface {
    // Debt
    CreateDebt(...) error
    GetAllDebts(...) []Debt
    // Income
    CreateIncome(...) error
    GetAllIncomes(...) []Income
    // Payment
    CreatePayment(...) error
    // User
    CreateUser(...) error
    // ... 20 métodos más
}

// ✅ CORRECTO: Interfaces pequeñas y específicas (como se implementó)
type DebtRepository interface { /* solo métodos de debt */ }
type IncomeRepository interface { /* solo métodos de income */ }
type UserRepository interface { /* solo métodos de user */ }
```

---

## D - Dependency Inversion Principle (DIP)
### "Depender de abstracciones, no de implementaciones concretas"

### Implementación en PayVue

Las capas superiores dependen de interfaces, no de implementaciones:

```go
// ✅ CORRECTO: Service depende de interfaz Repository
// pkg/domain/debt/service.go
type service struct {
    *Container  // Contiene Repository interface
}

func (s *service) CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error) {
    // s.Repository es una interfaz, no una implementación concreta
    return s.Repository.CreateDebt(ctx, debt)
}

// ✅ CORRECTO: Handler depende de interfaz Service
// pkg/rest/writer/debt/handler.go
type handler struct {
    debtService debt.Service  // Interface
}

// ✅ CORRECTO: Container inyecta las dependencias
// cmd/app/container/container.go
func New(cfg config.Config) *Container {
    db, _ := database.InitDB(cfg.DatabasePath)
    
    // Crear implementación concreta del repository
    debtRepository := debtRepo.NewRepository(db)
    
    // Inyectar como interfaz en el container del dominio
    debtContainer := &debt.Container{
        Repository: debtRepository,  // Repository interface
    }
    
    // Crear servicio con la dependencia inyectada
    debtService := debt.New(debtContainer)
    
    return &Container{
        DebtService: debtService,  // Service interface
    }
}
```

### Diagrama de Dependencias

```
┌─────────────────────────────────────────────────────────────┐
│                     CAPAS SUPERIORES                        │
│  ┌─────────┐     ┌─────────┐     ┌─────────────────────┐   │
│  │ Handler │────▶│ Service │────▶│ Repository Interface│   │
│  └─────────┘     └─────────┘     └─────────────────────┘   │
│       │               │                    ▲               │
│       │               │                    │               │
│       ▼               ▼                    │               │
│  debt.Service    debt.Repository           │               │
│  (interface)     (interface)               │               │
└────────────────────────────────────────────┼───────────────┘
                                             │
┌────────────────────────────────────────────┼───────────────┐
│                     CAPAS INFERIORES       │               │
│                                            │               │
│  ┌──────────────────────────────────┐      │               │
│  │     SQLite Repository            │──────┘               │
│  │  (implementación concreta)       │                      │
│  └──────────────────────────────────┘                      │
└─────────────────────────────────────────────────────────────┘

Las flechas apuntan hacia las ABSTRACCIONES (interfaces),
no hacia las implementaciones concretas.
```

---

## Resumen de Aplicación SOLID

| Principio | Aplicación en PayVue |
|-----------|---------------------|
| **SRP** | Separación Handler/Service/Repository |
| **OCP** | Interfaces permiten nuevas implementaciones |
| **LSP** | Mocks intercambiables con implementaciones reales |
| **ISP** | Interfaces específicas por dominio |
| **DIP** | Inyección de dependencias via Container |

### Beneficios Obtenidos

1. **Mantenibilidad**: Cambios localizados sin efectos secundarios
2. **Testabilidad**: Fácil crear mocks para pruebas unitarias
3. **Escalabilidad**: Agregar nuevas funcionalidades sin modificar código existente
4. **Legibilidad**: Código organizado por responsabilidades claras
5. **Flexibilidad**: Cambiar implementaciones (ej. base de datos) sin afectar la lógica de negocio

