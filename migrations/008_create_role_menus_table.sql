-- =============================================
-- 8. ROLE_MENUS
-- Menu apa saja yang bisa diakses role tertentu
-- =============================================
CREATE TABLE role_menus (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    menu_id UUID NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    can_view BOOLEAN DEFAULT true,
    can_create BOOLEAN DEFAULT false,
    can_edit BOOLEAN DEFAULT false,
    can_delete BOOLEAN DEFAULT false,
    
    PRIMARY KEY (role_id, menu_id)
);

-- Indexes
CREATE INDEX idx_role_menus_role_id ON role_menus(role_id);
CREATE INDEX idx_role_menus_menu_id ON role_menus(menu_id);