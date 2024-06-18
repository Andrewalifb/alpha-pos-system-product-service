CREATE TABLE pos_product_categories (
    category_id UUID PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_product_sub_categories (
    sub_category_id UUID PRIMARY KEY,
    sub_category_name VARCHAR(255) NOT NULL,
    category_id UUID REFERENCES pos_product_categories(category_id) NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_products (
    product_id UUID PRIMARY KEY,
    product_barcode_id VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    cost_price DECIMAL(10, 2),
    category_id UUID REFERENCES pos_product_categories(category_id) NOT NULL,
    sub_category_id UUID REFERENCES pos_product_sub_categories(sub_category_id) NOT NULL,
    stock_quantity INT NOT NULL,
    reorder_level INT,
    supplier_id UUID,
    product_description TEXT,
    active BOOLEAN DEFAULT TRUE,
    store_id UUID,
    branch_id UUID,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_suppliers (
    supplier_id UUID PRIMARY KEY,
    supplier_name VARCHAR(255) NOT NULL,
    contact_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_inventory_history (
    inventory_id UUID PRIMARY KEY,
    product_id UUID REFERENCES pos_products(product_id) NOT NULL,
    store_id UUID,
    date TIMESTAMP NOT NULL,
    quantity INT NOT NULL,
    branch_id UUID,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_promotions (
    promotion_id UUID PRIMARY KEY,
    product_id UUID REFERENCES pos_products(product_id) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    discount_rate DECIMAL(5, 2) NOT NULL,
    store_id UUID,
    branch_id UUID,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);
