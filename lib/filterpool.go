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
	"errors"
	"math/rand"
	"strconv"

	"time"

	"gopkg.in/mgo.v2/bson"
)

const POOL_COMMAND_NONE string = "none"
const POOL_COMMAND_ADD string = "add"
const POOL_COMMAND_REMOVE string = "remove"
const POOL_COMMAND_RESET string = "reset"

const DEPLOYMENT_STOPPING string = "stopping"
const DEPLOYMENT_STARTING string = "starting"
const DEPLOYMENT_RUNNING string = "running"

func GetFilterPoolCommand(poolid string, poolsize string) (result FilterPoolCommand, err error) {
	err = updateKnownFilterPools(poolid)
	if err != nil {
		return
	}
	poolIsUnsync, err := filterPoolIsUnsync(poolid, poolsize)
	if err != nil {
		return result, err
	}
	if poolIsUnsync {
		result.Command = POOL_COMMAND_RESET
		result.Assignment, err = GetFilterPoolAssignments(poolid)
		return result, err
	}

	starting, err := getStartingDeployments(poolid)
	if len(starting) > 0 {
		err = setDeploymentsRunning(poolid)
		if err != nil {
			return result, err
		}
		result.Command = POOL_COMMAND_ADD
		result.Assignment = starting
		return result, err
	}

	stopping, err := getStoppingDeployments(poolid)
	if len(stopping) > 0 {
		err = removeStoppingDeployments(poolid)
		if err != nil {
			return result, err
		}
		result.Command = POOL_COMMAND_REMOVE
		result.Assignment = stopping
		return result, err
	}

	result.Command = POOL_COMMAND_NONE
	return result, err
}

func GetFilterPoolAssignments(poolid string) (result []FilterDeployment, err error) {
	err = updateKnownFilterPools(poolid)
	if err != nil {
		return
	}

	err = removeStoppingDeployments(poolid)
	if err != nil {
		return
	}
	running, err := getRunningDeployments(poolid)
	if err != nil {
		return result, err
	}
	starting, err := getStartingDeployments(poolid)
	if err != nil {
		return result, err
	}

	err = setDeploymentsRunning(poolid)
	if err != nil {
		return result, err
	}

	result = append(running, starting...)
	return
}

func RemoveFilterPool(poolid string) (err error) {
	session, collection := getMongoFilterPoolCollection()
	defer session.Close()
	err = collection.Update(KnownFilterPool{FilterPool: poolid}, KnownFilterPool{FilterPool: poolid, Active: false})
	if err != nil {
		return err
	}
	running, err := getRunningDeployments(poolid)
	if err != nil {
		return err
	}
	starting, err := getStartingDeployments(poolid)
	if err != nil {
		return err
	}
	for _, deployment := range append(running, starting...) {
		_, err = reassignFilter(deployment.FilterId)
		if err != nil {
			return err
		}
	}
	err = removeStoppingDeployments(poolid)
	if err != nil {
		return err
	}
	return collection.Remove(KnownFilterPool{FilterPool: poolid})
}

func removeStoppingDeployments(poolid string) (err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	_, err = collection.RemoveAll(bson.M{"state": DEPLOYMENT_STOPPING, "filter_pool": poolid})
	return
}

func getStoppingDeployments(poolid string) (result []FilterDeployment, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	err = collection.Find(bson.M{"state": DEPLOYMENT_STOPPING, "filter_pool": poolid}).All(&result)
	return
}

func setDeploymentsRunning(poolid string) (err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	_, err = collection.UpdateAll(bson.M{"state": DEPLOYMENT_STARTING, "filter_pool": poolid}, bson.M{"$set": bson.M{"state": DEPLOYMENT_RUNNING}})
	return
}

func getStartingDeployments(poolid string) (result []FilterDeployment, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	err = collection.Find(bson.M{"state": DEPLOYMENT_STARTING, "filter_pool": poolid}).All(&result)
	return
}

func filterPoolIsUnsync(poolid string, transmittedPoolSize string) (result bool, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	count_running, err := collection.Find(bson.M{"state": DEPLOYMENT_RUNNING, "filter_pool": poolid}).Count()
	if err != nil {
		return
	}
	count_stopping, err := collection.Find(bson.M{"state": DEPLOYMENT_STOPPING, "filter_pool": poolid}).Count()
	if err != nil {
		return
	}
	poolSize, err := strconv.Atoi(transmittedPoolSize)
	return poolSize != count_stopping+count_running, err
}

func updateKnownFilterPools(poolid string) error {
	session, collection := getMongoFilterPoolCollection()
	defer session.Close()
	_, err := collection.Upsert(KnownFilterPool{FilterPool: poolid}, KnownFilterPool{FilterPool: poolid, LastContact: time.Now().Unix(), Active: true})
	return err
}

func getRunningDeployments(poolid string) (result []FilterDeployment, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	err = collection.Find(bson.M{"state": DEPLOYMENT_RUNNING, "filter_pool": poolid}).All(&result)
	return
}

func selectFilterPool() (result string, err error) {
	session, collection := getMongoFilterPoolCollection()
	defer session.Close()
	pools := []KnownFilterPool{}
	err = collection.Find(KnownFilterPool{Active: true}).All(&pools)
	if err != nil {
		return
	}
	if len(pools) == 0 {
		err = errors.New("no event filter pool available")
		return
	}
	result = pools[rand.Intn(len(pools))].FilterPool
	return
}

func reassignFilter(filterId string) (poolId string, err error) {
	session, collection := getMongoFilterCollection()
	defer session.Close()
	poolId, err = selectFilterPool()
	if err != nil {
		return
	}
	err = collection.Update(FilterDeployment{FilterId: filterId}, FilterDeployment{FilterPool: poolId, State: DEPLOYMENT_STARTING})
	return
}
