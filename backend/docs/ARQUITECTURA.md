# Arquitectura del Software - PayVue

## 1. Visión General

PayVue es una aplicación de gestión financiera personal construida con una **arquitectura de microservicios** que implementa el patrón **CQRS (Command Query Responsibility Segregation)**.

## 2. Arquitectura Implementada

### 2.1 Arquitectura de Capas + Microservicios

```
┌─────────────────────────────────────────────────────────────────┐
│                        FRONTEND (React)                         │
│                    Puerto: 3000 / Vercel                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API REST (Go + Chi)                        │
│                    Puerto: 8080 / Render                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  Auth Handler   │  │ Income Handler  │  │  Debt Handler   │  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  │
│           │                    │                    │           │
│  ┌────────▼────────────────────▼────────────────────▼────────┐  │
│  │                    SERVICE LAYER                          │  │
│  │         (Lógica de negocio y validaciones)                │  │
│  └────────────────────────────┬──────────────────────────────┘  │
│                               │                                 │
│  ┌────────────────────────────▼──────────────────────────────┐  │
│  │                   REPOSITORY LAYER                        │  │
│  │              (Acceso a datos - SQLite)                    │  │
│  └────────────────────────────┬──────────────────────────────┘  │
└───────────────────────────────┼─────────────────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │      SQLite DB        │
                    │    (payvue.db)        │
                    └───────────────────────┘
```

### 2.2 Patrón CQRS (Command Query Responsibility Segregation)

El sistema separa las operaciones de **lectura (Query)** de las operaciones de **escritura (Command)**:

```
┌──────────────────────────────────────────────────────────────┐
│                      CLIENTE (Frontend)                       │
└──────────────────────┬───────────────────┬───────────────────┘
                       │                   │
            ┌──────────▼──────────┐ ┌──────▼──────────┐
            │    READER SERVICE   │ │  WRITER SERVICE │
            │   (GET Operations)  │ │ (POST/PUT/DEL)  │
            │   Puerto: 8080      │ │  Puerto: 8081   │
            └──────────┬──────────┘ └──────┬──────────┘
                       │                   │
                       └─────────┬─────────┘
                                 │
                    ┌────────────▼────────────┐
                    │      BASE DE DATOS      │
                    │        (SQLite)         │
                    └─────────────────────────┘
```

**Beneficios de CQRS implementados:**
- **Escalabilidad**: Los servicios de lectura y escritura pueden escalar independientemente
- **Optimización**: Cada servicio puede optimizarse para su tipo de operación
- **Mantenibilidad**: Código más limpio y separación de responsabilidades
- **Flexibilidad**: Posibilidad de usar diferentes modelos de datos para lectura y escritura

### 2.3 Estructura de Carpetas (Clean Architecture)

```
payvue_proyecto_software/
├── cmd/                           # Puntos de entrada de la aplicación
│   ├── app/                       # Configuración compartida
│   │   ├── config/                # Carga de configuración
│   │   └── container/             # Inyección de dependencias
│   ├── reader/                    # Microservicio de lectura
│   │   └── main.go
│   ├── writer/                    # Microservicio de escritura
│   │   └── main.go
│   └── server/                    # Servidor unificado
│       └── main.go
├── pkg/                           # Código reutilizable
│   ├── domain/                    # Capa de dominio (entidades y lógica de negocio)
│   │   ├── debt/
│   │   ├── income/
│   │   ├── payment/
│   │   └── user/
│   ├── repository/                # Capa de acceso a datos
│   │   ├── database/
│   │   ├── debt/
│   │   ├── income/
│   │   ├── payment/
│   │   └── user/
│   ├── rest/                      # Capa de presentación (API REST)
│   │   ├── entities/              # DTOs
│   │   ├── reader/                # Handlers de lectura
│   │   └── writer/                # Handlers de escritura
│   └── utils/                     # Utilidades compartidas
├── frontend/                      # Aplicación React
└── docs/                          # Documentación
```

## 3. Componentes del Sistema

### 3.1 Backend (Go)

| Componente | Responsabilidad |
|------------|-----------------|
| **Config** | Carga de variables de entorno y configuración |
| **Container** | Inyección de dependencias (DI Container) |
| **Domain** | Entidades, interfaces y lógica de negocio |
| **Repository** | Implementación de acceso a datos |
| **REST Handlers** | Manejo de peticiones HTTP |
| **Middleware** | Logger, Recoverer, CORS, Timeout |

### 3.2 Frontend (React)

| Componente | Responsabilidad |
|------------|-----------------|
| **Pages** | Vistas principales (Login, Dashboard, Records) |
| **Components** | Componentes reutilizables (Sidebar, Forms) |
| **Config** | Configuración de API URL |
| **Utils** | Funciones auxiliares (alerts) |

## 4. Flujo de Datos

### 4.1 Ejemplo: Crear una Deuda

```
1. Usuario llena formulario en Frontend
                    │
                    ▼
2. Frontend envía POST /finances/debt
                    │
                    ▼
3. Router (Chi) recibe la petición
                    │
                    ▼
4. Middleware procesa (Logger, Auth, CORS)
                    │
                    ▼
5. Handler decodifica JSON y valida
                    │
                    ▼
6. Service aplica lógica de negocio
                    │
                    ▼
7. Repository persiste en SQLite
                    │
                    ▼
8. Respuesta JSON al Frontend
```

## 5. Tecnologías Utilizadas

### Backend
- **Lenguaje**: Go 1.21
- **Framework HTTP**: Chi Router
- **Base de Datos**: SQLite
- **Validación**: go-playground/validator
- **Variables de Entorno**: godotenv
- **Seguridad**: bcrypt para passwords

### Frontend
- **Framework**: React 19
- **UI Library**: Material-UI, Bootstrap
- **HTTP Client**: Axios
- **Routing**: React Router DOM
- **Charts**: Chart.js

### DevOps
- **Contenedores**: Docker
- **Orquestación**: Docker Compose
- **CI/CD**: GitHub Actions
- **Frontend Hosting**: Vercel
- **Backend Hosting**: Render

