FROM node:20-alpine AS base

RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY package*.json ./

# dev
FROM base AS dev
RUN npm install
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# build
FROM base AS builder
RUN npm install --only=production
COPY . .
RUN npm run build

# prod
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
