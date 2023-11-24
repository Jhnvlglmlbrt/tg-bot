<!-- <p align="center">
<img src="https://pepy.tech/badge/rss-aggregator" alt="https://pepy.tech/project/rss-aggregator">
<img src="https://pepy.tech/badge/rss-aggregator/month" alt="https://pepy.tech/project/rss-aggregator">
<img src="https://img.shields.io/github/license/Jhnvlglmlbrt/rss-aggregator.svg" alt="https://github.com/Jhnvlglmlbrt/rss-aggregator/blob/master/LICENSE"> -->

# ⚙️ Tg-bot

#### # Создание телеграм бота, который позволяет:

- Сохранять ссылки
- Отправлять случайную ссылку из списка
- Выводить список ссылок

## ❗ Requirements

- Создать tg бота у @BotFather
- token

## 💿 Installation

```
go get 
```

<!-- ## 💻 Example -->

## 🪛 How to use?

```
make run host="api.telegram.org" token="your-token"   
```

Чтобы бот сохранил ссылку - надо просто её отправить.

- **/start**
- **/help**  
- **/rnd** - отправляет 1 случайную ссылку клиенту.
- **/list** - отправляет список ссылок клиенту.
- **/remove your_url** - удаляет ссылку из списка.

## Roadmap

- [x] Улучшенный вывод списка
- [x] Удаление ссылок
- [ ] <del>Добавление категорий</del>
- [ ] Уведомления
- [ ] Импорт и экспорт ссылок