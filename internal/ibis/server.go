package ibis

import (
	"fmt"
	"gitlab.tesonero-computers.ru/ibis/aaa/internal/data"
	"gitlab.tesonero-computers.ru/ibis/aaa/internal/model"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	pb "gitlab.tesonero-computers.ru/ibis/authproto/pkg/endpoint"
)

type Server struct{}


var ibisSessionMap map[string]model.IbisSession
var defSessStatus = 0


func (s *Server) AuthEndpoint(stream pb.AAA_AuthEndpointServer) error {
	for {

		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}


		params, err := auth(in.Params)


		newUuid := uuid.NewV4().String()

		ibisSessionMap = make(map[string]model.IbisSession)

		var sess model.IbisSession
		sess.Account = in.Params.Account
		sess.Uuid    = newUuid
		sess.Status  = defSessStatus

		ibisSessionMap[newUuid] = sess

		// fmt.Println(sessionMap, "session map")

		resp := &pb.AuthResponse{
			Params: params,
			Uuid:   newUuid,
		}


		authMessage := "authentication success"

		if err != nil {
			log.Error().Err(err).Msg("")
			resp.Error = err.Error()
			authMessage = "authentication failed"
		}

		log.Info().Msgf("account: %d, uuid: %s, " + authMessage, resp.Params.Account, resp.Uuid)
		// log.Info().Msgf("account: %d, uuid: %s, " + authMessage, resp.Params.Account, resp.Params.Rnd)

		if err = stream.Send(resp); err != nil {
			return err
		}
	}
}

func (s *Server) SessionEndpoint(stream pb.AAA_SessionEndpointServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if in.Open {

			sessIbis := ibisSessionMap[in.Uuid];

			if sessIbis.Status == defSessStatus {

				_ , err := sessionOpen(in.Uuid)

				if err != nil {
					fmt.Println(err)
				}

				//dt := time.Now()
				//currDate := dt.Format("01-02-2006 15:04:05")
				//account := ibisSessionMap[in.Uuid].Account
				//
				//_, err := data.IbisSessionAdd(in.Uuid, currDate, account)
				//if err != nil {
				//	log.Info().Msgf("Not  session save: %s", in.Uuid)
				//}
				//
				//fmt.Println(ibisSessionMap[in.Uuid].Account)
				//
				//log.Info().Msgf("Request to open session: %s", in.Uuid)
			}
		}


		if in.Close {

			sessionClose(in.Uuid)
		}


		if err = stream.Send(&pb.SessionResponse{
			Uuid: in.Uuid,
		}); err != nil {
			return err
		}
	}
}


func sessionOpen(uuid string) (int, error) {

	dt := time.Now()
	currDate := dt.Format("01-02-2006 15:04:05")
	account := ibisSessionMap[uuid].Account

	sessId , err := data.IbisSessionAdd(uuid, currDate, account)
	if err != nil {
		log.Info().Msgf("Not  session save: %s", uuid)
		return sessId, err
	}

	// fmt.Println(ibisSessionMap[uuid].Account)

	log.Info().Msgf("Request to open session: %s", uuid)
	return sessId, nil

}


func sessionClose(uuid string) {

	dt := time.Now()
	currDate := dt.Format("01-02-2006 15:04:05")
	account := ibisSessionMap[uuid].Account

	_ , err := data.IbisSessionAdd(uuid, currDate, account)

	if err != nil {
		fmt.Println(err)
	}

}
