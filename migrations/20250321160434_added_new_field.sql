-- +goose Up
ALTER TABLE tokens ADD COLUMN contract_address TEXT;
ALTER TABLE tokens ADD COLUMN decimals INT;
UPDATE tokens SET contract_address = '';
UPDATE tokens SET decimals = 18;
ALTER TABLE tokens ALTER COLUMN contract_address SET NOT NULL;

-- +goose Down
ALTER TABLE tokens DROP COLUMN contract_address;
ALTER TABLE tokens DROP COLUMN decimals;