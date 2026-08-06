CREATE TABLE businesses (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    plan TEXT NOT NULL DEFAULT 'free',
    platform_notes TEXT NULL,
    suspended_at TIMESTAMPTZ NULL,
    timezone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE branches (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    address TEXT NULL,
    phone TEXT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (business_id, code)
);

CREATE TABLE employee_placements (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    business_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    position TEXT NOT NULL,
    employment_type TEXT NOT NULL,
    status TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
