drop table if exists posts;
drop table if exists sites;

CREATE TABLE `sites` (
                         `id`	BIGINT AUTO_INCREMENT PRIMARY KEY,
                         `name`	VARCHAR(100)	NOT NULL,
                         `created_at`	DATETIME(3) NULL,
                         `updated_at`	DATETIME(3)	NULL,
                         `deleted_at`	DATETIME(3)	NULL
);

CREATE TABLE `posts` (
                         `id`	BIGINT	AUTO_INCREMENT PRIMARY KEY,
                         `url`	VARCHAR(255)	NOT NULL,
                         `post_id`	VARCHAR(50) NULL,
                         `title`	VARCHAR(100)	NULL,
                         `company_name`	VARCHAR(100)	NULL,
                         `created_at`	DATETIME(3) NULL,
                         `updated_at`	DATETIME(3)	NULL,
                         `deleted_at`	DATETIME(3)	NULL,
                         `site_id`	BIGINT	NOT NULL,
                         FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE ON UPDATE CASCADE
);
