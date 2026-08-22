CREATE TABLE service_order_status_history (
    id UUID PRIMARY KEY,
    service_order_id UUID NOT NULL,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    from_status TEXT NOT NULL,
    to_status TEXT NOT NULL,
    changed_by_user_id UUID NULL,
    request_id TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX service_order_status_history_order_idx ON service_order_status_history (service_order_id, created_at);
