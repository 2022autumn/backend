package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
	"time"
)

//func QueryApplicationByAuthor(author_id string) (submit database.Application, notFound bool) {
//	err := global.DB.Where("author_id = ? AND status = 1", author_id).First(&submit).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return submit, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		panic(err)
//	} else {
//		return submit, false
//	}
//}
//func QueryUserIsScholar(user_id uint64) (submit database.Application, notFound bool) {
//	err := global.DB.Where("user_id = ? AND status = 1", user_id).First(&submit).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return submit, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		panic(err)
//	} else {
//		return submit, false
//	}
//}
//func CreateApplication(submit *database.Application) (err error) {
//	if err = global.DB.Create(&submit).Error; err != nil {
//		return err
//	}
//	return nil
//}
//func GetApplicationByID(application_id uint64) (application database.Application, notFound bool) {
//	err := global.DB.First(&application, application_id).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return application, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return application, true
//	} else {
//		return application, false
//	}
//}

//func MakeUserScholar(user database.User, application database.Application) {
//	user.Email = application.Email
//	user.AuthorName = application.AuthorName
//	user.UserType = 1
//	user.Fields = application.Fields
//	user.AuthorID = application.AuthorID
//	err := global.DB.Save(&user).Error
//	if err != nil {
//		panic(err)
//	}
//}

//	func QueryAllSubmit() (application []database.Application) {
//		global.DB.Find(&application)
//		return application
//	}
//
//	func QueryUncheckedSubmit() (applications []database.Application, notFound bool) {
//		err := global.DB.Where("status = ?", 0).Find(&applications).Error
//		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//			return applications, true
//		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//			panic(err)
//		} else {
//			return applications, false
//		}
//	}
func GetAppsByUserID(userID uint64, status int) (applications []database.Application, err error) {
	err = global.DB.Where("user_id = ? AND status = ?", userID, status).Find(&applications).Error
	return
}
func GetAppByID(appID uint64) (application database.Application, notFound bool) {
	err := global.DB.First(&application, appID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return application, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return application, false
	}
}
func GetApps(status int) (applications []database.Application, err error) {
	err = global.DB.Where("status = ?", status).Find(&applications).Error
	return
}

func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": EMAIL_USER,
		"pass": EMAIL_PASSWORD,
		"host": EMAIL_HOST,
		"port": EMAIL_PORT,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "ishare")) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", mailTo...)                                     //发送给多个用户
	m.SetHeader("Subject", subject)                                  //设置邮件主题
	m.SetBody("text/html", body)                                     //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

func SendVerifyCode(email string, code string) (err error) {
	subject := "来自ishare的申请验证码"
	// 邮件正文
	mailTo := []string{
		email,
	}
	body := "尊敬的用户您好，欢迎使用ishare学术交流平台，您的申请验证码是:\n"
	body += code + "\n"

	err = SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send code fail")
		//panic(err)
		return err
	}
	fmt.Println("send code successfully")
	return nil
}

func CreateVerifyCodeRecode(userID uint64, code string, email string) (err error) {
	rec := database.VerifyCode{
		Code:    code,
		UserID:  userID,
		Email:   email,
		GenTime: time.Now(),
	}
	if err = global.DB.Create(&rec).Error; err != nil {
		return err
	}
	return nil
}

func CheckVerifyCode(userID uint64, code int, email string) (rec database.VerifyCode, notFound bool) {
	rec = database.VerifyCode{}
	err := global.DB.Where("user_id = ? AND code = ? AND email = ?", userID, code, email).First(&rec).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return rec, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return rec, false
	}
}

func GetVerifiedUser() (num int) {
	verified := make([]database.User, 0)
	global.DB.Where("verified = ?", 1).Find(&verified)
	return len(verified)
}
