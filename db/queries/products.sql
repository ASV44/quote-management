-- name: CreateProducts :batchexec
INSERT INTO products(name, description, price, tax_rate, metadata)
VALUES($1, $2, $3, $4, $5);

-- name: GetProducts :many
SELECT *
FROM products p
WHERE p.name ILIKE '%' || @query::TEXT || '%'
  OR p.description ILIKE '%' || @query::TEXT || '%'
ORDER BY
    CASE WHEN @sortBy::TEXT = 'name' AND @sortOrder::TEXT = 'asc' THEN p.name END ASC,
    CASE WHEN @sortBy::TEXT = 'name' AND @sortOrder::TEXT = 'desc' THEN p.name END DESC,
    CASE WHEN @sortBy::TEXT = 'price' AND @sortOrder::TEXT = 'asc' THEN p.price END ASC,
    CASE WHEN @sortBy::TEXT = 'price' AND @sortOrder::TEXT = 'desc' THEN p.price END DESC,
    CASE WHEN @sortBy::TEXT = 'created_at' AND @sortOrder::TEXT = 'asc' THEN p.created_at END ASC,
    CASE WHEN @sortBy::TEXT = 'created_at' AND @sortOrder::TEXT = 'desc' THEN p.created_at END DESC
--     CASE
--     WHEN @sortOrder = 'asc' THEN
--         CASE WHEN @sortBy::TEXT = 'name' THEN p.name END
--         CASE WHEN @sortBy::TEXT = 'price' THEN p.price END
--         CASE WHEN @sortBy::TEXT = 'created_at' THEN p.created_at END
--     WHEN @sortOrder = 'desc' THEN
--         CASE WHEN @sortBy::TEXT = 'name' THEN p.name END DESC
--         CASE WHEN @sortBy::TEXT = 'price' THEN p.price END DESC
--         CASE WHEN @sortBy::TEXT = 'created_at' THEN p.created_at END DESC
--     END
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

-- name: UpdateProducts :batchexec
UPDATE products
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    tax_rate = COALESCE(sqlc.narg('tax_rate'), tax_rate),
    metadata = COALESCE(sqlc.narg('metadata'), metadata),
    updated_at = NOW()
WHERE id = @id;

-- name: DeleteProducts :exec
DELETE FROM products
WHERE id = ANY(@ids::INT[]);
