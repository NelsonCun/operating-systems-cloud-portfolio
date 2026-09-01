# PROYECTO 3 - SISTEMAS OPERATIVOS 1

Estudiante: Nelson Emanuel Cún Bálan
Carné: 201222010

## Objetivo
Implementar un sistema de monitoreo de datos del clima de diferentes municipios utilizando tecnologías de contenedores y orquestación con Kubernetes.

## Descripción Técnica
El proyecto consiste en desarrollar una aplicación que recopile datos climáticos de varios municipios y los almacene en una base de datos Valkey. La aplicación está contenida en un contenedor Docker y se orquesta utilizando Kubernetes para asegurar su escalabilidad y disponibilidad. Se utiliza un cluster de GCP para desplegar la aplicación.

### Componentes Principales
1. **Contenedor Docker**: La aplicación se empaqueta en un contenedor Docker que incluye todas las dependencias necesarias para su ejecución.
2. **Kubernetes**: Se utiliza Kubernetes para gestionar el despliegue, escalado y operación de la aplicación en el cluster de GCP.
3. **Base de Datos Valkey**: Se utiliza Valkey como base de datos para almacenar los datos climáticos recopilados.
4. **Cluster en GCP**: El cluster de Google Cloud Platform proporciona la infraestructura necesaria para ejecutar la aplicación de manera eficiente.
5. **VM zot**: Se utiliza una máquina virtual llamada "zot" para alojar ciertos servicios o componentes de la aplicación, especialmente aquellas imágenes que realicé, las cuales son la api-rust, el server_go, los writers tanto de kafka como de rabbitmq.
También los consumidores de kafka y rabbitmq.
6. **Monitoreo y Escalabilidad**: Se implementan mecanismos de monitoreo para asegurar el correcto funcionamiento de la aplicación y se configuran políticas de escalabilidad automática para manejar variaciones en la carga de trabajo.

### Despliegue
El despliegue de la aplicación se realiza mediante archivos YAML que definen los recursos de Kubernetes necesarios, incluyendo despliegues, servicios y políticas de escalabilidad.

Para poder ejecutarlo se deben seguir los siguientes pasos:
1. Clonar el repositorio del proyecto.
2. Construir las imágenes Docker necesarias y subirlas a un registro de contenedores accesible desde el cluster de GCP.
3. Aplicar los archivos YAML utilizando `kubectl apply -f <archivo.yaml>` para desplegar la aplicación en el cluster.
4. Verificar el estado de los pods y servicios utilizando `kubectl get pods` y `kubectl get services`.

### Arquitectura del Sistema
La arquitectura del sistema se basa en una estructura de microservicios donde cada componente de la aplicación se ejecuta en su propio contenedor. Los datos climáticos son recopilados por la aplicación principal y enviados a la base de datos Valkey para su almacenamiento y posterior análisis.

![Arquitectura del Sistema](./images/01.png)

### Locust
Para realizar pruebas de carga y rendimiento en la aplicación, se utiliza Locust, una herramienta de código abierto para pruebas de carga. Locust permite simular múltiples usuarios concurrentes que interactúan con la aplicación, lo que ayuda a identificar posibles cuellos de botella y optimizar el rendimiento.

![Pruebas con Locust](./images/02.png)

### API Rest
Se implementa una API REST para permitir la interacción con la aplicación. Esta recibe de parte de locust por medio de un ingress los datos climáticos y los envía a la base de datos Valkey.

### Servicios GRPC
Además de la API REST, se implementan servicios GRPC para la comunicación eficiente entre los diferentes componentes de la aplicación. GRPC es un framework de comunicación de alto rendimiento que utiliza HTTP/2 y Protobuf para la serialización de datos.

### Kafka
Se utiliza Apache Kafka como sistema de mensajería para manejar la transmisión de datos entre los diferentes componentes de la aplicación. Kafka permite una comunicación asíncrona y escalable, lo que es ideal para aplicaciones distribuidas.

### RabbitMQ
RabbitMQ se utiliza como otro sistema de mensajería para complementar a Kafka. Proporciona una solución robusta para la gestión de colas de mensajes y la comunicación entre servicios.

Tanto Kafka como RabbitMQ son utilizados para asegurar que los datos climáticos sean transmitidos de manera eficiente y confiable entre los diferentes componentes de la aplicación.

### Valkey
Valkey es la base de datos elegida para almacenar los datos climáticos recopilados por la aplicación. Valkey es una base de datos en memoria que ofrece alta velocidad y rendimiento, lo que es crucial para aplicaciones que requieren acceso rápido a grandes volúmenes de datos.

### Grafana
Grafana se utiliza para la visualización y monitoreo de los datos almacenados en Valkey. Permite crear dashboards personalizados que muestran métricas clave y tendencias en los datos climáticos, facilitando el análisis y la toma de decisiones.

![Dashboard en Grafana](./images/03.png)

### Zot
Zot se utiliza como un registro de contenedores privado para almacenar las imágenes Docker personalizadas utilizadas en el proyecto. Esto permite un acceso rápido y seguro a las imágenes desde el cluster de GCP. En este caso Zot se despliega en una VM aparte para alojar las imágenes que realicé.

![Zot en VM](./images/04.png)

### HPA (Horizontal Pod Autoscaler)
Se implementa el Horizontal Pod Autoscaler (HPA) para asegurar que la aplicación pueda escalar automáticamente en función de la carga de trabajo. HPA monitorea métricas como el uso de CPU y ajusta el número de pods en consecuencia, garantizando un rendimiento óptimo bajo diferentes condiciones de carga.

### VPA (Vertical Pod Autoscaler)
Además del HPA, se utiliza el Vertical Pod Autoscaler (VPA) para ajustar automáticamente los recursos asignados a cada pod. VPA analiza el uso de recursos y ajusta la memoria y CPU asignados a los pods para optimizar el rendimiento y la eficiencia.


### Consideraciones Finales
El proyecto busca demostrar la capacidad de utilizar tecnologías modernas de contenedores y orquestación para desarrollar aplicaciones escalables y resilientes en la nube. La elección de Valkey como base de datos y la implementación en un cluster de GCP aseguran un rendimiento óptimo y una alta disponibilidad de los datos climáticos recopilados.