FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
COPY apps/web/package.json apps/web/package.json
RUN npm ci

FROM deps AS build
COPY apps/web apps/web
WORKDIR /app/apps/web
ENV VITE_API_BASE_URL=
RUN npm run build

FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/package.json ./package.json
COPY --from=deps /app/apps/web/package.json ./apps/web/package.json
COPY --from=build /app/apps/web/dist ./dist
COPY apps/web/server.mjs ./server.mjs
EXPOSE 3000
CMD ["node", "server.mjs"]
