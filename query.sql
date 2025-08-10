-- name: SearchFood :many
SELECT * FROM foods
WHERE food LIKE $1
AND food_type = $2;

-- name: ListNullFoodType :many
SELECT * FROM foods
WHERE food_type IS NULL;

-- name: AddFood :exec
INSERT INTO foods (id, food, food_type)
VALUES ($1, $2, $3)
RETURNING *;
