# Stage 1: build CSS
FROM node:20-alpine AS css-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npx tailwindcss -i ./static/style.css -o ./static/output.css --minify

# Stage 2: build Go binary
FROM golang:1.23-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=css-builder /app/static/output.css ./static/output.css
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/web

# Stage 3: final lean image
FROM alpine:latest
WORKDIR /app
# needed for sqlite CGO
RUN apk add --no-cache sqlite-libs
COPY --from=go-builder /app/server .
COPY --from=go-builder /app/templates ./templates
COPY --from=css-builder /app/static ./static
EXPOSE 8080
CMD ["./server"]
