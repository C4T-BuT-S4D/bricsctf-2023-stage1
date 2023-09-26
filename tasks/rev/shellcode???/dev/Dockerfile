# Используем образ с ARM архитектурой в качестве базового
FROM gcc:latest

# Копируем ваш исходный код в контейнер
WORKDIR /app
COPY . .

# Компилируем x86 код с помощью кросс-компилятора
RUN gcc -s -o my_x86_program main.c -lm

# Запускаем программу x86
CMD tail -f /dev/null
