package repositories

import (
	"database/sql"
	"jobgolangcrawl/models"
	"log"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) GetAllPostIdsWithSite() (map[string]map[string]struct{}, error) {
	sitePostsMap := make(map[string]map[string]struct{})

	rows, err := r.DB.Query(
		"select s.name, post_id " +
			"from posts p join sites s on p.site_id = s.id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var siteName string
		var postID string
		if err = rows.Scan(&siteName, &postID); err != nil {
			return nil, err
		}
		if _, ok := sitePostsMap[siteName]; !ok {
			sitePostsMap[siteName] = make(map[string]struct{})
		}
		sitePostsMap[siteName][postID] = struct{}{}
	}
	return sitePostsMap, nil
}

func (r *PostRepository) InsertPosts(newDtos []*models.PostInsertDto) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO posts (post_id, url, title, company_name, site_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, dto := range newDtos {
		_, err = stmt.Exec(dto.PostId, dto.Url, dto.Title, dto.CompanyName, dto.SiteId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
