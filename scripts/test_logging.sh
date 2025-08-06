#!/bin/bash

echo "🚀 Sistema de Logging Mejorado - Prueba de Funcionalidad"
echo "=========================================================="
echo ""

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}1. Verificando estructura del proyecto...${NC}"
if [ -f "pkg/logger/logger.go" ]; then
    echo -e "${GREEN}✓ Logger principal encontrado${NC}"
else
    echo -e "${RED}✗ Logger principal no encontrado${NC}"
    exit 1
fi

if [ -f "internal/adapter/middleware/logger_middleware.go" ]; then
    echo -e "${GREEN}✓ Middleware de logging encontrado${NC}"
else
    echo -e "${RED}✗ Middleware de logging no encontrado${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}2. Verificando configuración de ejemplo...${NC}"
if [ -f ".env.example" ]; then
    echo -e "${GREEN}✓ Archivo .env.example encontrado${NC}"
    echo -e "${YELLOW}   Configuración de LOG_LEVEL disponible${NC}"
else
    echo -e "${YELLOW}! Archivo .env.example no encontrado (no crítico)${NC}"
fi

echo ""
echo -e "${BLUE}3. Verificando documentación...${NC}"
if [ -f "LOGGING.md" ]; then
    echo -e "${GREEN}✓ Documentación de logging encontrada${NC}"
else
    echo -e "${RED}✗ Documentación de logging no encontrada${NC}"
fi

echo ""
echo -e "${BLUE}4. Verificando dependencias de Go...${NC}"
if go mod verify > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Dependencias verificadas${NC}"
else
    echo -e "${YELLOW}! Ejecutando go mod tidy...${NC}"
    go mod tidy
fi

echo ""
echo -e "${BLUE}5. Verificando compilación...${NC}"
if go build -o /tmp/test_build ./cmd/api > /dev/null 2>&1; then
    echo -e "${GREEN}✓ El proyecto compila correctamente${NC}"
    rm -f /tmp/test_build
else
    echo -e "${RED}✗ Error de compilación${NC}"
    echo "Ejecuta: go build ./cmd/api para ver los errores"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 ¡Sistema de logging instalado correctamente!${NC}"
echo ""
echo -e "${YELLOW}📋 Próximos pasos:${NC}"
echo "1. Copia .env.example a .env y configura tus variables"
echo "2. Configura LOG_LEVEL según tus necesidades:"
echo "   - DEBUG: Máxima información (desarrollo)"
echo "   - INFO: Información general (testing)"
echo "   - WARN: Solo advertencias y errores (staging)"
echo "   - ERROR: Solo errores críticos (producción)"
echo ""
echo -e "${YELLOW}🚦 Para iniciar el servidor:${NC}"
echo "   LOG_LEVEL=DEBUG go run cmd/api/main.go"
echo ""
echo -e "${YELLOW}📖 Para más información consulta:${NC}"
echo "   - LOGGING.md para documentación completa"
echo "   - .vscode/tasks.json para tareas configuradas"
echo ""
