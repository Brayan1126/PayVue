# PayVue

Aplicación web de finanzas personales: registra ingresos, deudas y pagos (con recibo adjunto) y muestra un dashboard con el resumen y gráficas.

| Parte | Tecnología | Carpeta |
|-------|------------|---------|
| API | Go + chi + SQLite | [backend/](backend/) |
| Interfaz | React 19 (Create React App) + Chart.js | [frontend/](frontend/) |
| Wireframes | SVG listos para importar en Figma | [wireframes/](wireframes/) |

```
┌──────────────┐   HTTP + header X-User-ID   ┌──────────────────┐     ┌──────────┐
│  Frontend    │ ──────────────────────────▶ │  API (Go, :8080) │ ──▶ │  SQLite  │
│  React :3000 │                             │  cmd/server      │     │ payvue.db│
└──────────────┘                             └──────────────────┘     └──────────┘
                                                      │
                                                      ▼
                                               ./uploads (recibos)
```

---

## Inicio rápido

### Opción A: Docker (todo junto)

Requiere Docker Desktop.

```bash
docker compose up -d --build
```

- Frontend: http://localhost:3000
- API: http://localhost:8080 (health check en `/health`)

La base de datos queda en `backend/data/` y los recibos en `backend/uploads/`.

### Opción B: Sin Docker

Requiere Go 1.21+, un compilador de C (GCC/MinGW, lo necesita el driver de SQLite) y Node 18+.

```bash
# Terminal 1: API
cd backend
CGO_ENABLED=1 go run ./cmd/server

# Terminal 2: frontend
cd frontend
npm install
npm start
```

### Opción C: Solo frontend, sin backend (demo / capturas)

El frontend trae un **modo mock** con datos de ejemplo en memoria. Sirve para ver la app o tomar capturas sin instalar Go ni Docker.

```bash
cd frontend
echo REACT_APP_MOCK_API=true > .env.local
npm install
npm start
```

Inicia sesión con cualquier email y contraseña. Para volver a usar la API real, borra `.env.local` o pon la variable en `false`. Más detalle en [frontend/README.md](frontend/README.md).

---

## Funcionalidades

- **Autenticación:** registro y login con email y contraseña (bcrypt).
- **Ingresos:** crear, listar, editar y eliminar.
- **Deudas:** monto total, monto restante, cuota, número de cuotas, día de pago, fecha límite y tasa de interés.
- **Pagos:** se asocian a una deuda y se puede adjuntar un recibo (imagen o PDF, máximo 10 MB).
- **Dashboard:** ingresos, ahorro y gastos del mes, gráfica de ingresos, deudas restantes, próximas a vencer y gastos por mes.
- **Historial:** tablas filtrables de deudas, ingresos y pagos, con edición, eliminación y visor de recibos.

## Limitaciones conocidas

- **No hay autenticación real en la API.** El usuario se identifica solo con el header `X-User-ID`, que el frontend toma de `localStorage`. No hay tokens ni sesiones: cualquiera que conozca un ID puede leer o modificar los datos de ese usuario. Si no se envía el header, los `GET` devuelven los registros de **todos** los usuarios.
- "Recuperar contraseña", "Cambiar contraseña" y "Editar perfil" existen solo en el frontend; la API no tiene endpoints para ellos.
- No hay tests automatizados.

## Documentación adicional

- [backend/README.md](backend/README.md): API, endpoints, configuración y estructura.
- [frontend/README.md](frontend/README.md): pantallas, modo mock y despliegue.
- [backend/docs/](backend/docs/): documentación académica (arquitectura, patrones de diseño, SOLID, pruebas, despliegue).
- [CHANGELOG.md](CHANGELOG.md)
