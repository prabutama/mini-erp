CREATE TABLE resources (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    unit TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (business_id, branch_id, code)
);

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    resource_id UUID NOT NULL,
    movement_type TEXT NOT NULL,
    quantity NUMERIC NOT NULL,
    reason TEXT NULL,
    service_order_id UUID NULL,
    actor_user_id UUID NULL,
    request_id TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX stock_movements_resource_idx ON stock_movements (business_id, branch_id, resource_id);
