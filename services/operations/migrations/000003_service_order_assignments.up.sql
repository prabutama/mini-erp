CREATE TABLE service_order_assignments (
    id UUID PRIMARY KEY,
    service_order_id UUID NOT NULL,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    assigned_user_id UUID NOT NULL,
    assigned_by_user_id UUID NOT NULL,
    status TEXT NOT NULL,
    request_id TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX service_order_assignments_one_active_idx
    ON service_order_assignments (service_order_id)
    WHERE status = 'active';

CREATE INDEX service_order_assignments_user_idx ON service_order_assignments (business_id, assigned_user_id, status);
