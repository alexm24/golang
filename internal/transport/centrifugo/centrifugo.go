package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

type Centrifugo struct {
	cfg models.CentrifugoConfig
}

func NewCentrifugo(cfg models.CentrifugoConfig) *Centrifugo {
	return &Centrifugo{cfg}
}

func (c *Centrifugo) GetToken(user api.SUsername) (api.SToken, error) {
	ttl := time.Now().Add(time.Hour * 10)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": *user.Username,
		"exp": ttl.Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.cfg.TokenHmacSecret))

	if err != nil {
		tokenEmpty := ""
		timeNow := time.Now()
		return api.SToken{Exp: &timeNow, Token: &tokenEmpty}, err
	}

	return api.SToken{Exp: &ttl, Token: &tokenString}, nil
}

func (c *Centrifugo) Publish(channel string, msg interface{}) error {
	cmd := models.Centrifugo{
		Method: "publish",
		Params: models.Params{
			Channel: channel,
			Data:    msg,
		},
	}
	byteCmd, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	authKey := fmt.Sprintf("apikey %s", c.cfg.APIKey)

	req, err := http.NewRequest("POST", c.cfg.Url, bytes.NewBuffer(byteCmd))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", authKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
