CREATE TABLE audit_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_version INT NOT NULL,
    producer TEXT NOT NULL,
    business_id UUID NOT NULL,
    branch_id UUID NULL,
    actor_id UUID NULL,
    request_id TEXT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audit_events_business_idx ON audit_events (business_id, occurred_at DESC);
CREATE INDEX audit_events_branch_idx ON audit_events (business_id, branch_id, occurred_at DESC);

CREATE TABLE operation_snapshots (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    branch_id UUID NULL,
    snapshot_date DATE NOT NULL,
    open_orders INT NOT NULL DEFAULT 0,
    in_progress_orders INT NOT NULL DEFAULT 0,
    completed_orders INT NOT NULL DEFAULT 0,
    cancelled_orders INT NOT NULL DEFAULT 0,
    resources_used NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX operation_snapshots_scope_date_idx
    ON operation_snapshots (business_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid), snapshot_date);
