# Stage 1: Build the React application
FROM node:18 AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Stage 2: Production image with Express server
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/dist ./dist
COPY server.js ./server.js

EXPOSE 8080

ENV PORT=8080 \
    LEAVE_API_URL=http://localhost:8081 \
    MANAGER_API_URL=http://localhost:8082

CMD ["node", "server.js"]
