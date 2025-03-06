package services

import (
	"fmt"
	"jobgolangcrawl/models"
	"jobgolangcrawl/repositories"
	"time"
)

type PostService struct {
	PostRepository *repositories.PostRepository
	SiteRepository *repositories.SiteRepository
}

func NewPostService(
	postRepository *repositories.PostRepository,
	siteRepository *repositories.SiteRepository) *PostService {
	return &PostService{
		PostRepository: postRepository,
		SiteRepository: siteRepository,
	}
}

// CreateNewPosts: filtering dtos and insert data
func (s *PostService) CreateNewPosts(dtos []*models.PostRequestDto) ([]*models.PostRequestDto, error) {
	createNewPostsStart := time.Now()

	siteNameIdMap, siteIdNameMap, err := s.initSiteMap()

	err = s.insertNewSites(dtos, siteNameIdMap)
	if err != nil {
		return nil, err
	}

	siteNameIdMap, siteIdNameMap, err = s.initSiteMap()

	// Get all data from db
	sitePostMap, err := s.PostRepository.GetAllPostIdsWithSite()
	if err != nil {
		fmt.Println("Failed to get all posts: ", err)
		return nil, err
	}

	newDtos := s.filterNewPosts(dtos, sitePostMap, siteNameIdMap)

	// Insert data with site id
	fmt.Printf("New dtos count: %d\n", len(newDtos))
	err = s.PostRepository.InsertPosts(newDtos)
	if err != nil {
		fmt.Println("Failed to insert posts: ", err)
		return nil, err
	}

	result := s.convertToPostRequestDtos(newDtos, siteIdNameMap)

	createNewPostsEnd := time.Now()
	fmt.Printf("createNewPosts time: %v\n", createNewPostsEnd.Sub(createNewPostsStart))
	return result, nil
}

func (s *PostService) insertNewSites(dtos []*models.PostRequestDto, siteNameIdMap map[string]int) error {
	// Get New Sites from dtos
	newSites := make(map[string]struct{}, 0)
	for _, dto := range dtos {
		if _, ok := siteNameIdMap[dto.SiteName]; !ok {
			newSites[dto.SiteName] = struct{}{}
		}
	}
	fmt.Printf("Get new sites from dtos...\n  new sites count: %d\n", len(newSites))

	// Insert new sites
	if len(newSites) > 0 {
		err := s.SiteRepository.InsertSites(newSites)
		if err != nil {
			fmt.Println("Failed to insert new sites: ", err)
			return err
		}
	}
	fmt.Println("Insert sites succeeded")
	return nil
}

func (s *PostService) initSiteMap() (map[string]int, map[int]string, error) {
	// get all sites
	sites, err := s.SiteRepository.GetAllSites()
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("Get all sites...\n  sites count: %d\n", len(sites))

	siteNameIdMap := make(map[string]int)
	siteIdNameMap := make(map[int]string)
	for _, site := range sites {
		if _, ok := siteNameIdMap[site.Name]; !ok {
			siteNameIdMap[site.Name] = site.ID
		}
		if _, ok := siteIdNameMap[site.ID]; !ok {
			siteIdNameMap[site.ID] = site.Name
		}
	}
	return siteNameIdMap, siteIdNameMap, nil
}

func (s *PostService) filterNewPosts(dtos []*models.PostRequestDto, sitePostMap map[string]map[string]struct{}, siteNameIdMap map[string]int) []*models.PostInsertDto {
	// Filter new data
	newDtos := make([]*models.PostInsertDto, 0)
	for _, dto := range dtos {
		if _, ok := sitePostMap[dto.SiteName]; ok {
			if _, ok = sitePostMap[dto.SiteName][dto.PostId]; ok {
				continue
			}
		}

		siteId, ok := siteNameIdMap[dto.SiteName]
		if !ok {
			continue
		}
		newDtos = append(newDtos, models.ConvertToPostInsertDto(dto, siteId))
	}
	return newDtos
}

func (s *PostService) convertToPostRequestDtos(newDtos []*models.PostInsertDto, siteIdNameMap map[int]string) []*models.PostRequestDto {
	result := make([]*models.PostRequestDto, 0)
	for _, dto := range newDtos {
		siteName, ok := siteIdNameMap[dto.SiteId]
		if !ok {
			continue
		}
		result = append(result, models.ConvertToPostRequestDto(dto, siteName))
	}
	return result
}
