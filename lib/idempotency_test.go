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
	"flag"
	"github.com/SmartEnergyPlatform/event-broker/util"
	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


/*
TODO: idempotent events

DELETE
	1. process does not exist
	2. process exists

PUT
	1. process does not exist
	2. process exists
*/

func TestIdempotency(t *testing.T) {
	closer, mongoport, _, err := testHelper_getMongoDependency()
	defer closer()
	if err != nil {
		t.Error(err)
		return
	}

	configLocation := flag.String("config", "../config.json", "configuration file")
	flag.Parse()

	err = util.LoadConfig(*configLocation)
	if err != nil {
		t.Error(err)
		return
	}

	util.Config.MongoUrl = "mongodb://localhost:"+mongoport

	httpServer := httptest.NewServer(getRoutes())
	defer httpServer.Close()

	err = updateKnownFilterPools("fp1")
	if err != nil {
		t.Error(err)
		return
	}

	err = testHelper_putProcess("1", "f1")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := http.Get(httpServer.URL + "/filter")
	if err != nil {
		t.Error(err)
		return
	}
	filter := []FilterDeployment{}
	err = json.NewDecoder(resp.Body).Decode(&filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(filter) != 1 || filter[0].ProcessId != "1" || filter[0].FilterId != "f1" {
		t.Error("unexpected result:", filter)
		return
	}

	err = testHelper_putProcess("1", "f2")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = http.Get(httpServer.URL + "/filter")
	if err != nil {
		t.Error(err)
		return
	}
	filter = []FilterDeployment{}
	err = json.NewDecoder(resp.Body).Decode(&filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(filter) != 1 || filter[0].ProcessId != "1" || filter[0].FilterId != "f2" {
		t.Error("unexpected result:", filter)
		return
	}

	err = testHelper_deleteProcess("1")
	if err != nil {
		t.Error(err)
		return
	}

	//mock filterpool
	_, err = GetFilterPoolAssignments("fp1")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = http.Get(httpServer.URL + "/filter")
	if err != nil {
		t.Error(err)
		return
	}
	filter = []FilterDeployment{}
	err = json.NewDecoder(resp.Body).Decode(&filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(filter) != 0 {
		t.Error("unexpected result:", filter)
		return
	}

	err = testHelper_deleteProcess("1")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = http.Get(httpServer.URL + "/filter")
	if err != nil {
		t.Error(err)
		return
	}
	filter = []FilterDeployment{}
	err = json.NewDecoder(resp.Body).Decode(&filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(filter) != 0 {
		t.Error("unexpected result:", filter)
		return
	}

	//mock filterpool
	_, err = GetFilterPoolAssignments("fp1")
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = http.Get(httpServer.URL + "/filter")
	if err != nil {
		t.Error(err)
		return
	}
	filter = []FilterDeployment{}
	err = json.NewDecoder(resp.Body).Decode(&filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(filter) != 0 {
		t.Error("unexpected result:", filter)
		return
	}

}


func testHelper_getMongoDependency() (closer func(), hostPort string, ipAddress string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() {}, "", "", err
	}
	log.Println("start mongodb")
	mongo, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		return func() {}, "", "", err
	}
	hostPort = mongo.GetPort("27017/tcp")
	err = pool.Retry(func() error {
		log.Println("try mongodb connection...")
		sess, err := mgo.Dial("mongodb://localhost:" + hostPort)
		if err != nil {
			return err
		}
		defer sess.Close()
		return sess.Ping()
	})
	return func() { mongo.Close() }, hostPort, mongo.Container.NetworkSettings.IPAddress, err
}

func testHelper_putProcess(vid string, filterId string) error {
	return handleDeploymentEventsUpdate(DeploymentCommand{
		Id:            vid,
		Command:       "PUT",
		Deployment: DeploymentRequest{
			Process:AbstractProcess{
				MsgEvents: []MsgEvent{
					{
						FilterId:filterId,
					},
				},
			},
		},
	})
}

func testHelper_deleteProcess(vid string) error {
	return handleDeploymentEventsDelete(DeploymentCommand{
		Id:      vid,
		Command: "DELETE",
	})
}
