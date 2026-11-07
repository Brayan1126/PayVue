# Proceso de Despliegue - PayVue

## 1. Diagrama de Despliegue (UML)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              INTERNET                                        │
└─────────────────────────────────────────────────────────────────────────────┘
        │                                           │
        │ HTTPS                                     │ HTTPS
        ▼                                           ▼
┌───────────────────────┐                 ┌───────────────────────┐
│       VERCEL          │                 │       RENDER          │
│    (Frontend Host)    │                 │    (Backend Host)     │
│                       │                 │                       │
│  ┌─────────────────┐  │                 │  ┌─────────────────┐  │
│  │    CDN Edge     │  │                 │  │  Docker Runtime │  │
│  │    Network      │  │                 │  │    Container    │  │
│  └────────┬────────┘  │                 │  └────────┬────────┘  │
│           │           │                 │           │           │
│  ┌────────▼────────┐  │  HTTP/REST      │  ┌────────▼────────┐  │
│  │  React App      │──┼────────────────▶│  │   Go API        │  │
│  │  (Static SPA)   │  │                 │  │   Server        │  │
│  │                 │  │                 │  │                 │  │
│  │  - HTML/CSS/JS  │  │                 │  │  - Chi Router   │  │
│  │  - React 19     │  │                 │  │  - Handlers     │  │
│  │  - MUI/Bootstrap│  │                 │  │  - Services     │  │
│  └─────────────────┘  │                 │  └────────┬────────┘  │
│                       │                 │           │           │
│  Puerto: 443 (HTTPS)  │                 │  ┌────────▼────────┐  │
│                       │                 │  │ Persistent Disk │  │
└───────────────────────┘                 │  │   (1GB SSD)     │  │
                                          │  │                 │  │
                                          │  │  ┌───────────┐  │  │
                                          │  │  │ SQLite DB │  │  │
                                          │  │  │ payvue.db │  │  │
                                          │  │  └───────────┘  │  │
                                          │  │                 │  │
                                          │  │  ┌───────────┐  │  │
                                          │  │  │  Uploads  │  │  │
                                          │  │  │  Folder   │  │  │
                                          │  │  └───────────┘  │  │
                                          │  └─────────────────┘  │
                                          │                       │
                                          │  Puerto: 8080         │
                                          └───────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                        GITHUB REPOSITORY                                     │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    payvue_proyecto_software/                          │   │
│  │  ├── frontend/          → Despliega a Vercel (auto)                  │   │
│  │  ├── cmd/               → Código Go del backend                       │   │
│  │  ├── pkg/               → Paquetes Go                                 │   │
│  │  ├── Dockerfile         → Build del backend                           │   │
│  │  └── render.yaml        → Configuración de Render                     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Trigger: Push a main branch                                                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Proceso de Despliegue Paso a Paso

### 2.1 Pre-requisitos

- Cuenta en GitHub con el repositorio del proyecto
- Cuenta en Vercel (gratuita)
- Cuenta en Render (gratuita)

### 2.2 Despliegue del Frontend en Vercel

#### Paso 1: Conectar Repositorio

