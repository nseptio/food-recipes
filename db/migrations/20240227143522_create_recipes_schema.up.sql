CREATE TABLE Recipe (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(150) NOT NULL,
    cooking_time INT,
    serving int
);

CREATE TABLE Ingredient (
    id INT PRIMARY KEY AUTO_INCREMENT,
    recipe_id INT,
    name varchar(150) NOT NULL,
    quantity INT,
    unit varchar(50),
    FOREIGN KEY (recipe_id) REFERENCES Recipe(id)
);

CREATE TABLE Step (
    id INT PRIMARY KEY AUTO_INCREMENT,
    recipe_id INT,
    step_number INT,
    description TEXT,
    FOREIGN KEY (recipe_id) REFERENCES Recipe(id)
);

CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50),
    created_at BIGINT NOT NULL ,
    updated_at BIGINT NOT NULL
);

CREATE TABLE user_like_recipe (
    user_id INT,
    recipe_id INT,
    PRIMARY KEY (user_id, recipe_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (recipe_id) REFERENCES Recipe(id)
);
