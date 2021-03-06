package wechat

import (
	"github.com/num5/axiom"
	"github.com/KevinGong2013/wechat"
)

type weixin struct {
	axiom.BasicProvider
	bot *axiom.Robot
	users *axiom.UserMap
	wechat *wechat.WeChat
}

func NewWeChat(r *axiom.Robot) (axiom.Provider, error) {
	wx := new(weixin)

	wx.SetRobot(r)

	users := axiom.NewUserMap(r)

	wx.users = users

	wc, err := wechat.NewBot(nil)
	if err != nil {
		panic(err)
	}

	wx.wechat = wc

	return wx, nil
}

func (wx *weixin) Name() string {
	return "web wechat"
}

func (wx *weixin) Run() error {
	wx.wechat.Handle(`/login`, func(evt wechat.Event) {
		isSuccess := evt.Data.(int) == 1
		if isSuccess {
			log.Info(`登录成功......`)
		} else {
			log.Error(`登录失败......`)
		}
	})

	// 私聊
	wx.wechat.Handle(`/msg/solo`, func(evt wechat.Event) {
		wmsg := evt.Data.(wechat.EventMsgData)

		msg := &axiom.Message{
			FromUser: wx.newUsers(wx.wechat.ContactByUserName(wmsg.FromUserName)),
			Text: wmsg.Content,
		}

		err := wx.Receive(msg)

		if err != nil {
			log.Errorf(`receive msg error: %v`, err)
		}

	})

	// 微信群
	wx.wechat.Handle(`/msg/group`, func(evt wechat.Event) {
		wmsg := evt.Data.(wechat.EventMsgData)

		if wmsg.AtMe {

			msg := &axiom.Message{
				FromUser: wx.newUsers(wx.wechat.ContactByUserName(wmsg.FromUserName)),
				Text: wmsg.Content,
			}

			err := wx.Receive(msg)

			if err != nil {
				log.Errorf(`receive msg error: %v`, err)
			}
		}

	})

	wx.wechat.Go()

	return nil
}

func (wx *weixin) Close() error {
	return nil
}

func (wx *weixin) Receive(msg *axiom.Message) error {

	return wx.bot.Receive(msg)
}

func (wx *weixin) Send(res *axiom.Response, strings ...string) error {
	for _, str := range strings {
		err := wx.wechat.SendTextMsg(str, res.FromUserName())

		if err != nil {
			log.Errorf(`Send msg error: %v`, err)
		}
	}

	return nil
}

func (wx *weixin) Reply(res *axiom.Response, strings ...string) error {
	for _, str := range strings {
		err := wx.wechat.SendTextMsg(str, res.FromUserName())

		if err != nil {
			log.Errorf(`Send msg error: %v`, err)
		}
	}

	return nil
}

func (wx *weixin) newUsers(un *wechat.Contact) axiom.User {
	var user axiom.User
	user.ID = un.UserName
	user.Name = un.NickName
	user.Options = map[string]interface{} {
		"HeadImgURL": un.HeadImgURL,
		"HeadHash": un.HeadHash,
		"RemarkName": un.RemarkName,
		"DisplayName": un.DisplayName,
		"StarFriend": un.StarFriend,
		"Sex": un.Sex,
		"Signature": un.Signature,
		"VerifyFlag": un.VerifyFlag,
		"ContactFlag": un.ContactFlag,
		"HeadImgFlag": un.HeadImgFlag,
		"Province": un.Province,
		"City": un.City,
		"Alias": un.Alias,
		"EncryChatRoomID": un.EncryChatRoomID,
		"Type": un.Type,
		"MemberList": un.MemberList,
	}

	return user
}

func (wx *weixin) chatRoomMember(room_name string) (map[string]int, error) {

	stats := make(map[string]int)

	RoomContactList, err := wx.wechat.MembersOfGroup(room_name)
	if err != nil {
		return nil, err
	}

	man := 0
	woman := 0
	none := 0
	for _, v := range RoomContactList {

		member := wx.wechat.ContactByUserName(v.UserName)

		if member.Sex == 1 {
			man++
		} else if member.Sex == 2 {
			woman++
		} else {
			none++
		}

	}

	stats = map[string]int{
		"woman": woman,
		"man":   man,
		"none":  none,
	}

	return stats, nil
}

func init() {
	axiom.RegisterProvider(`wechat`, NewWeChat)
}
