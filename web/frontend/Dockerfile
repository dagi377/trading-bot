# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Create necessary directories and set permissions
RUN mkdir -p /var/cache/nginx /var/run /var/log/nginx /run && \
	touch /run/nginx.pid && \
	chown -R nginx:nginx /var/cache/nginx /var/run /var/log/nginx /run /etc/nginx/conf.d && \
	chmod -R 755 /var/cache/nginx /var/run /var/log/nginx /run /etc/nginx/conf.d && \
	chmod 644 /run/nginx.pid

# Copy built assets from build stage
COPY --from=build /app/build /usr/share/nginx/html
RUN chown -R nginx:nginx /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN chown nginx:nginx /etc/nginx/conf.d/default.conf

EXPOSE 80

# Use nginx user
USER nginx

CMD ["nginx", "-g", "daemon off;"]
