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
	"io/ioutil"
	"log"
)

func GetSwagger() (string, error) {
	file, err := ioutil.ReadFile(util.Config.SwaggerLocation)
	if err != nil {
		log.Println("error on config load: ", err)
		return "", err
	}
	return string(file), err
}
