## Ghost Orange

This is basicly a storage that can store various types of data like **`text data, binary data, credentials and bank card info`**. It provides simple tui that access server via http or https.

![general scheme](./assets/data_flow_scheme.svg)

Сделано и работает:
Аутентификация и хранение сессий
Клиентский интерфейс
Хранилище postgeSQL
Хранение данных учетных записей в зашифрованном виде

Требует доработки:
Формы добавление и изменения карт, файлов и текста
Механизм обновления хранимых данных в БД
MakeFile
Вывод данных карты после ввода короткого кода
Заменить кастомный gzip middleware
