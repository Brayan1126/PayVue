# Frontend PayVue

## Descripción

Este es el frontend de PayVue, una aplicación web para la gestión de finanzas personales, desarrollada en React. Permite registrar ingresos, deudas, pagos y visualizar reportes y estadísticas de manera sencilla.

---

## Estructura del proyecto

- **src/**
  - **components/**: Componentes reutilizables (Sidebar, FinanceForm, etc).
  - **pages/**: Vistas principales (Dashboard, Records, PaymentHistory, Login, Register).
  - **utils/**: Utilidades como alertas personalizadas.
  - **App.js**: Configuración de rutas principales.
  - **index.js**: Punto de entrada de la aplicación.

---

## Scripts disponibles

En la carpeta del frontend puedes ejecutar:

- `npm start`  
  Inicia la app en modo desarrollo.  
  Abre [http://localhost:3000](http://localhost:3000) en tu navegador.

- `npm run build`  
  Genera una versión optimizada para producción en la carpeta `build`.

- `npm test`  
  Ejecuta los tests de React (si tienes tests implementados).

---

## Dependencias principales

- **React**: Librería principal para la interfaz.
- **axios**: Para llamadas HTTP a la API backend.
- **react-router-dom**: Navegación entre páginas.
- **bootstrap** y **bootstrap-icons**: Estilos y componentes visuales.
- **chart.js** y **react-chartjs-2**: Gráficas y reportes.
- **sweetalert2**: Alertas y confirmaciones visuales.

---

## Configuración y uso

1. Instala las dependencias:
   ```
   npm install
   ```

2. Inicia el servidor de desarrollo:
   ```
   npm start
   ```

3. Asegúrate de que el backend esté corriendo en [http://localhost:8000](http://localhost:8000) para que la app funcione correctamente.

---

## Personalización

- Puedes modificar los estilos en `src/index.css` o en los archivos CSS de cada componente.
- Las rutas principales están definidas en `src/App.js`.
- Las alertas personalizadas están en `src/utils/alert.js`.

---

## Notas

- El frontend está preparado para consumir la API REST del backend PayVue.
- Si cambias la URL del backend, actualízala en los archivos donde se hacen peticiones con axios.
- Para producción, genera el build y sirve los archivos estáticos desde tu servidor preferido.

---

## Créditos

Desarrollado por [Tu Nombre o Equipo].

---