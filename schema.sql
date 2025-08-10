
CREATE TYPE food_type AS ENUM ('sandwhich', 'salad', 'soup');

CREATE TABLE foods (
  id BIGSERIAL PRIMARY KEY,
  food text NOT NULL,
  food_type food_type -- NULL if food can't be classified.  That should ofc never happen
);


