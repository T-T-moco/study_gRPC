package main

import (
	"io/ioutil"
	"log"
	"protobuf-lesson/pb"

	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id:          1,
		Name:        "Suzuki",
		Email:       "test@example.com",
		Occupation:  pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project:     map[string]*pb.Company_Project{"ProjectX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My name is Suzuki",
		},
		Birthday: &pb.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
	}

	binData, err := proto.Marshal(employee) // シリアライズしたい構造体を渡す
	if err != nil {
		log.Fatalln("Cant serialize", err)
	}

	if err := ioutil.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("Cant write", err)
	}
}
