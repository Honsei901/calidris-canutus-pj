package main

import (
	"calidris-canutus-pj/pb"
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id:          1,
		Name:        "Yamada",
		Email:       "test@example.com",
		Occupation:  pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project:     map[string]*pb.Company_Project{"ProjectX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My name is Taro",
		},
		Birthday: &pb.Date{
			Year:  1990,
			Month: 1,
			Day:   1,
		},
	}

	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("Can not serialize", err)
	}

	if err := os.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("Can not write", err)
	}

	in, err := os.ReadFile("test.bin")
	if err != nil {
		log.Fatalln("Can not read", err)
	}

	readEmployee := &pb.Employee{}
	if err = proto.Unmarshal(in, readEmployee); err != nil {
		log.Fatalln("Can not deserialize", err)
	}

	jsonData, err := protojson.Marshal(readEmployee)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(jsonData))

}
