#!/bin/bash
MODULE_PATH="../modulo-kernel/continfo.ko"
SYS_MODULE_PATH="../modulo-kernel/sysinfo.ko"

if [ -f "$MODULE_PATH" ]; then
    sudo insmod "$MODULE_PATH"
    echo "continfo módulo cargado"
fi

if [ -f "$SYS_MODULE_PATH" ]; then
    sudo insmod "$SYS_MODULE_PATH"
    echo "sysinfo módulo cargado"
fi
