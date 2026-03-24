ALTER TABLE kitchens
ADD COLUMN IF NOT EXISTS region TEXT NOT NULL DEFAULT 'global';

CREATE TABLE IF NOT EXISTS emission_factors (
    region TEXT PRIMARY KEY,
    grid_g_co2e_per_kwh NUMERIC NOT NULL,
    note TEXT
);

INSERT INTO emission_factors (region, grid_g_co2e_per_kwh, note)
VALUES ('global', 450, 'illustrative default'),
       ('EU', 350, 'approx EU average'),
       ('IN', 620, 'approx India grid')
ON CONFLICT (region) DO NOTHING;
