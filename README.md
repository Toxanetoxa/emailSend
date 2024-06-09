1. Определение интерфейсов и структур.
2. Реализация компонента отправки писем.
3. Добавление очереди и подключения к брокеру (без конкретной зависимости от Redis).
4. Добавление ограничения на количество отправленных писем.
5. Написание тестов с использованием моков.
6. Подготовка места для метрик.


тз
Компонент для отправки email через очередь
Нужен пакет, который реализует интерфейс компонента системы для отправки письма (емейла).
Он должен складывать их в очередь и отправлять, используя библиотеку для отправки емейлов.
Должна быть возможность установки лимита на кол-во отправленных писем (например, лимит в минуту).
Очередь нужно хранить в redis, компонент должен быть полностью stateless, чтобы забирать письма
из очереди и отправлять их могли сразу несколько экземпляров приложения. Прямой зависимости от redis
желательно избежать, что бы можно было заменить его на другие брокеры (rabbitmq, kafka).

Необходимо определить интерфейсы и покрыть всё тестами, используя мок для отправки емейлов.
Тесты должны покрыть все возможные кейсы.

Предусмотреть место для внедрения метрик в будущем (для prometheus),
чтобы считать, сколько отправлено, успешно\не успешно и тд. сами метрики пока не нужны

Делать покусочно, декомпозировав сначала задачу на более мелкие.
То есть итеративно, кусок сделан - можно показать. Но кусок должен быть логически завершенным и целостным.
Например, если не можешь сделать целиком все задание, но можешь сделать отдельные части, которые затем применить в проекте.