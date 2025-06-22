# drip

**drip** es una herramienta escrita en Go para desplegar y exponer modelos de inteligencia artificial en la red de forma eficiente, controlada y escalable.

## Características

- 🌐 **Exposición de modelos IA**: Sirve modelos a través de endpoints HTTP de forma sencilla y segura.
- 📊 **Métricas integradas**: Monitorea el rendimiento de los modelos con métricas listas para Prometheus.
- 🚦 **Tipos de release avanzados**: Soporta estrategias de despliegue como:
  - **A/B Testing**
  - **Canary Releases**
  - **Shadow Deployments**
- ⚙️ **Autoescalado configurable**: Ajusta automáticamente los recursos según la carga del sistema y las métricas personalizadas.
- 🛠️ **Configuración declarativa**: Todo el despliegue puede definirse mediante archivos YAML o JSON.

## Instalación

```bash
go install github.com/MiguelAMeloM/drip
