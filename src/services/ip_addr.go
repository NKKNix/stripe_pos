package services

import (
	httpclients "go-fiber-template/src/infrastructure/httpclients"
)

type IPService struct {
	IpHttpClient httpclients.IpInterface
}

type IIpService interface {
	GetIp() (string, error)
}

func NewIpService(ipHC httpclients.IpInterface) IIpService {
	return &IPService{
		IpHttpClient: ipHC,
	}
}

func (sv *IPService) GetIp() (string, error) {
	return sv.IpHttpClient.GetIp()
}
