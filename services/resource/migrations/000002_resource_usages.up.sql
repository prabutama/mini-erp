CREATE TABLE resource_usages (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    service_order_id UUID NOT NULL,
    resource_id UUID NOT NULL,
    quantity NUMERIC NOT NULL,
    reason TEXT NULL,
    recorded_by_user_id UUID NULL,
    stock_movement_id UUID NOT NULL,
    request_id TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX resource_usages_order_idx ON resource_usages (business_id, service_order_id, created_at);
CREATE INDEX resource_usages_resource_idx ON resource_usages (business_id, branch_id, resource_id);
