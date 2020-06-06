#Веб-сервер для социальной сети meet&go
Основные параметры указаны в конфигурационном файле (config.json).

    {
    "dataBase": {               //Настройки базы данных
        "username": "postgres", //Имя пользователя
        "password": "password", //Пароль
        "dbName": "name",       //Название базы данных
        "ip": "127.0.0.1",      //IP адрес базы данных
        "port": 5432,           //Порт для подключения
        "reInitDataBase": false //Указать true для переинициализации базы данных
    },
    "log": true                 //Включить режим логирования 
    }

Стоит учесть, что данный сервер работет только с базой данных PostgreSQL.

При указании парамемтра __reInitDataBase__ в конфигурационном файле удаляются все существующие таблицы со всеми данными
и пересоздаются (предполагалось, что его стоит использовать для инициализации таблиц в базе данных)

При указании параметра __log__ в конфигурационном файле веб-сервер будет осуществлять журналирование запросов 
в таблицу log в развернутой базе данных.

## Запросы

На данный момент реализованы следующие запросы:

`hostname/registration "POST"` - запрос на регистрацию

`hostname/login "POST"` - запрос на авторизацию

`hostname/profiles/{id} "GET"` - запрос на получение информации о профиле

`hostname/search "GET"` - запрос на поиск пользователей

`hostname/delete "POST"` - запрос на удаление аккаунта

`hostname/tags "GET"` - запрос на получение всех тэгов

`hostname/tags/add "POST"` - запрос на добавление себе тэгов

`hostname/tags/get/{id} "GET"` - запрос на получение тэгов пользователя

`hostname/tasks/{taskName} "GET"` - запрос на получение информации о задании

`hostname/tasks/{taskName}/tags "GET"` - запрос на получение

router.HandleFunc("/invite/user/{id}", handlers.InviteUser).Methods("POST")

router.HandleFunc("/invite/show", handlers.GetInvites).Methods("GET")

router.HandleFunc("/validate/show", handlers.GetTasksToValidate).Methods("GET")

router.HandleFunc("/quests/show", handlers.GetQuests).Methods("GET")

router.HandleFunc("/quest/status/change", handlers.ChangeQuestStatus).Methods("POST")

А так же несколько запросов для разработчиков:

`hostname/developers/getAccount "GET"` - запрос на получение аккаунта разработчика 

`hostname/developers/postTag "POST"` - запрос на добавление тэга

### hostname/registration "POST"
Для данного запроса требуется прикрепить 