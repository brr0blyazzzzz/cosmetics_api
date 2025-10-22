# cosmetics_api
Pепозиторий проекта RESTful API для управления базой данных косметических продуктов. 
# Cтруктура базы данных
База данных проекта состоит из четырех таблиц, связанных друг с другом, чтобы хранить информацию о косметических продуктах, их производителях и составе.

#### Таблица manufacturer: Содержит информацию о производителях.

    manufacturer_id (INTEGER, PRIMARY KEY)

    manufacturer_title (TEXT)

    country (TEXT)

    address (TEXT)

    contact_list (TEXT)

#### Таблица products: Хранит данные о продуктах.

    product_id (INTEGER, PRIMARY KEY)

    product_title (TEXT)

    product_description (TEXT)

    contraindications (TEXT, может быть NULL)

    application (TEXT)

    volume (TEXT)

    photo (TEXT)

    manufacturer_id (INTEGER, FOREIGN KEY, ссылается на manufacturer)

#### Таблица structure: Содержит список компонентов или ингредиентов.

    structure_id (INTEGER, PRIMARY KEY)

    structure_name (TEXT)

#### Таблица structure: Содержит список компонентов или ингредиентов.

    user_id (INTEGER, PRIMARY KEY)

    username (TEXT)

    password (TEXT, хешируется с помощью bcrypt)

#### Таблица product_structure: Таблица связи "многие ко многим" между продуктами и их составом.

    product_id (INTEGER, FOREIGN KEY, ссылается на products)

    structure_id (INTEGER, FOREIGN KEY, ссылается на structure)

# Cтруктура проекта
C:\Users\Polina\Desktop\учеба\project\
│
├── main.go                              # Точка входа в приложение
│
├── models\                              # Модели данных
│   └── models.go                        # Все структуры: User, Product, Manufacturer и др.
│
├── database\                            # Работа с базой данных
│   ├── database.go                      # Подключение к SQLite
│   └── Структура БД магазина.drawio     # Диаграмма структуры БД
│
├── repository\                          # CRUD-логика (работа с БД)
│   ├── manufacturer_repository.go       # Методы CRUD производителя
│   ├── product_repository.go            # Методы CRUD продукта
│   └── user_repository.go               # Методы для пользователей (регистрация, авторизация и т.д.)
│
├── handlers\                            # HTTP-обработчики запросов (контроллеры)
│   ├── manufacturer.go                  # CRUD-обработчики производителей
│   ├── product.go                       # CRUD-обработчики продуктов
│   ├── user.go                          # API-обработчики регистрации и логина (JSON)
│   ├── auth_forms.go                    # Обработка HTML-форм (login/register)
│   ├── jwt.go                           # Генерация и структура JWT-токенов
│   ├── admin_middleware.go              # Middleware проверки авторизации (JWT)
│   └── web.go                           # Обработка HTML-шаблонов (index.html, admin.html)
│
├── views\                               # HTML-шаблоны (frontend)
│   ├── index.html                       # Главная страница
│   ├── admin.html                       # Админ-панель
│   ├── login.html                       # Страница входа
│   └── register.html                    # Страница регистрации
│
├── queries\                             # Тестовые HTTP-запросы для проверки API
│   └── requests.http                    # Все HTTP-запросы
│
└── go.mod / go.sum                      # Go-модули и зависимости

