# Задание 2
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
docker build --rm --tag ddas/lesson2 .
docker run --rm -p 127.0.0.1:8080:8080 -d --name lesson2 ddas/lesson2
```