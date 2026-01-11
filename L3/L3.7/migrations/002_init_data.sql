-- Insert test users with different roles
INSERT INTO users (id, email, name, role)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'admin@warehouse.local', 'Admin User', 'admin'),
    ('22222222-2222-2222-2222-222222222222', 'manager@warehouse.local', 'Manager User', 'manager'),
    ('33333333-3333-3333-3333-333333333333', 'viewer@warehouse.local', 'Viewer User', 'viewer'),
    ('44444444-4444-4444-4444-444444444444', 'auditor@warehouse.local', 'Auditor User', 'auditor')
ON CONFLICT (email) DO NOTHING;

-- Insert sample items
INSERT INTO items (id, name, sku, quantity, reserved_qty, location, created_by, updated_by)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Laptop', 'SKU-001', 50, 5, 'A-001', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Monitor', 'SKU-002', 100, 10, 'A-002', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Keyboard', 'SKU-003', 200, 20, 'B-001', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Mouse', 'SKU-004', 150, 15, 'B-002', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111')
ON CONFLICT (sku) DO NOTHING;
