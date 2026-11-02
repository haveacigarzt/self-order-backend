CREATE TABLE IF NOT EXISTS permissions (
  id bigserial PRIMARY KEY,
  code text NOT NULL
);
CREATE TABLE IF NOT EXISTS roles_permissions (
  role_id bigint NOT NULL REFERENCES roles ON DELETE CASCADE,
  permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
  ('produk:read'),
  ('produk:write');