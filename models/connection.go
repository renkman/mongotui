package models

type Connection struct {
	Host       string
	Port       string
	User       string
	Password   string
	Replicaset string
	tls        bool
}
