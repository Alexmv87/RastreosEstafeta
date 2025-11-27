# RastreosEstafeta

Aplicación pequeña en Go que consulta el historial de seguimiento de guías en Estafeta (scraping) y expone un endpoint HTTP para procesar múltiples guías concurrentemente.

**Propósito**
- Proveer un servicio HTTP simple que reciba un array de guías y devuelva el estado/fecha/hora para cada una, obtenidos desde el servicio web de Estafeta mediante scraping con la librería `colly`.

**Endpoint**
- `POST /buscarGuia` : recibe datos en JSON y devuelve un array con `TrackingResponse`.

Formatos de entrada soportados
- Array JSON puro: `["8015905005552600009651","OTRA_GUIA"]`

Respuesta
- JSON: arreglo de objetos `TrackingResponse` con los campos:
  - `guia` (string)
  - `fecha` (string)
  - `hora` (string)
  - `status` (string)

Concurrencia
- La aplicación procesa las guías en paralelo usando gorutinas y un semáforo (canal buffered) para limitar la concurrencia. Por defecto el límite no está fijado, para agregar un límite se descomentarían las líneas de la función "buscarGuia"  referentes a las variables "maxConcurrent" (donde se indicaría la cantidad máxima de consultas simultáneas) y la variable sem.

Ejecución local (PowerShell)

1. Compilar:
```powershell
go build main.go
```

2. Ejecutar:
```powershell
.\main.exe
# o
go run main.go
```

3. Petición de prueba (array JSON puro):
```powershell
Invoke-RestMethod -Uri 'http://localhost:9000/buscarGuia' -Method Post -Body '["8015905005552600009651","OTRA_GUIA"]' -ContentType 'application/json'
```

4. Petición de prueba de guía unitaria :
```powershell
Invoke-RestMethod -Uri 'http://localhost:9000/buscarGuia/8015905005552600009651' -Method GET
```

Qué verás
- La respuesta HTTP contendrá el array de `TrackingResponse` con la información (o un mensaje en `status` si no hay información).

Notas y recomendaciones
- Cada invocación de `busqueda` crea su propio `colly.Collector`, por lo que es seguro ejecutar varias en paralelo.
- Ajusta `maxConcurrent` en `buscarGuia` según la carga y límites del servicio externo.
- Considerar añadir timeouts por guía (usando `context` o `select` con `time.After`) para evitar gorutinas bloqueadas en llamadas lentas.
- Para producción: añadir logs estructurados, manejo de errores más detallado y tests unitarios/integración.

Licencia
- Código proporcionado sin licencia específica. Añade el archivo `LICENSE` si quieres aplicar una licencia.

Contacto
- Repositorio: `RastreosEstafeta` (owner: `Alexmv87`)
