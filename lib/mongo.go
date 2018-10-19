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
	"log"
	"sync"

	"gopkg.in/mgo.v2"
)

var instance *mgo.Session
var once sync.Once

func GetDb() *mgo.Session {
	once.Do(func() {
		session, err := mgo.Dial(util.Config.MongoUrl)
		if err != nil {
			log.Fatal("error on connection to mongodb: ", err)
		}
		session.SetMode(mgo.Monotonic, true)
		instance = session
	})
	return instance.Copy()
}

func getMongoFilterCollection() (session *mgo.Session, collection *mgo.Collection) {
	session = GetDb()
	collection = session.DB(util.Config.MongoTable).C(util.Config.MongoDeploymentCollection)
	return
}

func getMongoFilterPoolCollection() (session *mgo.Session, collection *mgo.Collection) {
	session = GetDb()
	collection = session.DB(util.Config.MongoTable).C(util.Config.MongoFilterPoolCollection)
	return
}
