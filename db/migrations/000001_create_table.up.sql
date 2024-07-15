CREATE TABLE IF NOT EXISTS item_service_items(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name varchar NOT NULL,
	description text,
	category_id uuid REFERENCES item_service_item_categories(id),
	condition varchar NOT NULL,
	swap_preference jsonb NOT NULL,
	owner_id uuid REFERENCES Auth_Service_Users(id)
	status varchar NOT NULL, 
	created_at timestamp DEFAULT CURRENT_TIMESTAMP
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP
	deleted_at timestamp DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS item_service_swaps (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	offered_item_id uuid REFERENCES Item_Service_Items(id),
	requested_item_id uuid REFERENCES Item_Service_Items(id),
	requester_id uuid REFERENCES Auth_Service_Users(id),
	owner_id uuid REFERENCES Auth_Service_Users(id),
	status varchar NOT NULL,
	message text, 
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamp DEFAULT null
);

CREATE TABLE IF NOT EXISTS item_service_item_categories (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name varchar NOT NULL,
	description text,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamp DEFAULT null
);

CREATE TABLE IF NOT EXISTS item_service_recycling_centers (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name varchar NOT NULL,
	address text NOT NULL,
	accepted_materials jsonb NOT NULL, 
	working_hours text,
	contact_number varchar NOT NULL, 
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamp DEFAULT null
);

CREATE TABLE IF NOT EXISTS item_service_recycling_submissions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	center_id uuid REFERENCES Item_Service_Recycling_Centers(id),
	user_id uuid REFERENCES Auth_Service_Users(id),
	items jsonb NOT NULL,
	eco_points_earned integer NOT NULL,  
	created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_service_ratings (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid REFERENCES Auth_Service_Users(id), 
	rater_id uuid REFERENCES Auth_Service_Users(id), 
	rating decimal((2, 1)), 
	comment text, 
	swap_id uuid REFERENCES Item_Service_Swaps(id),
	created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_service_eco_challenges (
	id uuid PRIMARY KEY  DEFAULT gen_random_uuid(),
	title varchar NOT NULL,
	description text NOT NULL, 
	start_date timestamp NOT NULL, 
	end_date timestamp NOT NULL,
	reward_points int, 
	created_at timestamp DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS item_service_challenge_participations (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	challenge_id uuid REFERENCES Item_Service_Eco_Challenges(id),
	user_id uuid REFERENCES Auth_Service_Users(id),
	status varchar DEFAULT "joined",
	recycled_items_count int DEFAULT 0,
	joined_at timestamp  DEFAULT CURRENT_TIMESTAMP,
	created_at timestamp  DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp  DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS item_service_eco_tips (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	title varchar NOT NULL,
	content text NOT NULL, 
	created_at timestamp  DEFAULT CURRENT_TIMESTAMP
);