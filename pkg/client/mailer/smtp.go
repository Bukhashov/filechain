package mailer

import (
	"fmt"
	"net/smtp"
	"crypto/tls"
	"bytes"
	"html/template"
	
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/configs"
)

type Mailer interface {
	ParseTemplate(path string, data interface{})(err error)
	Send()(err error)
}

type mailer struct {
	config 	configs.Smtp
	logger 	*logging.Logger
	to		string
	subject string
	body 	string
}

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func(m *mailer) ParseTemplate(path string, data interface{})(err error) {
	t, err := template.ParseFiles(path); if err != nil {
		m.logger.Info(fmt.Sprintf("Mailer send error parse template %s", err))
		return err
	}

	buffer := new(bytes.Buffer)

	err = t.Execute(buffer, data); if err != nil {
		m.logger.Info(fmt.Sprintf("Mailer send error parse template t.ex %s", err))
		return err
	}

	m.body = buffer.String()
	
	return nil
}



func(m *mailer) Send()(err error) {
	auth := smtp.PlainAuth("", m.config.Mail, m.config.Password, m.config.Server)
	
	connect, err := tls.Dial("tcp",  m.config.Server+":"+ m.config.Port, 
		&tls.Config{ 
			ServerName: m.config.Server,
			},
		); 
		if err != nil {
			m.logger.Info(fmt.Sprintf("Mailer send error TSL Dail %s", err))
			return err
		}
	
	client, err := smtp.NewClient(connect, m.config.Server); if err != nil {
		m.logger.Info(fmt.Sprintf("Mailer send error new clinet %s", err))
		return err
	}
	
	massage := "To: " + m.to + "\r\nSubject: " + m.subject + "\r\n" + mime + "\r\n" + m.body

	// AUTH.
	err = client.Auth(auth); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error clinet auth %s", err))
		return err
	}
	// FROM.				
	err = client.Mail(m.config.Mail); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error clinet mail %s", err))
		return err
	}		
	// TO
	err = client.Rcpt(m.to); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error clinet rcpt %s", err))
		return err
	}

	w, _ := client.Data()
	// BODY MASSAGE
	_, err = w.Write([]byte(massage)); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error write file %s", err))
		return err
	}			

	err = w.Close(); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error clinet close %s", err))
		return err
	}
	err = client.Quit(); if err != nil { 
		m.logger.Info(fmt.Sprintf("Mailer send error clinet quit %s", err))
		return err
	}

	return nil
}

func NewMailer(config configs.Smtp, logger *logging.Logger, to, subject string) Mailer {
	return &mailer{
		config: config,
		logger: logger,
		to:		to,
		subject: subject,
	}
}