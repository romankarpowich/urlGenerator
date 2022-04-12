CREATE TABLE IF NOT EXISTS shorts (
                                      id BIGSERIAL,
                                      short varchar(10),
    source varchar(255),
    CONSTRAINT shorts_pkey PRIMARY KEY (id)
    );

select nextval('shorts_id_seq');