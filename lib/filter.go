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
	"github.com/SmartEnergyPlatform/event-broker/util"

	"gopkg.in/mgo.v2/bson"
)

func StopFilter(processId string) (err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	_, err = collection.UpdateAll(FilterDeployment{ProcessId: processId}, bson.M{"$set": FilterDeployment{State: DEPLOYMENT_STOPPING}})
	return err
}

func CreateFilter(processid string, filterid string, filter Filter) (err error) {
	if filter.Topic == "" {
		filter.Topic = createFilterTopic(filter)
	}
	session, collection := getMongoFilterCollection()
	defer session.Close()
	pool, err := selectFilterPool()
	if err != nil {
		return err
	}
	return collection.Insert(FilterDeployment{FilterPool: pool, Filter: filter, FilterId: filterid, ProcessId: processid, State: DEPLOYMENT_STARTING})
}

func createFilterTopic(filter Filter) string {
	return util.Config.DefaultKafkaTopic //strings.Replace(filter.ServiceId, "#", "_", -1)
}

func GetAllFilter() (result []FilterDeployment, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	err = collection.Find(nil).All(&result)
	return
}

func GetFilter(id string) (result FilterDeployment, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	err = collection.Find(FilterDeployment{FilterId: id}).One(&result)
	return
}
