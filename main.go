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
	NoOfEmployee float64       `json:"no_of_employee"`
	StrText      []string  `json:"str_text"`
	IntText      []float64     `json:"int_text"`
}

type address struct {
	Street   string `json:"street"`
	Landmark string `json:"landmark"`
	City     string `json:"city"`
	Pincode  float64    `json:"pincode"`
	State    string `json:"state"`
}

type sponser struct {
	Name string `json:"name"`
}

// type company struct {
// 	Name       string  `json:"name"`
// 	AgeInYears float64 `json:"age_in_years"`
// 	Origin     string  `json:"origin"`
// 	HeadOffice string  `json:"head_office"`
// 	Address    []struct {
// 		Street   string `json:"street"`
// 		Landmark string `json:"landmark"`
// 		City     string `json:"city"`
// 		Pincode  float64    `json:"pincode"`
// 		State    string `json:"state"`
// 	} `json:"address"`
// 	Sponsers struct {
// 		Name string `json:"name"`
// 	} `json:"sponsers"`
// 	Revenue      string   `json:"revenue"`
// 	NoOfEmployee float64      `json:"no_of_employee"`
// 	StrText      []string `json:"str_text"`
// 	IntText      []float64    `json:"int_text"`
// }


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

func setJSONKey(key string, val interface{}, src map[string]interface{}) (n int) {// n is used to count the total no of keys which are set by this func
	if _, ok := src[key]; ok { // if src has key then setting it directly
		src[key] = val
		n++
	}
	for _, v := range src {
		rVal := reflect.ValueOf(v)
		switch rVal.Kind() {
		case reflect.Map:
			n = n + setJSONKey(key, val, v.(map[string]interface{})) // recursive call for nested map
		case reflect.Slice:
			for _, m := range v.([]interface{}) {
				if reflect.ValueOf(m).Kind() == reflect.Map { // Slice can have map as element also
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


// PopulateStruct function populates fields of a struct using data from a map.
func populateStruct(src map[string]interface{}, st interface{}) {
	//getting reflection value of interface{} and then dereferencing to the structure that st points to
	stVal := reflect.ValueOf(st).Elem()

	for i:=0; i<stVal.NumField(); i++{
		//getting json tag of field of struct
		fJSONTag := stVal.Type().Field(i).Tag.Get("json")
		
		
		field := stVal.Field(i)
		v, ok := src[fJSONTag];
		if (!ok){
			//getting field name of struct if json Tag is not written
			fldName := stVal.Type().Field(i).Name
			v, ok = src[fldName]
		}
		if  ok {//if map has key of json tag value or field name
			if field.Type() == reflect.ValueOf(v).Type() {	
				field.Set(reflect.ValueOf(v))
			} else if field.Kind() == reflect.Slice{ 
				sliceElemType := field.Type().Elem()
				nestedSlice := reflect.MakeSlice(reflect.SliceOf(sliceElemType), len(v.([]interface{})), cap(v.([]interface{})))
				
				if sliceElemType.Kind() == reflect.Struct{
					//if it is slice of structs then call populateStruct recursively
					for i, ns := range v.([]interface{}) {
						populateStruct(ns.(map[string]interface{}), nestedSlice.Index(i).Addr().Interface())
					}
				}	else {
					for i, ns := range v.([]interface{}) {
						//if it is slice of types other than struct then
						//checking the type of each element and setting them if types are same
						if sliceElemType == reflect.TypeOf(ns){
							nestedSlice.Index(i).Set(reflect.ValueOf(ns))
						}
					}
				}
				field.Set(nestedSlice)

			} else if field.Kind() == reflect.Struct{
				if nestedMap , ok := v.(map[string]interface{}); ok {
					nestedStruct := reflect.New(field.Type()).Interface()
					populateStruct(nestedMap, nestedStruct)
					field.Set(reflect.ValueOf(nestedStruct).Elem())
				}
			}

		}
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
	fmt.Printf("\n\nStructure value after populating from json : \n\n%+v\n\n", i)
}
