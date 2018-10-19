/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"encoding/json"
	"github.com/SmartEnergyPlatform/event-broker/util"
	"log"
	"net/http"
	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/julienschmidt/httprouter"
)

func ApiStart() {
	log.Println("start server on port: ", util.Config.ServerPort)
	httpHandler := getRoutes()
	corsHandler := cors.New(httpHandler)
	logger := logger.New(corsHandler, util.Config.LogLevel)
	if util.Config.DecodeUrlFix == "true" {
		urldecodefix := NewUrlDecodeMiddleWare(logger)
		log.Println(http.ListenAndServe(":"+util.Config.ServerPort, urldecodefix))
	} else {
		log.Println(http.ListenAndServe(":"+util.Config.ServerPort, logger))
	}
}

func getRoutes() *httprouter.Router {

	router := httprouter.New()

	router.GET("/filter", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filter, err := GetAllFilter()
		if err == nil {
			response.To(res).Json(filter)
		} else {
			log.Println("error on GetAllFilter(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/filter/:id", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		filter, err := GetFilter(id)
		if err == nil {
			response.To(res).Json(filter)
		} else {
			log.Println("error on GetFilter(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.POST("/filter/:processid/:filterid", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filterid := ps.ByName("filterid")
		processid := ps.ByName("processid")
		filter := Filter{}
		err := json.NewDecoder(r.Body).Decode(&filter)
		if err != nil {
			response.To(res).DefaultError(err.Error(), 400)
		}
		err = CreateFilter(processid, filterid, filter)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on CreateDeployments(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.DELETE("/filter/:processid", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		processid := ps.ByName("processid")
		err := StopFilter(processid)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on DeleteEvent(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/pool/assignments/:poolid", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		poolid := ps.ByName("poolid")
		assignments, err := GetFilterPoolAssignments(poolid)
		if err == nil {
			response.To(res).Json(assignments)
		} else {
			log.Println("error on GetPoolAssignments(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/pool/error/:poolid/:filterid/:err", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		poolid := ps.ByName("poolid")
		filterid := ps.ByName("filterid")
		err := ps.ByName("err")
		log.Println("ERROR: pool thrown error: ", poolid, filterid, err)
		response.To(res).Text("ok")
		log.Println("TODO: handle pool error")
	})

	router.GET("/pool/command/:poolid/:poolsize", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		poolid := ps.ByName("poolid")
		poolsize := ps.ByName("poolsize")
		command, err := GetFilterPoolCommand(poolid, poolsize)
		if err == nil {
			response.To(res).Json(command)
		} else {
			log.Println("error on GetPoolCommand(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.DELETE("/pool/:poolid", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		poolid := ps.ByName("poolid")
		err := RemoveFilterPool(poolid)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on RemoveFilterPool(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/swagger", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		swagger, err := GetSwagger()
		if err == nil {
			response.To(res).Text(swagger)
		} else {
			log.Println("error on GeteventPoolCommand(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	return router
}
