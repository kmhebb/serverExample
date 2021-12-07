package slack

import (
	"fmt"

	cloud "github.com/kmhebb/serverExample"
	"github.com/slack-go/slack"

	"github.com/kmhebb/serverExample/log"
)

const ChannelID = "C02H8BU6A9X"

func New(token string, logger log.Logger) *Service {
	c := slack.New(token, slack.OptionDebug(true))
	return &Service{
		c: c,
		l: logger,
	}
}

type Service struct {
	c *slack.Client
	l log.Logger
}

func (s Service) TestConnection() error {
	_, err := s.c.AuthTest()
	if err != nil {
		return fmt.Errorf("slack failed to connect: %w", err)
	}
	//fmt.Printf("slack auth test: %+v", resp)
	return nil
}

// func (s Service) NewCustomerAsync(ctx cloud.Context) {
// 	attachment := slack.Attachment{
// 		Fields: []slack.AttachmentField{
// 			slack.AttachmentField{
// 				Title: "To:",
// 				Value: fmt.Sprintf("<%s>%s", to, name),
// 			},
// 			slack.AttachmentField{
// 				Title: "CID:",
// 				Value: strconv.FormatInt(cid, 10),
// 			},
// 			slack.AttachmentField{
// 				Title: "Name:",
// 				Value: dlr,
// 			},
// 		},
// 	}
// 	if _, _, err := s.c.PostMessageContext(
// 		ctx.Ctx,
// 		ChannelID,
// 		slack.MsgOptionText("New Customer", false),
// 		slack.MsgOptionAttachments(attachment),
// 	); err != nil {
// 		//s.l.Log(ctx, log.Error, "Slack message failed", "err", err)
// 	}
// }

func (s Service) NewUserAsync(ctx cloud.Context, name, to, pass string) {
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "To:",
				Value: fmt.Sprintf("<%s>%s", to, name),
			},
			slack.AttachmentField{
				Title: "Password:",
				Value: pass,
			},
		},
	}
	_, _, err := s.c.PostMessageContext(
		ctx.Ctx,
		ChannelID,
		slack.MsgOptionText("New User", false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		fmt.Errorf("slack message for new user failed: %w", err)
		//s.l.Log(ctx, log.Error, "Slack message failed", "err", err)
	}
}

func (s Service) ResetPasswordAsync(ctx cloud.Context, name, to, token string) {
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "To:",
				Value: fmt.Sprintf("<%s>%s", to, name),
			},
			slack.AttachmentField{
				Title: "Token:",
				Value: token,
			},
		},
	}
	if _, _, err := s.c.PostMessageContext(
		ctx.Ctx,
		ChannelID,
		slack.MsgOptionText("Reset Password", false),
		slack.MsgOptionAttachments(attachment),
	); err != nil {
		//s.l.Log(ctx, log.Error, "Slack message failed", "err", err)
	}
}

func (s Service) NewPasswordAsync(ctx cloud.Context, name, to, pass string) {
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "To:",
				Value: fmt.Sprintf("<%s>%s", to, name),
			},
			slack.AttachmentField{
				Title: "Password:",
				Value: pass,
			},
		},
	}
	if _, _, err := s.c.PostMessageContext(
		ctx.Ctx,
		ChannelID,
		slack.MsgOptionText("New Password", false),
		slack.MsgOptionAttachments(attachment),
	); err != nil {
		//s.l.Log(ctx, log.Error, "Slack message failed", "err", err)
	}
}

func (s Service) ValidateEmailAsync(ctx cloud.Context, to, code string) {
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "To:",
				Value: fmt.Sprintf("<%s>", to),
			},
			slack.AttachmentField{
				Title: "Code:",
				Value: code,
			},
		},
	}
	if _, _, err := s.c.PostMessageContext(
		ctx.Ctx,
		ChannelID,
		slack.MsgOptionText("Validate Email", false),
		slack.MsgOptionAttachments(attachment),
	); err != nil {
		//s.l.Log(ctx, log.Error, "Slack message failed", "err", err)
	}
}

func (s Service) Close() chan int {
	c := make(chan int, 1)
	c <- 1
	return c
}
