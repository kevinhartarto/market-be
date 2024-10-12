create table public.role (
    id              UUID PRIMARY KEY,
    name            text not null,
    can_view        boolean default false,
    can_add         boolean default false,
    can_edit        boolean default false,
    can_delete      boolean default false,
    can_buy         boolean default false,
    can_wishlist    boolean default false,
    is_admin        boolean default false,
    is_owner        boolean default false,
    deprecated      boolean default false
);

create table public.user (
    id          UUID PRIMARY KEY default gen_random_uuid(),
    email       text unique,
    username    text,
    password    text not null,
    role        UUID references public.role(id),
    verified    boolean default false ,
    active      boolean default true,
    created_at  timestamp,
    updated_at  timestamp
);



