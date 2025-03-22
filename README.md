# Webhook Automation Service

Bu servis, mesajları bir webhook URL'sine otomatik olarak gönderen ve durumlarını takip eden bir sistemdir.

## Özellikler

- Mesajları belirli aralıklarla (20 saniye) otomatik gönderme
- PostgreSQL'de mesaj durumu takibi
- Redis'te mesaj detayları önbellekleme
- REST API endpoints
- Swagger dokümantasyonu

## Teknolojiler

- Go 1.23
- PostgreSQL 15
- Redis 7
- Docker & Docker Compose
- Gin Web Framework
- GORM

## Kurulum

1. Projeyi klonlayın:
```bash
git clone https://github.com/your-username/webhook-automation.git
cd webhook-automation
```

2. Docker Compose ile çalıştırın:
```bash
docker compose up --build
```

## API Endpoints

### Mesaj Gönderme Servisi

- **POST /api/v1/messages/start**: Mesaj gönderme servisini başlatır
- **POST /api/v1/messages/stop**: Mesaj gönderme servisini durdurur
- **GET /api/v1/messages/sent**: Gönderilmiş mesajları listeler

### Swagger Dokümantasyonu

Swagger dokümantasyonuna erişmek için:
```
http://localhost:8080/swagger/index.html
```

## Veritabanı Şeması

### Messages Tablosu
```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    message_id VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Redis Önbellek

Mesaj detayları Redis'te şu formatta saklanır:
```json
{
    "message_id": "unique-message-id",
    "sent_at": "2024-03-22T23:09:07Z"
}
```
Key formatı: `message:{id}`

## Kontrol ve Test

1. Servisi başlatın:
```bash
curl -X POST http://localhost:8080/api/v1/messages/start
```

2. Gönderilen mesajları kontrol edin:
```bash
curl -X GET http://localhost:8080/api/v1/messages/sent
```

3. Redis'teki mesaj detaylarını kontrol edin:
```bash
docker exec -it webhook-automation-redis-1 redis-cli KEYS "message:*"
docker exec -it webhook-automation-redis-1 redis-cli GET "message:1"
```

4. PostgreSQL'deki mesajları kontrol edin:
```bash
docker exec -it webhook-automation-postgres-1 psql -U postgres -d insider_messages -c "SELECT * FROM messages WHERE is_sent = true;"
```

5. Servisi durdurun:
```bash
curl -X POST http://localhost:8080/api/v1/messages/stop
```