# Задание 1

## Запуск
```bash
docker build --rm --tag ddas/lesson1 .
docker run --name lesson1 -p 127.0.0.1:8080:8080 -d --rm ddas/lesson1
```