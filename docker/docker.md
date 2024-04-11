## docker

### Задание 1. Создание Docker контейнера

Создайте Docker контейнер с простым веб-сервером (например, на основе Nginx),
который отображает "Hello, Docker!" при обращении к корневому URL.

- Создайте Dockerfile
- Создайте index.html
- Постройте и запустите контейнер
- Теперь, когда вы открываете браузер и переходите по адресу http://localhost:8080,
  вы должны увидеть "Hello, Docker!"

//docker run -d nginx
docker pull nginx

### Задание 2. Передача аргументов в Docker контейнер

Измените предыдущий Dockerfile так, чтобы текст приветствия можно было передать в контейнер через аргументы.

### Задание 3. Работа с многими контейнерами

Создайте файл docer-compose чтобы запустить одновременно два контейнера один с веб сервером Nginx другой с базой данных Postgres

### Задание 4. Передача данных между контейнерами

Создайте контейнер с приложением на которое отправляет запрос к базе данных созданной в предыдущем задании и выводит результат  
Создайте для приложения  
Создайте простое приложение на
Создайте для обоих контейнеров Запустите контейнеры  
Теперь ваше приложение отправляет запрос к базе данных запущенной в соседнем контейнере

---

не прошел запрос, пока не завел сеть

### Задание 5. Монтирование томов

Создайте контейнер с веб сервером в который монтируется локальный каталог например с страницами внутрь контейнера

### Задание 6. Использование многоконтейнерного подхода

Создайте многоконтейнерное приложение включающее веб сервер и базу данных каждый из которых представлен своим собственным контейнером +

### Задание 7. Сборка образа на основе

Создайте который включает в себя сборку и запуск вашего приложения +

ПРОБЛЕМА -- build работает из bat-файла если dockerfile лежит в ./docker/app.dockerfile
но dockercompose build не может сделать COPY
приходится вытаскивать docker-compose.yaml && app.dockerfile в корень приложения (./)

### Задание 8. Использование сетей

Создайте два контейнера и поместите их в одну и ту же сеть чтобы они могли взаимодействовать друг с другом по именам контейнеров +

docker network create mynet

### Задание 9. Использование переменных окружения в

Используйте переменные окружения в своем для настройки параметров приложения +

### Задание 10. Работа с

Загрузите ваш образ на и поделитесь им с другими