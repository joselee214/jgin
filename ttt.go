package main

import (
	// "encoding/json"
	// "reflect"
	"strconv"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type A struct {
	A int
	B string
}


func mainzzz() {

		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
		    // handle error
		}
		defer c.Close()

sessionmaxLifetime, _ := strconv.Atoi("3600")
	aaa, errc := c.Do("SETEX", "0c6b8552a9895000ec4158f685753c43", sessionmaxLifetime,"11111{}")
	fmt.Println("zzzzz=>>>>",aaa,errc)

	// real := A{A:1,B:"zzz"}
	// reflected := reflect.New(reflect.TypeOf(real)).Elem().Interface()
	// fmt.Println(real, reflect.TypeOf(real),    reflect.TypeOf(&real),reflect.TypeOf(reflect.TypeOf(real)))
	// fmt.Println(reflected,reflect.TypeOf(reflected))

	// fmt.Println("=======")
	// fmt.Println("=======")

	// realx := A{A:1,B:"zzz"}
	// realxxxx := A{A:222,B:"zzz"}

	// jsbyte,_ := json.Marshal(realx)
	// jsonstr := string(jsbyte)

	// passed := &realx  //入参
	// passedxxx := &realxxxx

	// fmt.Println("----------1")
	// fmt.Println("realx:",realx,"        reflect.TypeOf(&realx):",reflect.TypeOf(&realx),"       jsonstr:",jsonstr )
	// fmt.Println("reflect.TypeOf(passed):",reflect.TypeOf(passed), "        		reflect.TypeOf(passed).Elem():",  reflect.TypeOf(passed).Elem(), "        		passed:", passed,"        		passedxxx:",passedxxx)

	// fmt.Println("===reflect.ValueOf(passed).Interface(): ",reflect.ValueOf(passed).Interface())

	//reflectedx := reflect.New(reflect.TypeOf(passed)).Elem().Interface()
	//
	//newed := &reflected
	//fmt.Println(reflectedx,reflect.TypeOf(reflectedx),reflectedx,newed,reflect.TypeOf(&reflected))
	//
	//fmt.Println("----------")
	//err:=json.Unmarshal([]byte(jsonstr),reflectedx)
	//fmt.Println(reflectedx,err,passedxxx)
	//err1:=json.Unmarshal([]byte(jsonstr),passedxxx)
	//fmt.Println(reflectedx,err1,passedxxx)

}
