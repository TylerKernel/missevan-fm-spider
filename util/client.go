package util

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	channels []*Channel
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) AddChannel(channel *Channel) {
	c.channels = append(c.channels, channel)
	return
}

func (c *Client) getChannel() (*Channel, error) {
	if len(c.channels) == 0 {
		return nil, errors.New("client channel is empty")
	}
	channelIndex := 0
	for index, channel := range c.channels {
		if channel.reqCount >= c.channels[channelIndex].reqCount {
			channelIndex = index
		}
	}
	channel := c.channels[channelIndex]
	if !channel.isCookiesValid() {
		err := channel.updateCookies()
		if err != nil {
			return nil, err
		}
	}
	return channel, nil
}

func (c *Client) Get(url string) (resp *resty.Response, err error) {
	channel, err := c.getChannel()
	if err != nil {
		return
	}
	return channel.NewRequest().Get(url)
}

type Channel struct {
	*resty.Client
	sync.RWMutex
	reqCount    int64
	lastReqTime time.Time
}

func (channel *Channel) updateCookies() error {
	channel.Cookies = make([]*http.Cookie, 0)
	resp, err := channel.R().Get("https://fm.missevan.com/api/user/info")
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.New("update cookies failed")
	}
	channel.SetCookies(resp.Cookies())
	return nil
}

func (channel *Channel) isCookiesValid() bool {
	if len(channel.Cookies) == 0 {
		return false
	}
	for _, cookie := range channel.Cookies {
		err := cookie.Valid()
		if err != nil {
			return false
		}
	}
	return true
}

func (channel *Channel) NewRequest() *resty.Request {
	resp := channel.R()
	resp.SetCookies(channel.Cookies)
	return resp
}

func NewChannelWithLocalAddr(localAddr net.Addr) *Channel {
	return &Channel{
		Client: resty.NewWithLocalAddr(localAddr),
	}
}
