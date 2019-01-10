package lib

import (
	"flag"
	"github.com/SmartEnergyPlatform/event-broker/util"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
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

	err = amqpEventHandling([]byte(event_example))

	if err != nil {
		t.Error(err)
		return
	}

	filter, err := GetFilter("fid_64eef49e-8d4d-493b-992c-29a7179ecfd9")

	if err != nil {
		t.Error(err)
		return
	}

	if filter.ProcessId == "" || filter.FilterId != "fid_64eef49e-8d4d-493b-992c-29a7179ecfd9" {
		t.Error("unexpected result", filter)
		return
	}
}







const event_example = `{"command":"PUT","id":"baeda93a-177f-4a89-8899-025d41875555","owner":"ad0391ee-a6ad-438f-8dbf-d8d2cac5005c","deployment":{"process":{"xml":"","name":"alarm_test","abstract_tasks":[{"tasks":[{"id":"Task_10avate","values":{"inputs":{"value":{"b":0,"bri":0,"g":0,"on":true,"r":0}},"outputs":{"status":0}},"parameter":[],"service_id":"iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed","label":"set_state"}],"device_type_id":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","options":[{"id":"iot#829fcf4f-b5f4-4daa-a15d-e52b4b762fe3","name":"Hue color lamp 3","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:0a:70:e7-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#7d2dbdeb-5e63-48e9-9d50-fb44ecb706fd","name":"Hue color lamp 2","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:4e:a3:7c-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"}],"selected":{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":[],"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},"state":""},{"tasks":[{"id":"Task_1l82idm","values":{"inputs":{"value":{"b":0,"bri":0,"g":0,"on":true,"r":0}},"outputs":{"status":0}},"parameter":[],"service_id":"iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed","label":"set_state"}],"device_type_id":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","options":[{"id":"iot#829fcf4f-b5f4-4daa-a15d-e52b4b762fe3","name":"Hue color lamp 3","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:0a:70:e7-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#7d2dbdeb-5e63-48e9-9d50-fb44ecb706fd","name":"Hue color lamp 2","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:4e:a3:7c-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"}],"selected":{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":[],"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},"state":""},{"tasks":[{"id":"Task_062eahn","values":{},"parameter":[],"service_id":"iot#deb8c74e-d556-4ab6-b2d8-4f1a74d528a6","label":"on"}],"device_type_id":"iot#7075b34a-23b5-49ba-9723-867613205e0c","options":[{"id":"iot#89f3b32e-5255-4537-b0ef-14ca61abe012","name":"Fibaro Switch (#4)","device_type":"iot#7075b34a-23b5-49ba-9723-867613205e0c","uri":"ZWAY_183e91ad-1aaa-44af-9816-b6cc98889a79_ZWayVDev_zway_4-0-37","tags":["zway_device_group:Fibaro Wall Plug FGWPx-102 ZW5"],"user_tags":null,"gateway":"iot#6e00e48e-8f31-4fdf-85f4-fe3e5c93fda7","img":"https://i.imgur.com/EQrfmsX.png"},{"id":"iot#e1aeb45b-8f03-458b-b14c-f770dc3fb13c","name":"Siren (#5)","device_type":"iot#7075b34a-23b5-49ba-9723-867613205e0c","uri":"ZWAY_183e91ad-1aaa-44af-9816-b6cc98889a79_ZWayVDev_zway_5-0-37","tags":["zway_device_group:Aeotec Indoor Siren"],"user_tags":null,"gateway":"iot#6e00e48e-8f31-4fdf-85f4-fe3e5c93fda7","img":"https://i.imgur.com/EQrfmsX.png"}],"selected":{"id":"iot#e1aeb45b-8f03-458b-b14c-f770dc3fb13c","name":"Siren (#5)","device_type":"iot#7075b34a-23b5-49ba-9723-867613205e0c","uri":"ZWAY_183e91ad-1aaa-44af-9816-b6cc98889a79_ZWayVDev_zway_5-0-37","tags":["zway_device_group:Aeotec Indoor Siren"],"user_tags":[],"gateway":"iot#6e00e48e-8f31-4fdf-85f4-fe3e5c93fda7","img":"https://i.imgur.com/EQrfmsX.png"},"state":""},{"tasks":[{"id":"Task_1s7fede","values":{"inputs":{"value":{"b":0,"bri":0,"g":0,"on":true,"r":0}},"outputs":{"status":0}},"parameter":[],"service_id":"iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed","label":"set_state"}],"device_type_id":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","options":[{"id":"iot#829fcf4f-b5f4-4daa-a15d-e52b4b762fe3","name":"Hue color lamp 3","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:0a:70:e7-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#7d2dbdeb-5e63-48e9-9d50-fb44ecb706fd","name":"Hue color lamp 2","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:4e:a3:7c-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"}],"selected":{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":[],"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},"state":""},{"tasks":[{"id":"Task_1w1172b","values":{"inputs":{"value":{"b":0,"bri":0,"g":0,"on":true,"r":0}},"outputs":{"status":0}},"parameter":[],"service_id":"iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed","label":"set_state"}],"device_type_id":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","options":[{"id":"iot#829fcf4f-b5f4-4daa-a15d-e52b4b762fe3","name":"Hue color lamp 3","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:0a:70:e7-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},{"id":"iot#7d2dbdeb-5e63-48e9-9d50-fb44ecb706fd","name":"Hue color lamp 2","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:4e:a3:7c-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":null,"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"}],"selected":{"id":"iot#6d8339da-b4bf-4e85-9909-d3e89440bf29","name":"Hue color lamp 1","device_type":"iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d","uri":"00:17:88:01:02:07:ba:ce-0b","tags":["manufacturer:Philips","type:Extended color light"],"user_tags":[],"gateway":"iot#2fc8d6a6-f223-4050-aad6-56373de47140","img":"https://i.imgur.com/OZOqLcR.png"},"state":""}],"abstract_data_export_tasks":[],"receive_tasks":null,"msg_events":[{"filter_id":"fid_64eef49e-8d4d-493b-992c-29a7179ecfd9","shape_id":"IntermediateCatchEvent_03l35pd","filter":{"scope":"all","rules":[{"path":"$.value.metrics.level","scope":"all","operator":"==","value":"on"}],"topic":"","device_id":"iot#3eac0cd1-1b58-4fd7-a106-4f7d501174ca","service_id":"iot#1710c870-04c8-403f-8788-95f953eccf1e"}},{"filter_id":"fid_9fe0099d-a691-41a5-8927-7ddd390cf16a","shape_id":"IntermediateCatchEvent_0isg4wc","filter":{"scope":"all","rules":[{"path":"$.value.metrics.level","scope":"all","operator":"==","value":"on"}],"topic":"","device_id":"iot#4d01a743-5c0d-41a6-8e1c-07e9a09e58a7","service_id":"iot#e1c204c4-7961-4cee-8a5c-3e9f0aac92c2"}}],"time_events":[]},"svg":""}}`