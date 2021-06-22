Демон - gRPC сервис, поднимает Bidirectional streaming RPC.

## Описание
Принимает от клиента поток 2d точек в формате:

**(X:0...10 000, Y:0...10 000) uint32**

По запросу в другой поток отдает три ближайшие точки-соседа от заданной.
Для поиска используется простая реализация кривой Мортона (Z-order curve)

### Конфигурация
Через командную строку задается порт сервера gRPC. По умолчанию порт: 8080

`app-name -port 8080`

### Сборка и запуск
```
git clone https://github.com/Haba1234/coordinateStoradge.git
cd coordinateStoradge
make build
make run
```
### Docker
```
sudo docker build -t service .
sudo docker run -d --name service -p 8080:8080 service
```