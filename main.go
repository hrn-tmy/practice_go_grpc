package main

import (
	"fmt"
	"go-grpc/pb"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id:         1,
		Name:       "Suzuki",
		Email:      "test@example.com",
		Occupation: pb.Occupation_ENGINEER,
		Phone:      []string{"080-1234-5678", "090-1234-5678"},
		Project:    map[string]*pb.Company_Project{"ProjextX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My Name Is Suzuki",
		},
		Birthday: &pb.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
	}
	// バイナリ形式にシリアライズする
	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("Can't serialize", err)
	}
	if err := os.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("Can't Write", err)
	}

	in, err := os.ReadFile("test.bin")
	if err != nil {
		log.Fatalln("Can't read file", err)
	}

	readEmployee := &pb.Employee{}

	// バイナリ形式のファイルをデシリアライズする
	err = proto.Unmarshal(in, readEmployee)
	if err != nil {
		log.Fatalln("Can't read file", err)
	}
	fmt.Println(readEmployee)
}
