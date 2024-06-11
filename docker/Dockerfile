FROM ubuntu:latest
LABEL authors="antonsotnik"

ENTRYPOINT ["top", "-b"]

# Используем официальный образ Golang в качестве базового образа
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем содержимое текущей директории (где находится Dockerfile) внутрь контейнера в директорию /app
COPY . .

# Собираем приложение
RUN go build cmd/main.go

# Команда для запуска приложения при старте контейнера
CMD ["./main"]
