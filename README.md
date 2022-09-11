# Структура проекта
```
├── cmd
│   ├── agent                                    Точка запуска агента
│   └── server                                   Точка запуска сервера
│
├── internal
│   ├── agent
│   │   ├── config                               Конфиг агента
│   │   ├── memory                               Функциональность работы с метриками агента: создание, обновление, считывание
│   │   └── service                              Функционал агента
│   │       ├── interfaces.go                              Интерфейсы, которыми пользуется агент         
│   │       ├── mocks                                      Mock'и интерфейсов (используются для тестов)
│   │       ├── request.go                                 Реализация http-запросов
│   │       ├── request_test.go                            Юнит-тесты http-запросов
│   │       └── service.go                                 Реализация создания, запуска агента
│   └── server
│       ├── api
│       │   ├── api.go                                     Реализация API сервера
│       │   ├── compress.go                                Middleware'ы архивации/деархивации тел запросов/ответов
│       │   ├── errors.go                                  Возможные ошибки
│       │   ├── handlerMetric.go                           Обработчики чтения, обновления метрик
│       │   ├── handlerMetricPathParams.go                 Обработчики чтения, обновления метрик (тело запроса в URL)
│       │   ├── handlerMetric_test.go                      Юнит-тесты http-запросов
│       │   ├── handlerPage.go                             Обработчик открытия стартовой страницы
│       │   ├── handlerPing.go                             Обработчик-пинг хранилища данных
│       │   ├── interfaces.go                              Интерфейсы, которыми пользуется API
│       │   ├── mocks                                      Mock'и интерфейсов (используются для тестов)
│       │   ├── page.go                                    Внутренняя реализации чтения стартовой страницы
│       ├── config                               Конфиг сервера
│       ├── model                                Модели работы со списком метрик
│       ├── service                              Слой приложения
│       └── storage                              Слой хранилища
│           ├── interfaces.go                              Интерфейсы, которыми пользуется хранилище                    
│           ├── memory                                     in-memory хранилище             
│           ├── pg                                         postgres хранилище
│           └── storage.go                                 Конструктор хранилища и определение интерфейса хранилища
├── pkg
│   └── metric                                   Реализация единицы метрики
```

# Блок-схема работы агента
```mermaid
graph TD;
    A[Agent start] -->|goroutine upd additional metrics| C{upd}
    A[Agent start] -->|goroutine upd basic metrics| B{upd}
    A[Agent start] -->|goroutine report metrics| D{report}

    C -->|repeat| C
    C -->|shutdown| END

    B -->|repeat| B
    B -->|shutdown| END

    D -->|repeat| D
    D -->|shutdown| END
```

# Блок-схема работы сервера
```mermaid
flowchart TB
    SERVER[Server]

    subgraph API
    Handlers --> HCRUD([CRUD metrics])
    Handlers --> HPAGE([page])
    Handlers --> HPING([ping db conn])
    end

    subgraph SERVICE
    SCRUD([CRUD metrics])
    SHUTDOWN(Shutdown)
    end

    SERVER --> API

    API <--> SERVICE

    SERVICE <--> STORAGE[(Storage)]
``` 

# Обновление шаблона

Чтобы получать обновления автотестов и других частей шаблона, выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-devops-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.