1. Ir a [vercel.com](https://vercel.com) e iniciar sesión con GitHub
2. Click en "Add New Project"
3. Seleccionar el repositorio `payvue_proyecto_software`

#### Paso 2: Configurar Build

```
Root Directory: frontend
Framework Preset: Create React App
Build Command: npm run build
Output Directory: build
Install Command: npm install
```

#### Paso 3: Configurar Variables de Entorno

| Variable | Valor |
|----------|-------|
| `REACT_APP_API_URL` | `https://payvue-api.onrender.com` |

#### Paso 4: Desplegar

1. Click en "Deploy"
2. Esperar a que el build termine (~2-3 minutos)
3. Obtener URL de producción (ej: `https://payvue.vercel.app`)

### 2.3 Despliegue del Backend en Render

#### Paso 1: Crear Nuevo Servicio

1. Ir a [render.com](https://render.com) e iniciar sesión
2. Click en "New" → "Web Service"
3. Conectar repositorio de GitHub

#### Paso 2: Configurar Servicio

```
Name: payvue-api
Region: Oregon (US West)
Branch: main
Root Directory: (dejar vacío)
Runtime: Docker
```

#### Paso 3: Configurar Dockerfile

```
Dockerfile Path: ./Dockerfile
Docker Build Context Directory: .
```

#### Paso 4: Configurar Variables de Entorno

| Variable | Valor |
|----------|-------|
| `PORT` | `8080` |
| `ENVIRONMENT` | `production` |
| `DATABASE_PATH` | `/app/data/payvue.db` |

#### Paso 5: Agregar Disco Persistente

1. En la sección "Disks", click "Add Disk"
2. Configurar:
   - Name: `payvue-data`
   - Mount Path: `/app/data`
   - Size: `1 GB`

#### Paso 6: Desplegar

1. Click en "Create Web Service"
2. Esperar build del Docker (~5-10 minutos)
3. Verificar health check en `/health`

---

## 3. Verificación del Despliegue

### 3.1 Checklist de Verificación

| # | Verificación | Comando/URL | Esperado |
|---|--------------|-------------|----------|
| 1 | Backend Health | `GET https://payvue-api.onrender.com/health` | "OK - PayVue API Server" |
| 2 | Frontend carga | Visitar URL de Vercel | Página de login |
| 3 | CORS funciona | Login desde frontend | Sin errores CORS |
| 4 | BD persistente | Crear datos, reiniciar, verificar | Datos persisten |
| 5 | Registro usuario | POST /auth/register | 201 Created |
| 6 | Login | POST /auth/login | 200 OK |

### 3.2 Pruebas Post-Despliegue

```bash
# Verificar health del backend
curl https://payvue-api.onrender.com/health

# Probar registro
curl -X POST https://payvue-api.onrender.com/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456"}'

# Probar login
curl -X POST https://payvue-api.onrender.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456"}'

# Probar GET incomes
curl https://payvue-api.onrender.com/finances/income
```

---

## 4. Arquitectura de Despliegue

### 4.1 Componentes

```
┌─────────────────────────────────────────────────────────────────┐
│                    ARQUITECTURA DE PRODUCCIÓN                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐  │
│  │    VERCEL    │      │    RENDER    │      │   GITHUB     │  │
│  │              │      │              │      │              │  │
│  │  ┌────────┐  │      │  ┌────────┐  │      │  ┌────────┐  │  │
│  │  │Frontend│  │◀────▶│  │Backend │  │◀────▶│  │  Repo  │  │  │
│  │  │ React  │  │ API  │  │   Go   │  │ Push │  │  Code  │  │  │
│  │  └────────┘  │      │  └────────┘  │      │  └────────┘  │  │
│  │              │      │       │      │      │              │  │
│  │  CDN Global  │      │  ┌────▼───┐  │      │  CI/CD       │  │
│  │  SSL Auto    │      │  │ SQLite │  │      │  Webhooks    │  │
│  │  Builds Auto │      │  │   DB   │  │      │              │  │
│  └──────────────┘      │  └────────┘  │      └──────────────┘  │
│                        │              │                         │
│                        │  Disk 1GB    │                         │
│                        │  SSL Auto    │                         │
│                        │  Health Check│                         │
│                        └──────────────┘                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Flujo de CI/CD

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐
│  Code   │────▶│  Push   │────▶│  Build  │────▶│ Deploy  │
│ Change  │     │ GitHub  │     │  Auto   │     │  Auto   │
└─────────┘     └─────────┘     └─────────┘     └─────────┘
                     │               │               │
                     │          ┌────┴────┐          │
                     │          ▼         ▼          │
                     │     ┌─────────┐ ┌─────────┐   │
                     │     │ Vercel  │ │ Render  │   │
                     │     │  Build  │ │  Build  │   │
                     │     └────┬────┘ └────┬────┘   │
                     │          │           │        │
                     │          ▼           ▼        │
                     │     ┌─────────┐ ┌─────────┐   │
                     └────▶│Frontend │ │Backend  │◀──┘
                           │   Live  │ │   Live  │
                           └─────────┘ └─────────┘
```

---

## 5. Configuración de Dominio Personalizado (Opcional)

### 5.1 Frontend (Vercel)

1. En Vercel Dashboard → Settings → Domains
2. Agregar dominio: `app.midominio.com`
3. Configurar DNS:
   ```
   CNAME: app → cname.vercel-dns.com
   ```

### 5.2 Backend (Render)

1. En Render Dashboard → Service → Settings → Custom Domain
2. Agregar dominio: `api.midominio.com`
3. Configurar DNS:
   ```
   CNAME: api → payvue-api.onrender.com
   ```

---

## 6. Monitoreo y Logs

### 6.1 Render (Backend)

- Dashboard → Service → Logs (logs en tiempo real)
- Métricas: CPU, Memory, Network
- Alertas configurables

### 6.2 Vercel (Frontend)

- Dashboard → Project → Analytics
- Métricas: Page Views, Visitors, Performance
- Edge Functions logs

---

## 7. Rollback y Versionado

### 7.1 Vercel

```bash
# Ver deployments anteriores
vercel ls

# Rollback a deployment anterior
vercel rollback <deployment-url>
```

### 7.2 Render

1. Dashboard → Service → Events
2. Click en deployment anterior
3. "Rollback to this deploy"

---

## 8. Costos Estimados

| Servicio | Plan | Costo |
|----------|------|-------|
| Vercel | Hobby | $0/mes |
| Render | Free | $0/mes |
| **Total** | | **$0/mes** |

**Nota**: Los planes gratuitos tienen limitaciones:
- Render: El servicio se "duerme" después de 15 min de inactividad
- Vercel: Límite de 100GB bandwidth/mes

---

## 9. Troubleshooting

### Error: "Service unavailable" en Render

**Causa**: El servicio está dormido (plan gratuito)
**Solución**: Esperar ~30 segundos, el servicio se activará automáticamente

### Error: CORS en producción

**Causa**: `CORS_ALLOWED_ORIGINS` no incluye el dominio de Vercel
**Solución**: Agregar URL de Vercel a la configuración CORS

### Error: Base de datos vacía tras redeploy

**Causa**: No se configuró disco persistente
**Solución**: Agregar disco en Render Dashboard → Service → Disks

### Error: Build falla en Render

**Causa**: Dependencias de CGO (SQLite)
**Solución**: El Dockerfile ya incluye `gcc musl-dev sqlite-dev`

