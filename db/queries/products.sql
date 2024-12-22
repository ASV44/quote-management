-- name: CreateProducts :batchexec
INSERT INTO products(name, description, price, tax_rate, metadata)
VALUES($1, $2, $3, $4, $5);

-- name: GetProducts :many
SELECT *
FROM products
WHERE name ILIKE '%' || @query::TEXT || '%'
  OR description ILIKE '%' || @query::TEXT || '%'
ORDER BY @sortBy::TEXT
LIMIT $1
OFFSET $2;

-- name: GetProductsTotalCount :one
SELECT COUNT(*)
FROM products
WHERE name ILIKE '%' || @query::TEXT || '%'
  OR description ILIKE '%' || @query::TEXT || '%';

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = @id;