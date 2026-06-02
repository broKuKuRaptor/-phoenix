# Phoenix

1) Надо создать еще один микросервис в этом проекте.
- Название `currency`, расположение `cmd/currency` и `internal/currency`. 
- В нем так же будет основной сервис `CurrencyService` унаследованный от `common.BaseService`
- В нем не будет создаваться экземпляр `apiserver.ApiServer` пока.
- От будет обмениваться сообщениями AMQP Protocol по каналам RabbitMQ ( в частности с сервисом `accounts`), паттерн Publish/Subscriber. 
2) В сервис `accounts` так же добавит функционал обмена сообщениями AMQP Protocol по каналам RabbitMQ, паттерн Publish/Subscriber.
3) Использовать пакет `github.com/rabbitmq/amqp091-go`, создать обобщенный AMQP broker, код вынести в `internal/broker`.
4) Реализовать следующий алгоритм:
- Сервис `accounts` в методе AccountsService.updateCurrenciesSupportStatus() посылает сообщение
`ping:CurrenciesSupportStatus`.
- Все запущенные сервисы `currency` отвечают на это сообщение сообщением `pong:CurrenciesSupportStatus`. `pong:CurrenciesSupportStatus` дожен иметь тело с полями uid - уникадбный идентификатор инстанса (например адрес инстанса в памяти) и ts - метка времени.
- Сервис `accounts` получает эти сообщения и печатает в консоль.
- Прослушивание сообщений сделать в гороутинах прерываемых через контекст, при вызове метода `BaseService.Shutdown`.
- Максимально использовать имеющийся код. Для работы с AMQP Protocol взять в качестве примера
  [показанные в руководстве](https://www.rabbitmq.com/tutorials/tutorial-three-go#publishsubscribe) фрагмены кода.


Напиши как понял задачу. Задавай вопросы. Приступаешь к кодированию только после как согласуем план.