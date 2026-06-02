-- =============================================
-- SEED DATA - Initial data untuk P0
-- =============================================

-- 1. Seed Roles
INSERT INTO roles (id, name, display_name, description, is_system, is_active) VALUES
(gen_random_uuid(), 'admin', 'Administrator', 'Full access to everything', true, true),
(gen_random_uuid(), 'supervisor', 'Supervisor', 'Can approve requests and view all data', true, true),
(gen_random_uuid(), 'warehouse_operator', 'Warehouse Operator', 'Can create inventory transactions', true, true),
(gen_random_uuid(), 'viewer', 'Viewer', 'Read-only access', true, true);

-- 2. Seed Admin User
-- Password: admin123 (bcrypt hash)
INSERT INTO users (id, email, username, password, name, role, is_active) VALUES (
    gen_random_uuid(),
    'admin_user_m@yopmail.com',
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr/.cZqGYdHcJg3kE8x5XZpQqXQGqLq',
    'Super Admin',
    'admin',
    true
);

-- 3. Seed Permissions
INSERT INTO permissions (id, resource, action, name, description) VALUES
-- Items
(gen_random_uuid(), 'items', 'create', 'Create Items', 'Can create new items'),
(gen_random_uuid(), 'items', 'read', 'View Items', 'Can view items list and detail'),
(gen_random_uuid(), 'items', 'update', 'Update Items', 'Can edit existing items'),
(gen_random_uuid(), 'items', 'delete', 'Delete Items', 'Can delete items'),
-- Categories
(gen_random_uuid(), 'categories', 'create', 'Create Categories', 'Can create new categories'),
(gen_random_uuid(), 'categories', 'read', 'View Categories', 'Can view categories'),
(gen_random_uuid(), 'categories', 'update', 'Update Categories', 'Can edit categories'),
(gen_random_uuid(), 'categories', 'delete', 'Delete Categories', 'Can delete categories'),
-- UOM
(gen_random_uuid(), 'uom', 'create', 'Create UOM', 'Can create units of measurement'),
(gen_random_uuid(), 'uom', 'read', 'View UOM', 'Can view units of measurement'),
(gen_random_uuid(), 'uom', 'update', 'Update UOM', 'Can edit units of measurement'),
(gen_random_uuid(), 'uom', 'delete', 'Delete UOM', 'Can delete units of measurement'),
-- Transactions
(gen_random_uuid(), 'transactions', 'create', 'Create Transactions', 'Can create inventory transactions'),
(gen_random_uuid(), 'transactions', 'read', 'View Transactions', 'Can view transactions'),
(gen_random_uuid(), 'transactions', 'approve', 'Approve Transactions', 'Can approve transactions'),
-- Suppliers
(gen_random_uuid(), 'suppliers', 'create', 'Create Suppliers', 'Can create suppliers'),
(gen_random_uuid(), 'suppliers', 'read', 'View Suppliers', 'Can view suppliers'),
(gen_random_uuid(), 'suppliers', 'update', 'Update Suppliers', 'Can edit suppliers'),
(gen_random_uuid(), 'suppliers', 'delete', 'Delete Suppliers', 'Can delete suppliers'),
-- Warehouses
(gen_random_uuid(), 'warehouses', 'create', 'Create Warehouses', 'Can create warehouses'),
(gen_random_uuid(), 'warehouses', 'read', 'View Warehouses', 'Can view warehouses'),
(gen_random_uuid(), 'warehouses', 'update', 'Update Warehouses', 'Can edit warehouses'),
(gen_random_uuid(), 'warehouses', 'delete', 'Delete Warehouses', 'Can delete warehouses'),
-- User Management
(gen_random_uuid(), 'users', 'create', 'Create Users', 'Can create users'),
(gen_random_uuid(), 'users', 'read', 'View Users', 'Can view users'),
(gen_random_uuid(), 'users', 'update', 'Update Users', 'Can update users'),
(gen_random_uuid(), 'users', 'delete', 'Delete Users', 'Can delete users'),
(gen_random_uuid(), 'roles', 'manage', 'Manage Roles', 'Can manage roles and permissions');

-- 4. Assign Permissions ke Roles (untuk P1 nanti)
-- Admin: dapat semua permission (INSERT ALL)
-- Caranya nanti di P1 akan assign via query

-- 5. Seed Menus (untuk P2 nanti)
INSERT INTO menus (id, name, icon, path, sort_order, is_active) VALUES
(gen_random_uuid(), 'Dashboard', 'home', '/dashboard', 1, true),
(gen_random_uuid(), 'Master Data', 'database', '#', 2, true),
(gen_random_uuid(), 'Items', 'package', '/items', 1, true),
(gen_random_uuid(), 'Categories', 'folder', '/categories', 2, true),
(gen_random_uuid(), 'UOM', 'ruler', '/uom', 3, true),
(gen_random_uuid(), 'Suppliers', 'truck', '/suppliers', 4, true),
(gen_random_uuid(), 'Warehouses', 'building', '/warehouses', 5, true),
(gen_random_uuid(), 'Transactions', 'refresh-cw', '/transactions', 3, true),
(gen_random_uuid(), 'Reports', 'bar-chart', '/reports', 4, true),
(gen_random_uuid(), 'User Management', 'users', '/users', 5, true);

-- Update parent_id untuk submenu
UPDATE menus SET parent_id = (SELECT id FROM menus WHERE name = 'Master Data') WHERE name IN ('Items', 'Categories', 'UOM', 'Suppliers', 'Warehouses');