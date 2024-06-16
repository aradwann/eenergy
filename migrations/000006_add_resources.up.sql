CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    notes VARCHAR, 
    coords GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL UNIQUE, 
    name VARCHAR(255) NOT NULL,
    notes VARCHAR, 
    owner VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (location_id) REFERENCES locations(id),
    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE INDEX locations_coords_idx ON locations USING GIST (coords);
