package services

import (
	"errors"
	"fmt"
	"jobgolangcrawl/config"
	"jobgolangcrawl/models"
	"log"
	"net/smtp"
	"sort"
	"strings"
	"time"
)

type MailService struct {
	dtos   []*models.PostRequestDto
	config *config.Config
}

func NewMailService(dtos []*models.PostRequestDto, config *config.Config) *MailService {
	sort.Slice(dtos, func(i, j int) bool {
		return strings.Compare(dtos[i].SiteName, dtos[j].SiteName) < 0
	})
	return &MailService{
		dtos:   dtos,
		config: config,
	}
}

func (s *MailService) SendMail() error {
	mailStart := time.Now()
	// Gmail SMTP 서버 설정
	smtpHost := s.config.Mail.SmtpHost
	smtpPort := s.config.Mail.SmtpPort

	// 발신자 정보
	sender := s.config.Mail.Sender
	password := s.config.Mail.Password

	// 수신자 정보
	receiver := s.config.Mail.Receiver

	// 메일 내용
	if len(s.dtos) <= 0 {
		return errors.New("No posts found")
	}
	subject := "Subject: JobGolang Finds New JobPosts! [" + time.Now().Format("2006-01-02 15:04") + "]\n\n"
	body := strings.Builder{}
	curSiteName := ""
	for _, dto := range s.dtos {
		if curSiteName != dto.SiteName {
			curSiteName = dto.SiteName
			body.WriteString(fmt.Sprintln(strings.ToUpper(curSiteName)))
		}
		body.WriteString(fmt.Sprintf("%s  |  %s\n", dto.CompanyName, dto.Title))
		body.WriteString(fmt.Sprintf("%s\n\n", dto.Url))
	}
	message := []byte(subject + "\n" + body.String())

	// 인증 설정
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// 메일 전송
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{receiver}, message)
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully!")
	mailEnd := time.Now()
	elapsed := mailEnd.Sub(mailStart)
	log.Println("Mail Process took: ", elapsed)
	return nil
}
