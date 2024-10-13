create table public.role (
    id              UUID PRIMARY KEY default uuid_generate_v4(),
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

create table public.account (
    id          UUID PRIMARY KEY default uuid_generate_v4(),
    email       text unique,
    username    text,
    password    text not null,
    role        UUID references public.role(id),
    verified    boolean default false ,
    active      boolean default true,
    created_at  timestamp,
    updated_at  timestamp
);

create table public.cart (
    id          UUID PRIMARY KEY references public.account(id),
    content     json,
    updated_at  timestamp
);

create table public.brand (
    id          UUID PRIMARY KEY default uuid_generate_v4(),
    name        text,
    logo        text,
    on_sale     boolean default false,
    active      boolean default true,
    owner       UUID references public.account(id),
    created_at  timestamp,
    updated_by  UUID,
    updated_at  timestamp
);

create table public.category (
    id          UUID PRIMARY KEY default uuid_generate_v4(),
    name        text,
    description text,
    featured    boolean default false,
    active      boolean default true,
    owner       UUID references public.account(id),
    created_at  timestamp,
    updated_by  UUID,
    updated_at  timestamp
);

create table public.product (
    id              UUID PRIMARY KEY default uuid_generate_v4(),
    name            text,
    image           json,
    price           int default 0,
    colour          json,
    brand           UUID references public.brand(id),
    categories      json,
    size            json,
    on_sale         boolean default false,
    sale_price      int,
    sale_percent    int,
    stock           int default 0,
    is_new          boolean default true,
    description     text,
    owner           UUID references public.account(id),
    active          boolean default true,
    created_at      timestamp,
    updated_by      UUID,
    updated_at      timestamp
);



