# buffer

В данной работе в качестве буфера данных выступает kafka
Для хранения оффсетов используется поле в структуре BufferUS
При завершении программы текущий оффсет сохраняется в файл .kafka_offset
При старте программы оффсет берется из файла .kafka_offset, если файла не существует, то он создается и оффсет устанавливается в 0

В качестве конфигурации используются переменные окружения
Для инициализации переменных окружения используется скрипт env.sh (т.к. задание учебное - чувствительные данные не скрыты)
Для разворачивания окружения кафки использовать docker-compose

В качестве примера записи данных в буфер приведена команда:
curl -i "http://localhost:8080/api/v1/facts/get?period_start=2024-12-01&period_end=2024-12-31&period_key=month&indicator_to_mo_id=227373&offset=0"
, где offset - это опциональный параметр, служащий для указания смещения от начала списка (используется в случае если в предыдущем запросе не все данные записались)

В качестве примера извлечения данных из буфера и записи их на сервер приведена команда:
curl -i "http://localhost:8080/api/v1/fact/save"