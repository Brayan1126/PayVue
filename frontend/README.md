# PayVue Frontend

Interfaz web de PayVue en **React 19** (Create React App). Consume la API de [backend/](../backend/) para gestionar ingresos, deudas y pagos.

**Librerías principales:** react-router-dom 7 · axios · Chart.js + react-chartjs-2 · Bootstrap 5 + bootstrap-icons · MUI · SweetAlert2

---

## Ejecutar

Requiere Node 18+.

```bash
cd frontend
npm install
npm start          # http://localhost:3000
```

La app llama a la API en `REACT_APP_API_URL`; si no está definida, usa `http://localhost:8080`.

### Modo mock (sin backend)

Para ver la app sin levantar la API (por ejemplo, para tomar capturas), crea `frontend/.env.local`:

```env
REACT_APP_MOCK_API=true
```

y reinicia `npm start`. En este modo:

- El login y el registro aceptan cualquier email y contraseña.
- El dashboard y el historial cargan datos de ejemplo (ingresos, 3 deudas y algunos pagos).
- Crear, editar y eliminar funcionan en memoria; los cambios **se pierden al recargar** la página.
- "Ver recibo" no muestra imagen, porque los pagos de ejemplo no tienen archivo.

El mock está en [src/config/mockApi.js](src/config/mockApi.js) y se activa desde [src/config/api.js](src/config/api.js). Para volver a la API real, borra `.env.local` o pon la variable en `false`. Git ignora `.env.local`.

### Variables de entorno

| Variable | Descripción | Por defecto |
|----------|-------------|-------------|
| `REACT_APP_API_URL` | URL base de la API | `http://localhost:8080` |
| `REACT_APP_MOCK_API` | `true` usa datos falsos en lugar de la API | sin definir (API real) |

Son variables de CRA: se fijan **al compilar**, así que hay que reiniciar `npm start` o volver a hacer el build después de cambiarlas.

---

## Pantallas

| Ruta | Página | Descripción |
|------|--------|-------------|
| `/` | `Login` | Inicio de sesión |
| `/register` | `Register` | Crear cuenta |
| `/forgot-password` | `ForgotPassword` | Recuperar contraseña (solo interfaz, no hay endpoint) |
| `/dashboard` | `Dashboard` | Resumen, gráficas, deudas restantes y próximas a vencer |
| `/add-debt` | `AddDebt` | Formulario de deuda |
| `/add-income` | `AddIncome` | Formulario de ingreso |
| `/add-payment` | `AddPayment` | Formulario de pago con recibo |
| `/history` | `History` | Tablas filtrables con edición, eliminación y visor de recibos |
| `/change-password` | `ChangePassword` | Cambiar contraseña (solo interfaz, no hay endpoint) |

Todas las rutas excepto login, registro y recuperar contraseña están protegidas por `PrivateRoute`: si no hay `user` en `localStorage`, redirige a `/`.

El Dashboard y el Historial recargan los datos cada 5 segundos.

## Cómo se comunica con la API

Todas las llamadas pasan por la instancia `api` de [src/config/api.js](src/config/api.js):

1. Al hacer login, se guarda `{ user_id, email }` en `localStorage` con la clave `user`.
2. Un interceptor de axios añade el header `X-User-ID` a cada petición.
3. "Cerrar sesión" borra esa clave de `localStorage`.

No se usan tokens. Las limitaciones que esto implica están descritas en el [README del backend](../backend/README.md#identificación-del-usuario).

## Estructura

```
src/
├── App.js               # Rutas y PrivateRoute
├── index.js / index.css # Entrada y estilos globales
├── config/
│   ├── api.js           # Instancia axios, interceptor X-User-ID, helpers de usuario
│   └── mockApi.js       # Datos de ejemplo para el modo mock
├── components/
│   ├── Sidebar.js       # Menú lateral
│   ├── ProfilePanel.js  # Panel lateral de perfil
│   └── Toast.js         # Notificaciones
├── pages/               # Una página por ruta (ver tabla)
└── utils/alert.js       # Alertas con SweetAlert2
```

## Build y despliegue

```bash
npm run build      # genera build/
```

- **Docker:** el [Dockerfile](Dockerfile) compila con Node 18 y sirve `build/` con nginx ([nginx.conf](nginx.conf)) en el puerto 80. `docker-compose` lo publica en el 3000 y pasa `REACT_APP_API_URL` como argumento de build.
- **Vercel:** [vercel.json](vercel.json) compila con `@vercel/static-build` y redirige todas las rutas a `index.html`. `REACT_APP_API_URL` se toma del secreto `@api_url`.
