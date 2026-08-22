CREATE TABLE workflows (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (business_id, name)
);

CREATE TABLE workflow_statuses (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL,
    business_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_initial BOOLEAN NOT NULL DEFAULT false,
    is_terminal BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workflow_id, code)
);

CREATE TABLE workflow_transitions (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL,
    business_id UUID NOT NULL,
    from_status_code TEXT NOT NULL,
    to_status_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workflow_id, from_status_code, to_status_code)
);

CREATE INDEX workflows_business_idx ON workflows (business_id);
CREATE INDEX workflow_statuses_workflow_idx ON workflow_statuses (workflow_id);
CREATE INDEX workflow_transitions_workflow_idx ON workflow_transitions (workflow_id);
