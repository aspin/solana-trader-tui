package store

import (
	"encoding/json"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	"log"
	"os"
)

type App struct {
	Err      error
	UI       UI
	Settings Settings
	Provider *provider.GRPCClient
}

type UI struct {
	WindowWidth  int
	WindowHeight int
}

type Settings struct {
	AuthHeader        string
	PrivateKey        solana.PrivateKey
	PublicKey         solana.PublicKey
	OpenOrdersAddress solana.PublicKey
	Project           pb.Project
}

func NewFromFile(filename string) App {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("could not read file (%v): %v", filename, err)
		return App{}
	}

	m := struct {
		AuthHeader        string           `json:"authHeader"`
		PrivateKey        string           `json:"privateKey"`
		PublicKey         solana.PublicKey `json:"publicKey"`
		OpenOrdersAddress solana.PublicKey `json:"openOrdersAddress"`
		Project           string           `json:"project"`
	}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Printf("could not unmarshal json: %v, bytes: %v", err, string(b))
		return App{}
	}

	s := Settings{
		AuthHeader:        m.AuthHeader,
		PublicKey:         m.PublicKey,
		OpenOrdersAddress: m.OpenOrdersAddress,
	}
	s.PrivateKey, err = solana.PrivateKeyFromBase58(m.PrivateKey)
	if err != nil {
		log.Printf("could not deserialize private key: %v", err)
		return App{}
	}

	project, ok := pb.Project_value[m.Project]
	if !ok {
		log.Printf("could not deserialize project: %v, value: %v", err, m.Project)
		return App{}
	}
	s.Project = pb.Project(project)

	return App{
		Settings: s,
	}
}

func (a App) NeedsInit() bool {
	return a.Settings.AuthHeader == ""
}
