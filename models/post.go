package models

type PostRequestDto struct {
	PostId      string
	Url         string
	Title       string
	CompanyName string
	SiteName    string
}

func NewPostRequestDto(postId, url, title, companyName string, siteName string) *PostRequestDto {
	return &PostRequestDto{
		PostId:      postId,
		Url:         url,
		Title:       title,
		CompanyName: companyName,
		SiteName:    siteName,
	}
}

type PostInsertDto struct {
	PostId      string
	Url         string
	Title       string
	CompanyName string
	SiteId      int
}

func NewPostInsertDto(postId, url, title, companyName string, siteId int) *PostInsertDto {
	return &PostInsertDto{
		PostId:      postId,
		Url:         url,
		Title:       title,
		CompanyName: companyName,
		SiteId:      siteId,
	}
}

func ConvertToPostInsertDto(dto *PostRequestDto, siteId int) *PostInsertDto {
	return &PostInsertDto{
		PostId:      dto.PostId,
		Url:         dto.Url,
		Title:       dto.Title,
		CompanyName: dto.CompanyName,
		SiteId:      siteId,
	}
}

func ConvertToPostRequestDto(dto *PostInsertDto, name string) *PostRequestDto {
	return &PostRequestDto{
		PostId:      dto.PostId,
		Url:         dto.Url,
		Title:       dto.Title,
		CompanyName: dto.CompanyName,
		SiteName:    name,
	}
}
