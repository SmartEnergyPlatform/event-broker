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
	"github.com/SmartEnergyPlatform/amqp-wrapper-lib"
	"github.com/SmartEnergyPlatform/event-broker/util"
	"log"
)

var amqp *amqp_wrapper_lib.Connection

type MsgEvent struct {
	FilterId string `json:"filter_id,omitempty"`
	ShapeId  string `json:"shape_id"`
	Filter   Filter `json:"filter"`
	State    string `json:"state,omitempty" bson:"-"`
}

type AbstractProcess struct {
	ReceiveTasks  []MsgEvent     `json:"receive_tasks"`
	MsgEvents     []MsgEvent     `json:"msg_events"`
}

type DeploymentRequest struct {
	Process AbstractProcess `json:"process"`
	Svg     string          `json:"svg"`
}

type DeploymentCommand struct {
	Command    string          		 	`json:"command"`
	Id         string           		`json:"id"`
	Owner      string           		`json:"owner"`
	Deployment DeploymentRequest	`json:"deployment"`
	DeploymentXml string				`json:"deployment_xml"`
}

func InitEventSourcing()(err error){
	amqp, err = amqp_wrapper_lib.Init(util.Config.AmqpUrl, []string{util.Config.AmqpDeploymentTopic}, util.Config.AmqpReconnectTimeout)
	if err != nil {
		return err
	}

	//event delete
	err = amqp.Consume(util.Config.AmqpConsumerName + "_" +util.Config.AmqpDeploymentTopic, util.Config.AmqpDeploymentTopic, func(delivery []byte) error {
		command := DeploymentCommand{}
		err = json.Unmarshal(delivery, &command)
		if err != nil {
			log.Println("ERROR: unable to parse amqp event as json \n", err, "\n --> ignore event \n", string(delivery))
			return nil
		}
		log.Println("amqp receive ", string(delivery))
		switch command.Command {
		case "POST":
			log.Println("WARNING: deprecated event type", command)
			return nil
		case "PUT":
			return handleDeploymentEventsUpdate(command)
		case "DELETE":
			return handleDeploymentEventsDelete(command)
		default:
			log.Println("WARNING: unknown event type", string(delivery))
			return nil
		}
	})

	return err
}

func handleDeploymentEventsDelete(command DeploymentCommand) error {
	return StopFilter(command.Id)
}

func handleDeploymentEventsUpdate(command DeploymentCommand)error{
	//log.Println("DEBUG: handleDeploymentEventsUpdate() ", command)
	return deployEvents(append(command.Deployment.Process.MsgEvents, command.Deployment.Process.ReceiveTasks...), command.Id)
}

func CloseEventSourcing(){
	amqp.Close()
}

func deployEvents(events []MsgEvent, processid string) (err error) {
	//log.Println("DEBUG: deployEvents() ", events)
	for _, event := range events {
		log.Println("DEBUG: deployEvents().loop ", err)
		err = SetFilter(processid, event.FilterId, event.Filter)
		if err != nil {
			log.Println("ERROR: unable to set filter ", err)
			return
		}
	}
	return
}
