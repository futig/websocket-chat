# WebSocket chat
Чат на WebSocket'ах (по мотивам https://xkcd.com/1305/)

Чтобы запустить, нужно сделать следующее:

1. в файле main.go изменить `ip` и `port` (по дефолту `ip = 0.0.0.0`, `port = 8080`)
2. запустить main.go

---

### Краткое описание протокола общения с сервером:

* `ping`, `pong` — ping-pong
* `{'mtype': 'INIT', 'id': 'ABCDEFGH'}` — инициализация, отправляется при входе пользователя в систему
    - `id`: идентификатор пользователя
* `{'mtype': 'TEXT', 'id': 'IDFROM', 'to': 'IDTO', 'text': 'message text'}` — отправка сообщения
    - `id`: идентификатор отправителя
    - `to`: идентификатор получателя (если пустой, сообщение отправляется всем)
    - `text`: текст сообщения
* `{'mtype': 'MSG', 'id': 'IDFROM', 'text': 'message'}` — сообщение с текстом `text` от пользователя `id`
* `{'mtype': 'DM', 'id': 'IDFROM', 'text': 'message'}` — «приватное» сообщение с текстом `text` от пользователя `id`
* `{'mtype': 'USER_ENTER', 'id': 'USERID'}` — служебное сообщение «пользователь `id` вошёл в чат»
* `{'mtype': 'USER_LEAVE', 'id': 'USERID'}` — служебное сообщение «пользователь `id` вышел из чата»

---

### Компоненты:
**Хаб** (hub.go): Центральный компонент, который управляет всеми подключениями клиентов. Он отвечает за регистрацию, разрегистрацию клиентов, рассылку сообщений всем или конкретным пользователям.

**Клиент** (client.go): Представляет отдельного пользователя. Читает сообщения от WebSocket-соединения и отправляет их в хаб для дальнейшей обработки.

**Обработчики** (handlers.go): Обрабатывают HTTP-запросы. В данном случае, ServeWs устанавливает WebSocket-соединение, ожидает инициализационное сообщение (INIT), регистрирует клиента в хабе и запускает горутины для чтения и записи сообщений.



