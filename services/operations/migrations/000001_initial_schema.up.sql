CREATE TABLE service_definitions (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    description TEXT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (business_id, code)
);

CREATE TABLE service_orders (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    service_definition_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX service_orders_business_branch_idx ON service_orders (business_id, branch_id);
