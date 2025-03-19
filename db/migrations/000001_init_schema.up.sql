CREATE TABLE `products` (
	`id` INTEGER NOT NULL AUTO_INCREMENT UNIQUE,
	`name` VARCHAR(255) NOT NULL,
	`image` VARCHAR(255) NOT NULL,
	`category` VARCHAR(255) NOT NULL,
	`description` TEXT,
	`rating` INTEGER NOT NULL,
	`num_reviews` INTEGER NOT NULL DEFAULT 0,
	`price` DECIMAL(10,2) NOT NULL,
	`count_in_stock` INTEGER NOT NULL,
	`created_at` DATETIME DEFAULT now(),
	`updated_at` DATETIME,
	PRIMARY KEY(`id`)
);


CREATE TABLE `orders` (
	`id` INTEGER NOT NULL AUTO_INCREMENT UNIQUE,
	`payment_method` VARCHAR(255) NOT NULL,
	`tax_price` DECIMAL(10,2) NOT NULL,
	`shipping_price` DECIMAL(10,2) NOT NULL,
	`total_price` DECIMAL(10,2) NOT NULL,
	`created_at` DATETIME DEFAULT now(),
	`updated_at` DATETIME,
	PRIMARY KEY(`id`)
);


CREATE TABLE `order_items` (
	`id` INTEGER NOT NULL AUTO_INCREMENT UNIQUE,
	`order_id` INTEGER NOT NULL,
	`product_id` INTEGER NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`quantity` INTEGER NOT NULL,
	`image` VARCHAR(255) NOT NULL,
	`price` INTEGER NOT NULL,
	PRIMARY KEY(`id`)
);


ALTER TABLE `order_items`
ADD FOREIGN KEY(`order_id`) REFERENCES `orders`(`id`)
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `order_items`
ADD FOREIGN KEY(`product_id`) REFERENCES `products`(`id`)
ON UPDATE NO ACTION ON DELETE NO ACTION;