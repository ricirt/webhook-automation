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


INSERT INTO messages (content, phone_number, is_sent, created_at, updated_at) 
VALUES 
    ('Merhaba, bu bir test mesajıdır.', '+905551111111', false, NOW(), NOW()),
    ('İkinci test mesajı.', '+905552222222', false, NOW(), NOW()),
    ('Üçüncü test mesajı.', '+905553333333', false, NOW(), NOW()),
    ('Dördüncü test mesajı.', '+905554444444', false, NOW(), NOW()),
    ('Beşinci test mesajı.', '+905555555555', false, NOW(), NOW()),
    ('Altıncı test mesajı.', '+905556666666', false, NOW(), NOW()),
    ('Yedinci test mesajı.', '+905557777777', false, NOW(), NOW()),
    ('Sekizinci test mesajı.', '+905558888888', false, NOW(), NOW()),
    ('Dokuzuncu test mesajı.', '+905559999999', false, NOW(), NOW()),
    ('Onuncu test mesajı.', '+905550000000', false, NOW(), NOW()),
    ('On birinci test mesajı.', '+905551111112', false, NOW(), NOW()),
    ('On ikinci test mesajı.', '+905552222223', false, NOW(), NOW()),
    ('On üçüncü test mesajı.', '+905553333334', false, NOW(), NOW()),
    ('On dördüncü test mesajı.', '+905554444445', false, NOW(), NOW()),
    ('On beşinci test mesajı.', '+905555555556', false, NOW(), NOW()),
    ('On altıncı test mesajı.', '+905556666667', false, NOW(), NOW()),
    ('On yedinci test mesajı.', '+905557777778', false, NOW(), NOW()),
    ('On sekizinci test mesajı.', '+905558888889', false, NOW(), NOW()),
    ('On dokuzuncu test mesajı.', '+905559999990', false, NOW(), NOW()),
    ('Yirminci test mesajı.', '+905550000001', false, NOW(), NOW());
