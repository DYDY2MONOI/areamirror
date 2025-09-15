# Docker - Area Web

Configuration Docker pour l'application React/Vite avec multi-stage build.

## 🚀 Commandes

```bash
# Développement (hot-reload sur localhost:5173)
npm run docker:dev

# Production (Nginx sur localhost:80)
npm run docker:prod

# Construire l'image
npm run docker:build-image

# Nettoyer les conteneurs
npm run docker:clean
```

## 📁 Structure

- **dev** : Environnement de développement avec hot-reload
- **builder** : Construction de l'application
- **production** : Serveur Nginx optimisé