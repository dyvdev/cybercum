#	Для запуска нужен config.json:
айди бота(получаем у botfather):  
`"BotId": ""`  

Есть два режима:
- Случайные фразы из фиксированного набора без учета контекста:  
  `"EnablePhrases": false`  
  `"DefaultPhrases": []`

    если режим включен, то будет работать только он

- Генерация случайных фраз в ответ на сообщение в чате:  
`"EnableSemen": true`  
  бот пишет сообщения каждые "Ratio" сообещений в чате  
  `"Ratio": 50`  
  примерная длина максимального сообщения  
  `"Length": 50`  
  Название текстового файла, из которого бот возьмет базовые знания. Можно не указывать  
  `"DefaultDataFileName": ""`

Первый кум, на команды которого бот реагирует. Кум может добавлять новых кумов  
`"MainCum": ""`  
#	chats.json  
После того, как бот будет добавлен в чат, конфиг чата добавится в `chats.json`
#	Запуск:  
`go run main/main.go`  
