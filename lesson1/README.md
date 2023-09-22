# Задание 1
- HTTP backend
- Запускается и ждет
- и обрабатывает 2 роута (http обработчика)
- /replace
- принимает body и сохраняет
- и отдаёт 200 ОК
- /get
- возвращает сохранённый body

## Запуск
```bash
docker build --rm --tag ddas/lesson1 .
docker run --name lesson1 -p 127.0.0.1:8080:8080 -d --rm ddas/lesson1
```
