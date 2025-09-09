FROM node:24-alpine3.21 AS build-client
WORKDIR /build
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build:ui

FROM golang:1.25.1-alpine3.21 AS build-server
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM scratch
WORKDIR /bin
COPY --from=build-client /build/dist ./dist
COPY --from=build-server /build/main ./main
EXPOSE 3001
CMD ["/bin/main"]
