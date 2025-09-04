#!/bin/bash

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
