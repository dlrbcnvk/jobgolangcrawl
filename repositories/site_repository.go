package repositories

import (
	"database/sql"
	"jobgolangcrawl/models"
	"log"
)

type SiteRepository struct {
	DB *sql.DB
}

func NewSiteRepository(db *sql.DB) *SiteRepository {
	return &SiteRepository{DB: db}
}

func (r *SiteRepository) GetAllSites() ([]*models.Site, error) {
	// Site 테이블의 전체 데이터 조회
	rows, err := r.DB.Query("SELECT id, name FROM sites")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 결과를 담을 슬라이스 선언
	var sites []*models.Site

	// 결과를 순회하며 구조체 슬라이스에 추가
	for rows.Next() {
		var site models.Site
		if err := rows.Scan(&site.ID, &site.Name); err != nil {
			log.Fatal(err)
		}
		sites = append(sites, &site)
	}

	// rows 순회 중 오류 확인
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return sites, nil
}

func (r *SiteRepository) InsertSites(sites map[string]struct{}) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		"INSERT INTO sites(name, created_at, updated_at) " +
			"VALUES(?, NOW(), NOW()) " +
			"ON DUPLICATE KEY UPDATE updated_at = NOW()")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for site := range sites {
		_, err = stmt.Exec(site)
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
