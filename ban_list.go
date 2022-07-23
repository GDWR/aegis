package main

import (
	"net"
)

type BanList struct {
	banMap map[string]bool
}

func NewBanList() *BanList {
	return &BanList{
		banMap: make(map[string]bool),
	}
}

func (b BanList) IsBanned(addr string) bool {
	_, banned := b.GetBannedAddr(addr)
	return banned
}

func (b BanList) GetBannedAddr(addr string) (*string, bool) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	names, err := net.LookupAddr(host)
	if err != nil {
		names = make([]string, 0)
	}

	for _, name := range append(names, host) {
		banned := b.banMap[name]

		if banned {
			return &name, true
		}
	}

	return nil, false
}

func (b BanList) AddBan(addr string) {
	b.banMap[addr] = true
}

func (b BanList) RemoveBan(addr string) {
	resolvedAddr, _ := b.GetBannedAddr(addr)

	if resolvedAddr != nil {
		b.banMap[*resolvedAddr] = false
	}
}
