# ARCHIVOS Y COMANDOS UTILIZADOS EN LA TAREA 1
Este repositorio contiene los archivos y comandos utilizados en la Tarea 1.

Curso: Sistemas Operativos 1  

- Nombre: Nelson Emanuel Cún Bálan 
- Carné: 201222010
- Sección: P 

---

## Navegación de directorios

### Comando 'cd'
El comando `cd` (change directory) se utiliza para cambiar el directorio de trabajo actual en la línea de comandos.

![cd command](./images/01.png)

### Comando 'ls'
El comando `ls` (list) se utiliza para listar los archivos y directorios en el directorio actual.

![ls command](./images/02.png)

### Comando 'pwd'
El comando `pwd` (print working directory) se utiliza para mostrar la ruta del directorio de trabajo actual.

![pwd command](./images/03.png)

## Manipulación de archivos

### Comando 'touch'
El comando `touch` se utiliza para crear un archivo vacío o actualizar la marca de tiempo de un archivo existente.

![touch command](./images/04.png)

Archivo creado: `prueba_touch.txt`

![touch command result](./images/05.png)

### Comando 'cp'
El comando `cp` (copy) se utiliza para copiar archivos o directorios de una ubicación a otra.

Agregué un mensaje de prueba al archivo `prueba_touch.txt`.

![cp command](./images/06.png)

Uso del comando `cp` para copiar el archivo `prueba_touch.txt` a un nuevo archivo llamado `prueba_copia.txt`.

![cp command result](./images/07.png)

Resultado de la copia.

![cp command result](./images/08.png)

### Comando 'mv'
El comando `mv` (move) se utiliza para mover o renombrar archivos o directorios.

Uso del comando `mv` para renombrar el archivo `prueba_copia.txt` a `prueba_renombrado.txt`.

![mv command](./images/09.png)

Uso del comando `mv` para mover el archivo `prueba_renombrado.txt` a un nuevo directorio llamado `nueva_carpeta`.

![mv command result](./images/10.png)

### Comando 'rm'
El comando `rm` (remove) se utiliza para eliminar archivos.

![rm command](./images/11.png)

## Visualización de contenido

### Comando 'cat'
El comando `cat` (concatenate) se utiliza para mostrar el contenido de un archivo en la terminal.

![cat command](./images/12.png)

### Comando 'more'
El comando `more` se utiliza para visualizar el contenido de un archivo página por página.

![more command page1](./images/13.png)
![more command page2](./images/14.png)
![more command page3](./images/15.png)
![more command page4](./images/16.png)
![more command page5](./images/17.png)

### Comando 'less'
El comando `less` se utiliza para visualizar el contenido de un archivo de manera interactiva, permitiendo desplazarse hacia adelante y hacia atrás.

![less command](./images/18.png)

## Gestión de permisos

### Comando 'chmod'
El comando `chmod` (change mode) se utiliza para cambiar los permisos de archivos o directorios.

Usé 'chmod 644 prueba.txt' para establecer los permisos de lectura y escritura para el propietario, y solo lectura para el grupo y otros usuarios.

* (Usuario) 6 = lectura (4) + escritura (2)
* (Grupo) 7 = lectura (4) + ejecución (1)
* (Otros) 4 = lectura (4)

![chmod command](./images/19.png)

### Comando 'chown'
El comando `chown` (change owner) se utiliza para cambiar el propietario de un archivo o directorio.

![chown command](./images/20.png)

## Script simulacion_crear_contenedores.sh

Creación del archivo del script `simulacion_crear_contenedores.sh` utilizando el editor de texto `nano`.

![nano command](./images/21.png)

El script contiene los siguientes comandos:

```bash
!/bin/bash

# Generar un número aleatorio entre 1 y 4
NUM_ARCHIVOS=$((RANDOM % 4 + 1))

echo "El número de archivos creados es de $NUM_ARCHIVOS."

# Nombre aleatorio
generar_nombre() {
    tr -dc 'A-Za-z' < /dev/urandom | head -c 6
}

# Crear archivos
for ((i=1; i<=NUM_ARCHIVOS; i++))
do
    NOMBRE_ALEATORIO=$(generar_nombre)
    NOMBRE_ARCHIVO="contenedor_201222010_${NOMBRE_ALEATORIO}.txt"
    
    # Crear archivo y escribir su propio nombre como contenido
    echo "$NOMBRE_ARCHIVO" > "$NOMBRE_ARCHIVO"
    
    echo "Archivo creado: $NOMBRE_ARCHIVO"
done

echo "Todos los archivos han sido creados."
```

### Ejecución del script

Di los permisos de ejecución al script utilizando el comando `chmod +x simulacion_crear_contenedores.sh`, lo cual significa que el archivo puede ser ejecutado como un programa.

![chmod +x command](./images/22.png)

Ejecución del script utilizando el comando `./simulacion_crear_contenedores.sh`.

![script execution](./images/23.png)


### Contenido dentro de los archivos creados

![file content 1](./images/24.png)
