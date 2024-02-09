package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type company struct {
	Name         string    `json:"name"`
	AgeInYears   float64   `json:"age_in_years"`
	Origin       string    `json:"origin"`
	HeadOffice   string    `json:"head_office"`
	Address      []address `json:"address"`
	Sponsers     sponser   `json:"sponsers"`
	Revenue      string    `json:"revenue"`
	NoOfEmployee int       `json:"no_of_employee"`
	StrText      []string  `json:"str_text"`
	IntText      []int     `json:"int_text"`
}

type address struct {
	Street   string `json:"street"`
	Landmark string `json:"landmark"`
	City     string `json:"city"`
	Pincode  int    `json:"pincode"`
	State    string `json:"state"`
}

type sponser struct {
	Name string `json:"name"`
}

func printJSON(d interface{}, idn string) { //idn string is used for indentation purpose
	v := reflect.ValueOf(d) //getting value of data as reflection object

	switch v.Kind() {
	case reflect.Map:
		fmt.Printf("{")
		for _, k := range v.MapKeys() {
			fmt.Printf("\n%v\t%v : ", idn, k)
			printJSON(v.MapIndex(k).Interface(), idn+"\t") //Interface method converts the reflect.Value into an interface{},
		}
		fmt.Printf("\n    %v}", idn)
	case reflect.Slice:
		fmt.Printf("[")
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("\n%v\t%v : ", idn, i)
			printJSON(v.Index(i).Interface(), idn+"\t")
		}
		fmt.Printf("\n    %v]", idn)
	default:
		fmt.Printf("%v, (type : %v, kind : %v)", v.Interface(), v.Type(), v.Kind())
	}
}

func setJSONKey(key string, val interface{}, src map[string]interface{}) (n int) {
	if _, ok := src[key]; ok {
		src[key] = val
		n++
	}
	for _, v := range src {
		rVal := reflect.ValueOf(v)
		switch rVal.Kind() {
		case reflect.Map:
			n = n + setJSONKey(key, val, v.(map[string]interface{}))
		case reflect.Slice:
			for _, m := range v.([]interface{}) {
				if reflect.ValueOf(m).Kind() == reflect.Map {
					n = n + setJSONKey(key, val, m.(map[string]interface{}))
				}
			}
		}
	}
	return n
}

func deleteJSONKey(key string, src map[string]interface{}) (n int) {
	if _, ok := src[key]; ok {
		delete(src, key)
		n++
	}
	for _, v := range src {
		rVal := reflect.ValueOf(v)
		switch rVal.Kind() {
		case reflect.Map:
			n = n + deleteJSONKey(key, v.(map[string]interface{}))
		case reflect.Slice:
			for _, m := range v.([]interface{}) {
				if reflect.ValueOf(m).Kind() == reflect.Map {
					n = n + deleteJSONKey(key, m.(map[string]interface{}))
				}
			}
		}
	}
	return n
}

func populateStruct(src map[string]interface{}, i interface{}) {
	//marshaling the map[string]interface{} to jsonDataBytes
	jsonDataBytes, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}

	//unmarshaling the jsonDataBytes to structure
	err = json.Unmarshal(jsonDataBytes, i.(*company))
	if err != nil {
		panic(err)
	}
}

func main() {
	var inp string = `{
		"name" : "Tolexo Online Pvt. Ltd",
		"age_in_years" : 8.5,
		"origin" : "Noida",
		"head_office" : "Noida, Uttar Pradesh",
		"address" : [
			{
				"street" : "91 Springboard",
				"landmark" : "Axis Bank",
				"city" : "Noida",
				"pincode" : 201301,
				"state" : "Uttar Pradesh"
			},
			{
				"street" : "91 Springboard",
				"landmark" : "Axis Bank",
				"city" : "Noida",
				"pincode" : 201301,
				
				"state" : "Uttar Pradesh"
			}
		],
		"sponsers" : {
			"name" : "One"
		},
		"revenue" : "19.8 million$",
		"no_of_employee" : 630,
		"str_text" : ["one","two"],
		"int_text" : [1,3,4],
		"city": "abc"
	}`
	var mp map[string]interface{}
	err := json.Unmarshal([]byte(inp), &mp) // decode JSON data into interface{}

	if err != nil {
		panic(err)
	}

	new_city := "New Delhi"
	set_key := "city"
	fmt.Print("\nMap before setting key :\n")
	printJSON(mp, "")

	fmt.Printf("\n\nFound %v values having key = %v and set to %v successfully", setJSONKey(set_key, new_city, mp), set_key, new_city)
	fmt.Print("\n\nMap after setting key :\n")
	printJSON(mp, "")

	fmt.Printf("\n\nFound %v values having key = %v and deleted successfully", deleteJSONKey(set_key, mp), set_key)
	fmt.Print("\n\nMap after deleting key :\n")
	printJSON(mp, "")

	var i company
	populateStruct(mp, &i)
	fmt.Printf("\n\nStructure value after populating from json : %v", i)
}
