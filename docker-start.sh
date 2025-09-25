#!/bin/bash

echo "Démarrage de la stack Area..."

docker-compose down
docker-compose up --build -d

echo "Services démarrés :"
echo "- Frontend: http://localhost:3000"
echo "- Backend: http://localhost:8080"
echo "- Database: localhost:5433"

echo "Pour voir les logs: docker-compose logs -f"
echo "Pour arrêter: docker-compose down"
