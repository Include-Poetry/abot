package main

import (
	"flag"
	"log"
	"math/rand"
	"os"

	"github.com/avabot/ava/shared/datatypes"
	"github.com/avabot/ava/shared/pkg"
	"github.com/jmoiron/sqlx"
)

var port = flag.Int("port", 0, "Port used to communicate with Ava.")
var db *sqlx.DB

type Wish string

func main() {
	flag.Parse()
	trigger := &datatypes.StructuredInput{
		Commands: []string{"wish"},
	}
	db = connectDB()
	p, err := pkg.NewPackage("wish", *port, trigger)
	if err != nil {
		log.Fatalln("creating package", p.Config.Name, err)
	}
	wish := new(Wish)
	if err := p.Register(wish); err != nil {
		log.Fatalln("registering package ", err)
	}
}

func (p *Wish) Run(m *datatypes.Message, respMsg *datatypes.ResponseMsg) error {
	resp := m.NewResponse()
	q := `INSERT INTO wishes (userid, sentence) VALUES ($1, $2)`
	_, err := db.Exec(q, m.User.ID, m.Input.Sentence)
	if err != nil {
		return err
	}
	n := rand.Intn(5)
	switch n {
	case 0:
		resp.Sentence = "Your wish is my command!"
	case 1:
		resp.Sentence = "I'll make some calls."
	case 2:
		resp.Sentence = "I hope to start doing that soon, too."
	case 3:
		resp.Sentence = "Roger that!"
	case 4:
		resp.Sentence = "I wish I could do that now, too. Soon, I hope."
	}
	return pkg.SaveResponse(respMsg, resp)
}

func (p *Wish) FollowUp(m *datatypes.Message,
	respMsg *datatypes.ResponseMsg) error {
	return pkg.SaveResponse(respMsg, m.NewResponse())
}

func connectDB() *sqlx.DB {
	log.Println("connecting to db")
	var db *sqlx.DB
	var err error
	if os.Getenv("AVA_ENV") == "production" {
		db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	} else {
		db, err = sqlx.Connect("postgres",
			"user=egtann dbname=ava sslmode=disable")
	}
	if err != nil {
		log.Println("err: could not connect to db", err)
	}
	return db
}
