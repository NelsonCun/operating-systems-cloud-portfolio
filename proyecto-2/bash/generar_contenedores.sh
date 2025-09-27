#!/bin/bash
LOG_FILE="$(dirname "$0")/../go-daemon/logs/docker_creacion.log"
IMAGENES=("alto_cpu_ram_1" "alto_cpu_ram_2" "bajo_consumo")

mkdir -p "$(dirname "$LOG_FILE")"
echo "==== Creación de contenedores $(date) ====" >> "$LOG_FILE"

for i in $(seq 1 10); do
    IMG=${IMAGENES[$RANDOM % ${#IMAGENES[@]}]}
    NOMBRE="cont_$IMG_$RANDOM"
    docker run -d --name "$NOMBRE" "$IMG"
    echo "Contenedor $NOMBRE creado con imagen $IMG" >> "$LOG_FILE"
done

docker ps -aq --filter "status=exited" | xargs -r docker rm
