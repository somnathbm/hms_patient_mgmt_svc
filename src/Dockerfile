FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api/
COPY ./db ./db/
COPY ./models ./models/
COPY ./utils ./utils/
COPY ./metrics ./metrics/
COPY main.go ./main.go

RUN go build -o ./dist/hms-patient-mgmt.sh

FROM golang:alpine AS run
WORKDIR /app
COPY --from=build --chmod=500 /app/dist/hms-patient-mgmt.sh ./
EXPOSE 8080
EXPOSE 8081
CMD ["./hms-patient-mgmt.sh"]
