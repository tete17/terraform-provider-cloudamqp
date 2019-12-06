package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (api *API) waitUntilReady(id string) (map[string]interface{}, error) {
	log.Printf("[DEBUG] go-api::instance::waitUntilReady waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		response, err := api.sling.New().Path("/api/instances/").Get(id).Receive(&data, &failed)

		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("waitUntilReady failed, status: %v, message: %s", response.StatusCode, failed))
		}
		if data["ready"] == true {
			data["id"] = id
			return data, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) waitUntilDeletion(id string) error {
	log.Printf("[DEBUG] go-api::instance::waitUntilDeletion waiting")
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	for {
		response, err := api.sling.New().Path("/api/instances/").Get(id).Receive(&data, &failed)

		if err != nil {
			return err
		}
		if response.StatusCode == 404 {
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::create params: %v", params)
	response, err := api.sling.New().Post("/api/instances").BodyJSON(params).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::waitUntilReady data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("CreateInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	data["id"] = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64)
	log.Printf("[DEBUG] go-api::instance::create id set: %v", data["id"])
	return api.waitUntilReady(data["id"].(string))
}

func (api *API) ReadInstance(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::read instance id: %v", id)
	response, err := api.sling.New().Path("/api/instances/").Get(id).Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) ReadInstances() ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	response, err := api.sling.New().Get("/api/instances").Receive(&data, &failed)
	log.Printf("[DEBUG] go-api::instance::read data: %v", data)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadInstances failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) UpdateInstance(id string, params map[string]interface{}) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::update instance id: %v, params: %v", id, params)
	response, err := api.sling.New().Put("/api/instances/"+id).BodyJSON(params).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("UpdateInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return nil
}

func (api *API) DeleteInstance(id string) error {
	failed := make(map[string]interface{})
	log.Printf("[DEBUG] go-api::instance::delete instance id: %v", id)
	response, err := api.sling.New().Path("/api/instances/").Delete(id).Receive(nil, &failed)

	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New(fmt.Sprintf("DeleteInstance failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilDeletion(id)
}
